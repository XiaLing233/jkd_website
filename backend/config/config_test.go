package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupEnv 备份并覆盖环境变量，测试结束后自动恢复。
func setupEnv(t *testing.T, vals map[string]string) {
	t.Helper()
	for k, v := range vals {
		t.Setenv(k, v)
	}
}

func TestLoadDefaults(t *testing.T) {
	setupEnv(t, map[string]string{
		"DB_HOST":        "testhost",
		"DB_PORT":        "9999",
		"DB_R_USER":      "testuser",
		"DB_R_PASSWORD":  "testpass",
		"DB_META":        "test_meta",
		"SERVER_PORT":    "9999",
	})

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "testhost", cfg.DBHost)
	assert.Equal(t, "9999", cfg.DBPort)
	assert.Equal(t, "testuser", cfg.DBRUser)
	assert.Equal(t, "testpass", cfg.DBRPass)
	assert.Equal(t, "test_meta", cfg.DBMeta)
	assert.Equal(t, "9999", cfg.ServerPort)
}

func TestLoadMissingPassword(t *testing.T) {
	os.Unsetenv("DB_R_PASSWORD")
	_, err := Load()
	assert.Error(t, err)
}

func TestMetaDSN(t *testing.T) {
	cfg := &Config{
		DBRUser: "root", DBRPass: "secret",
		DBHost: "localhost", DBPort: "3306",
		DBMeta: "test_meta",
	}
	assert.Contains(t, cfg.MetaDSN(), "root:secret@tcp(localhost:3306)/test_meta")
	assert.Contains(t, cfg.MetaDSN(), "charset=utf8mb4")
}

func TestCalendarDSN(t *testing.T) {
	cfg := &Config{
		DBRUser: "root", DBRPass: "secret",
		DBHost: "xk-mysql", DBPort: "3306",
	}
	assert.Contains(t, cfg.CalendarDSN("calendar_122_b"),
		"root:secret@tcp(xk-mysql:3306)/calendar_122_b")
}
