package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
)

// Recovery 捕获 panic 并返回统一错误响应，避免服务进程因单个请求崩溃。
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered",
			zap.Any("recovered", recovered),
			zap.String("request_id", requestID(c)),
		)
		response.Error(c, http.StatusInternalServerError, apperror.CodeInternal, "internal server error")
	})
}
