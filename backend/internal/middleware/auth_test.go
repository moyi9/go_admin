package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go_web/internal/config"
	"go_web/internal/pkg/security"
)

func TestJWTAuthRejectsMissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})

	r := gin.New()
	r.Use(JWTAuth(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "missing authorization header")
}

func TestJWTAuthRejectsMalformedHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})

	r := gin.New()
	r.Use(JWTAuth(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	tests := []struct {
		name   string
		header string
	}{
		{"basic auth instead of bearer", "Basic YWxhZGRpbjpvcGVuc2VzYW1l"},
		{"no bearer prefix", "some-token-value"},
		{"empty bearer", "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.header)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, http.StatusUnauthorized, rec.Code)
			require.Contains(t, rec.Body.String(), "invalid authorization header")
		})
	}
}

func TestJWTAuthRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})

	r := gin.New()
	r.Use(JWTAuth(jwtManager))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid token")
}

func TestJWTAuthSetsContextValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})

	token, err := jwtManager.Generate(42, "alice", []string{"admin", "editor"})
	require.NoError(t, err)

	r := gin.New()
	r.Use(JWTAuth(jwtManager))
	r.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		roles, _ := c.Get("roles")
		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
			"roles":    roles,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, `"user_id":42`)
	require.Contains(t, body, `"username":"alice"`)
	require.Contains(t, body, `"roles":["admin","editor"]`)
}
