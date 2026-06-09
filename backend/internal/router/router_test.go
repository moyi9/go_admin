package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"encoding/json"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go_web/internal/config"
	"go_web/internal/modules/audit"
	"go_web/internal/modules/auth"
	"go_web/internal/modules/notification"
	"go_web/internal/modules/rbac"
	"go_web/internal/pkg/security"
	"gorm.io/gorm"
)

func TestHealthz(t *testing.T) {
	engine := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestProtectedRouteRequiresJWT(t *testing.T) {
	engine := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginSuccess(t *testing.T) {
	engine := newTestRouter(t)

	body := strings.NewReader(`{"username":"admin","password":"admin123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
			User        struct {
				ID       uint   `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"data"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, 0, resp.Code)
	require.NotEmpty(t, resp.Data.AccessToken)
	require.Equal(t, "Bearer", resp.Data.TokenType)
	require.Equal(t, "admin", resp.Data.User.Username)
}

func TestLoginWrongPassword(t *testing.T) {
	engine := newTestRouter(t)

	body := strings.NewReader(`{"username":"admin","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid username or password")
}

func TestAuthMeWithValidToken(t *testing.T) {
	engine := newTestRouter(t)

	// 先登录获取 token
	loginBody := strings.NewReader(`{"username":"admin","password":"admin123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", loginBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	var loginResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&loginResp)

	// 用 token 访问 /auth/me
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req2.Header.Set("Authorization", "Bearer "+loginResp.Data.AccessToken)
	rec2 := httptest.NewRecorder()
	engine.ServeHTTP(rec2, req2)

	require.Equal(t, http.StatusOK, rec2.Code)

	var meResp struct {
		Data struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}
	json.NewDecoder(rec2.Body).Decode(&meResp)
	require.Equal(t, "admin", meResp.Data.Username)
	require.Equal(t, "admin@example.com", meResp.Data.Email)
}

func TestListUsersWithPagination(t *testing.T) {
	engine := newTestRouter(t)

	// 登录
	loginBody := strings.NewReader(`{"username":"admin","password":"admin123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", loginBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	var loginResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&loginResp)

	// 分页查用户
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&page_size=10", nil)
	req2.Header.Set("Authorization", "Bearer "+loginResp.Data.AccessToken)
	rec2 := httptest.NewRecorder()
	engine.ServeHTTP(rec2, req2)

	require.Equal(t, http.StatusOK, rec2.Code)
	require.Contains(t, rec2.Body.String(), `"total":`)
	require.Contains(t, rec2.Body.String(), `"page":`)
	require.Contains(t, rec2.Body.String(), `"page_size":`)
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(rbac.Models()...))
	require.NoError(t, db.AutoMigrate(audit.Models()...))
	require.NoError(t, db.AutoMigrate(notification.Models()...))
	require.NoError(t, rbac.Seed(db, config.SeedConfig{
		Enabled:       true,
		AdminUsername: "admin",
		AdminEmail:    "admin@example.com",
		AdminPassword: "admin123456",
	}))

	enforcer, err := casbin.NewEnforcer("../../configs/casbin_model.conf")
	require.NoError(t, err)
	jwtManager := security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})
	repo := rbac.NewRepository(db)
	rbacService := rbac.NewService(repo, enforcer)
	require.NoError(t, rbacService.SyncPolicies())
	auditRepo := audit.NewRepository(db)
	auditService := audit.NewService(auditRepo)
	notifRepo := notification.NewRepository(db)
	notifService := notification.NewService(notifRepo)
	authService := auth.NewService(repo, jwtManager)

	return New(Deps{
		Config: &config.Config{
			Server:  config.ServerConfig{Mode: "test"},
			CORS:    config.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Authorization", "Content-Type", "X-Request-ID"}},
			Swagger: config.SwaggerConfig{Enabled: false},
		},
		Logger:      zap.NewNop(),
		JWT:         jwtManager,
		Enforcer:    enforcer,
		AuthHandler:  auth.NewHandler(authService, auditService, notifService),
		RBACHandler: rbac.NewHandler(rbacService, auditService),
	})
}
