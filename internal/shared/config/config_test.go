package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Interface(t *testing.T) {
	t.Run("config implements ConfigReader interface", func(t *testing.T) {
		cfg := &config{
			server: serverConfig{
				Host: "localhost",
				Port: 8080,
			},
			database: databaseConfig{
				DSN: "postgres://test:test@localhost:5432/test",
			},
			security: securityConfig{
				JWTSecret:                  "test-secret",
				EncryptionKey:              "test-key",
				JWTExpireDuration:          24 * time.Hour,
				RefreshTokenExpireDuration: 30 * 24 * time.Hour,
			},
			logging: loggingConfig{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
		}

		var reader ConfigReader = cfg
		assert.NotNil(t, reader)
	})
}

func TestConfig_Methods(t *testing.T) {
	cfg := &config{
		server: serverConfig{
			Host: "localhost",
			Port: 8080,
		},
		database: databaseConfig{
			DSN: "postgres://test:test@localhost:5432/test",
		},
		security: securityConfig{
			JWTSecret:                  "test-secret",
			EncryptionKey:              "test-key",
			JWTExpireDuration:          24 * time.Hour,
			RefreshTokenExpireDuration: 30 * 24 * time.Hour,
		},
		logging: loggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}

	t.Run("GetDSN", func(t *testing.T) {
		dsn := cfg.GetDSN()
		assert.NotEmpty(t, dsn)
		assert.Contains(t, dsn, "postgres://")
	})

	t.Run("GetServerAddr", func(t *testing.T) {
		addr := cfg.GetServerAddr()
		assert.Equal(t, "localhost:8080", addr)
	})

	t.Run("GetJWTSecret", func(t *testing.T) {
		secret := cfg.GetJWTSecret()
		assert.Equal(t, "test-secret", secret)
	})

	t.Run("GetEncryptionKey", func(t *testing.T) {
		key := cfg.GetEncryptionKey()
		assert.Equal(t, "test-key", key)
	})

	t.Run("GetJWTExpireDuration", func(t *testing.T) {
		duration := cfg.GetJWTExpireDuration()
		assert.Equal(t, 24*time.Hour, duration)
	})

	t.Run("GetRefreshTokenExpireDuration", func(t *testing.T) {
		duration := cfg.GetRefreshTokenExpireDuration()
		assert.Equal(t, 30*24*time.Hour, duration)
	})

	t.Run("GetLoggingLevel", func(t *testing.T) {
		level := cfg.GetLoggingLevel()
		assert.Equal(t, "info", level)
	})

	t.Run("GetLoggingFormat", func(t *testing.T) {
		format := cfg.GetLoggingFormat()
		assert.Equal(t, "json", format)
	})

	t.Run("GetLoggingOutput", func(t *testing.T) {
		output := cfg.GetLoggingOutput()
		assert.Equal(t, "stdout", output)
	})
}

func TestConfig_DefaultValues(t *testing.T) {
	cfg := &config{
		server: serverConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		database: databaseConfig{
			DSN: "postgres://vaultfactory:password@localhost:5432/vaultfactory?sslmode=disable",
		},
		security: securityConfig{
			JWTSecret:                  "default-secret",
			EncryptionKey:              "default-key",
			JWTExpireDuration:          24 * time.Hour,
			RefreshTokenExpireDuration: 30 * 24 * time.Hour,
		},
		logging: loggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}

	t.Run("GetDSN returns default value", func(t *testing.T) {
		dsn := cfg.GetDSN()
		assert.Equal(t, "postgres://vaultfactory:password@localhost:5432/vaultfactory?sslmode=disable", dsn)
	})

	t.Run("GetServerAddr returns formatted address", func(t *testing.T) {
		addr := cfg.GetServerAddr()
		assert.Equal(t, "0.0.0.0:8080", addr)
	})

	t.Run("GetJWTSecret returns default value", func(t *testing.T) {
		secret := cfg.GetJWTSecret()
		assert.Equal(t, "default-secret", secret)
	})

	t.Run("GetEncryptionKey returns default value", func(t *testing.T) {
		key := cfg.GetEncryptionKey()
		assert.Equal(t, "default-key", key)
	})
}
