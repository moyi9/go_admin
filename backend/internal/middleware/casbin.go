package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
)

// Authorize 使用 Casbin 校验当前用户是否允许访问当前路由和 HTTP 方法。
// 权限粒度是 method + Gin 路由模板，例如 GET /api/v1/users/:id。
func Authorize(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user_id")
		if !ok {
			response.Error(c, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		userID, ok := v.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, apperror.CodeUnauthorized, "invalid user id")
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(fmt.Sprintf("%d", userID), c.FullPath(), c.Request.Method)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, apperror.CodeInternal, "authorization error")
			c.Abort()
			return
		}
		if !allowed {
			response.Error(c, http.StatusForbidden, apperror.CodeForbidden, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
