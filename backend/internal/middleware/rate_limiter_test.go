package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupRateLimiterTest(rate int, window time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(NewRateLimiter(rate, window).Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func doRequest(t *testing.T, r *gin.Engine) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec.Code
}

func TestRateLimiterAllowsWithinLimit(t *testing.T) {
	r := setupRateLimiterTest(3, time.Minute)

	for i := 0; i < 3; i++ {
		require.Equal(t, http.StatusOK, doRequest(t, r), "request %d should be allowed", i+1)
	}
}

func TestRateLimiterBlocksExceedingLimit(t *testing.T) {
	r := setupRateLimiterTest(2, time.Minute)

	require.Equal(t, http.StatusOK, doRequest(t, r))
	require.Equal(t, http.StatusOK, doRequest(t, r))
	require.Equal(t, http.StatusTooManyRequests, doRequest(t, r))
}

func TestRateLimiterSeparatesByIP(t *testing.T) {
	r := setupRateLimiterTest(1, time.Minute)

	// IP 1 用完额度
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "10.0.0.1:11111"
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)
	require.Equal(t, http.StatusOK, rec1.Code)

	req1b := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1b.RemoteAddr = "10.0.0.1:11111"
	rec1b := httptest.NewRecorder()
	r.ServeHTTP(rec1b, req1b)
	require.Equal(t, http.StatusTooManyRequests, rec1b.Code)

	// IP 2 不受影响
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "10.0.0.2:22222"
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	require.Equal(t, http.StatusOK, rec2.Code)
}

func TestRateLimiterNewWindowResets(t *testing.T) {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    2,
		window:  200 * time.Millisecond,
	}

	key := "test-ip"
	require.True(t, rl.allow(key))
	require.True(t, rl.allow(key))
	require.False(t, rl.allow(key))

	time.Sleep(300 * time.Millisecond)

	require.True(t, rl.allow(key), "should reset after window passes")
}
