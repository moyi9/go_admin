package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadReadsYAMLAndEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(path, []byte(`
server:
  port: 8080
jwt:
  secret: from-file
  access_token_ttl: 1h
postgres:
  host: localhost
  port: 5432
  user: app
  password: pass
  database: app
redis:
  addr: localhost:6379
`), 0600)
	require.NoError(t, err)

	t.Setenv("GO_WEB_SERVER_PORT", "9090")
	t.Setenv("GO_WEB_JWT_SECRET", "from-env")

	cfg, err := Load(path)

	require.NoError(t, err)
	require.Equal(t, 9090, cfg.Server.Port)
	require.Equal(t, "from-env", cfg.JWT.Secret)
	require.Equal(t, "localhost:6379", cfg.Redis.Addr)
}
