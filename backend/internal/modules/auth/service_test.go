package auth

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"go_web/internal/config"
	"go_web/internal/modules/rbac"
	"go_web/internal/pkg/security"
	"gorm.io/gorm"
)

func TestLoginWithSeededAdmin(t *testing.T) {
	db := newTestDB(t)
	err := rbac.Seed(db, config.SeedConfig{
		Enabled:       true,
		AdminUsername: "admin",
		AdminEmail:    "admin@example.com",
		AdminPassword: "admin123456",
	})
	require.NoError(t, err)

	repo := rbac.NewRepository(db)
	service := NewService(repo, security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	}))

	result, err := service.Login(LoginRequest{Username: "admin", Password: "admin123456"})

	require.NoError(t, err)
	require.NotEmpty(t, result.AccessToken)
	require.Equal(t, "admin", result.User.Username)
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	db := newTestDB(t)
	err := rbac.Seed(db, config.SeedConfig{
		Enabled:       true,
		AdminUsername: "admin",
		AdminEmail:    "admin@example.com",
		AdminPassword: "admin123456",
	})
	require.NoError(t, err)

	repo := rbac.NewRepository(db)
	service := NewService(repo, security.NewJWTManager(config.JWTConfig{
		Issuer:         "go_web",
		Secret:         "test-secret",
		AccessTokenTTL: time.Hour,
	}))

	_, err = service.Login(LoginRequest{Username: "admin", Password: "wrong-password"})

	require.Error(t, err)
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(rbac.Models()...))
	return db
}
