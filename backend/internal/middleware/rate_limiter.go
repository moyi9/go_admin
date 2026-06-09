package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
)

// RateLimiter 基于内存的滑动窗口限流器。
// 按客户端 IP 限制请求频率，适用于登录等敏感接口的防暴力破解。
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    int
	window  time.Duration
}

type bucket struct {
	count    int
	windowID int64 // 当前窗口的起始 unix 秒
}

// NewRateLimiter 创建一个限流器。
// rate: 每个 window 时间窗口内允许的最大请求数。
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		window:  window,
	}
	// 后台定期清理过期条目，防止内存泄漏。
	go rl.cleanup()
	return rl
}

// Middleware 返回 Gin 中间件处理函数。
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !rl.allow(key) {
			response.Error(c, http.StatusTooManyRequests, apperror.CodeTooManyRequests, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().UnixNano()
	windowNanos := int64(rl.window)
	if windowNanos < 1 {
		windowNanos = 1
	}
	currentWindow := now / windowNanos

	b, ok := rl.buckets[key]
	if !ok || b.windowID != currentWindow {
		rl.buckets[key] = &bucket{count: 1, windowID: currentWindow}
		return true
	}

	if b.count >= rl.rate {
		return false
	}
	b.count++
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now().UnixNano()
		windowNanos := int64(rl.window)
		if windowNanos < 1 {
			windowNanos = 1
		}
		currentWindow := now / windowNanos
		for key, b := range rl.buckets {
			if currentWindow-b.windowID > 1 {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}
