package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
	"go_web/internal/pkg/security"
)

// JWTAuth 校验 Bearer Token，并把 user_id、username、roles 写入 Gin Context。
func JWTAuth(jwtManager *security.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, apperror.CodeUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, apperror.CodeUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtManager.Parse(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, apperror.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}
