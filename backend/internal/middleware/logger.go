package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AccessLog 记录每个 HTTP 请求的核心信息，便于排查慢请求和错误请求。
func AccessLog(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info("http request",
			zap.String("request_id", requestID(c)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func requestID(c *gin.Context) string {
	value, _ := c.Get("request_id")
	if s, ok := value.(string); ok {
		return s
	}
	return c.GetHeader(RequestIDHeader)
}
