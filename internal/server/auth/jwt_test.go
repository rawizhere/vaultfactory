package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

func TestJWTService_GenerateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", time.Hour)

	t.Run("successful token generation", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		token, err := jwtService.GenerateToken(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", time.Hour)

	t.Run("successful refresh token generation", func(t *testing.T) {
		token, err := jwtService.GenerateRefreshToken()

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Len(t, token, 64) // 32 bytes = 64 hex characters
	})
}

func TestJWTService_ValidateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", time.Hour)

	t.Run("valid token", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		token, err := jwtService.GenerateToken(user)
		assert.NoError(t, err)

		claims, err := jwtService.ValidateToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
	})

	t.Run("invalid token", func(t *testing.T) {
		claims, err := jwtService.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("expired token", func(t *testing.T) {
		// Создаем JWT сервис с очень коротким временем жизни токена
		shortJWTService := NewJWTService("test-secret", time.Millisecond)

		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		token, err := shortJWTService.GenerateToken(user)
		assert.NoError(t, err)

		// Ждем истечения токена
		time.Sleep(10 * time.Millisecond)

		claims, err := shortJWTService.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("wrong secret", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		// Генерируем токен с одним секретом
		token, err := jwtService.GenerateToken(user)
		assert.NoError(t, err)

		// Создаем JWT сервис с другим секретом
		wrongJWTService := NewJWTService("wrong-secret", time.Hour)

		claims, err := wrongJWTService.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
