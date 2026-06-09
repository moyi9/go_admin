package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go_web/internal/config"
)

func TestJWTGenerateAndParse(t *testing.T) {
	manager := NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})

	token, err := manager.Generate(1, "admin", []string{"admin"})
	require.NoError(t, err)

	claims, err := manager.Parse(token)
	require.NoError(t, err)
	require.Equal(t, uint(1), claims.UserID)
	require.Equal(t, "admin", claims.Username)
	require.Equal(t, []string{"admin"}, claims.Roles)
}

func TestJWTRejectsInvalidToken(t *testing.T) {
	manager := NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	})

	_, err := manager.Parse("not-a-token")

	require.Error(t, err)
}
