package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// setupAuthorizeTest 创建带有 Authorize 中间件的测试路由。
// userID 通过一个前置中间件注入，模拟真实 JWTAuth → Authorize 的调用链。
func setupAuthorizeTest(t *testing.T, userIDToSet any) (*casbin.Enforcer, *gin.Engine) {
	t.Helper()
	enforcer, err := casbin.NewEnforcer("../../configs/casbin_model.conf")
	require.NoError(t, err)

	// userID=1 → role:admin → GET + POST /api/v1/users
	_, err = enforcer.AddGroupingPolicy("1", "admin")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("admin", "/api/v1/users", "GET")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("admin", "/api/v1/users", "POST")
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 注入 user_id（模拟 JWT 中间件行为）
	r.Use(func(c *gin.Context) {
		if userIDToSet != nil {
			c.Set("user_id", userIDToSet)
		}
		c.Next()
	})
	r.Use(Authorize(enforcer))

	r.GET("/api/v1/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"users": true})
	})
	r.POST("/api/v1/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"created": true})
	})

	return enforcer, r
}

func TestAuthorizeAllowsMatchingPolicy(t *testing.T) {
	_, r := setupAuthorizeTest(t, uint(1))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthorizeRejectsMissingUserID(t *testing.T) {
	_, r := setupAuthorizeTest(t, nil) // 不注入 user_id

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthorizeRejectsWrongUserIDType(t *testing.T) {
	_, r := setupAuthorizeTest(t, "string-instead-of-uint")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid user id")
}

func TestAuthorizeDeniesNoMatchingPolicy(t *testing.T) {
	// userID=2 没有分配任何角色，即使请求存在的路由也会被拒
	_, r := setupAuthorizeTest(t, uint(2))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorizeDeniesWrongMethod(t *testing.T) {
	// userID=1 有 GET 但没有 DELETE
	_, r := setupAuthorizeTest(t, uint(1))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthorizeEnforcerIntegration(t *testing.T) {
	enforcer, err := casbin.NewEnforcer("../../configs/casbin_model.conf")
	require.NoError(t, err)

	_, err = enforcer.AddGroupingPolicy("42", "editor")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("editor", "/api/v1/users/:id", "PUT")
	require.NoError(t, err)

	// keyMatch2 应匹配 :id 参数
	allowed, err := enforcer.Enforce("42", "/api/v1/users/5", "PUT")
	require.NoError(t, err)
	require.True(t, allowed)

	allowed, err = enforcer.Enforce("42", "/api/v1/users/5", "DELETE")
	require.NoError(t, err)
	require.False(t, allowed)
}

func TestAuthorizeUserIDFormatRoundTrip(t *testing.T) {
	// 确保 middleware 中 fmt.Sprintf("%%d", uint) 的输出与 Casbin 策略一致
	userID := uint(7)
	formatted := fmt.Sprintf("%d", userID)
	require.Equal(t, "7", formatted)

	enforcer, err := casbin.NewEnforcer("../../configs/casbin_model.conf")
	require.NoError(t, err)
	_, err = enforcer.AddGroupingPolicy("7", "viewer")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("viewer", "/api/v1/roles", "GET")
	require.NoError(t, err)

	allowed, err := enforcer.Enforce(formatted, "/api/v1/roles", "GET")
	require.NoError(t, err)
	require.True(t, allowed)
}
