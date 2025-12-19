package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvStr_PanicsOnEmpty(t *testing.T) {
	key := "MISSING_ENV_VAR"
	_ = os.Unsetenv(key)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing env var '%s', but did not panic", key)
		}
	}()

	_ = LoadEnvStr(key)
}

func TestLoadEnvStr_ReturnsValue(t *testing.T) {
	key := "EXISTING_ENV_VAR"
	val := "test_value"
	_ = os.Setenv(key, val)

	result := LoadEnvStr(key)
	assert.Equal(t, val, result)
}

func TestLoadEnvConfiguration(t *testing.T) {
	_ = os.Setenv(DB_HOST, "localhost")
	_ = os.Setenv(DB_USER, "postgres")
	_ = os.Setenv(DB_PASSWORD, "secret")
	_ = os.Setenv(DB_NAME, "testdb")
	_ = os.Setenv(DB_PORT, "5432")
	_ = os.Setenv(DB_SCHEMA, "public")
	_ = os.Setenv(Env, "test")
	_ = os.Setenv(PORT, "8080")

	cfg := LoadEnvConfiguration()

	assert.Equal(t, "localhost", cfg.DBConnectionParams.Host)
	assert.Equal(t, "postgres", cfg.DBConnectionParams.User)
	assert.Equal(t, "secret", cfg.DBConnectionParams.Password)
	assert.Equal(t, "testdb", cfg.DBConnectionParams.DBName)
	assert.Equal(t, "5432", cfg.DBConnectionParams.DBPort)
	assert.Equal(t, "public", cfg.DBConnectionParams.DBSchema)

	assert.Equal(t, "test", cfg.ServiceParams.Env)
	assert.Equal(t, "8080", cfg.ServiceParams.Port)
}
