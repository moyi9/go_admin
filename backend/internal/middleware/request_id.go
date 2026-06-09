package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID 复用客户端传入的 X-Request-ID，或为请求生成新的 ID。
// 后续日志、错误响应和链路排查都会使用同一个 request_id。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("request_id", requestID)
		c.Header(RequestIDHeader, requestID)
		c.Next()
	}
}
