package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/server/auth"
	"github.com/tempizhere/vaultfactory/internal/server/service/mocks"
	"github.com/tempizhere/vaultfactory/internal/shared/crypto"
	"github.com/tempizhere/vaultfactory/internal/shared/logger"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

func TestAuthService_Register(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService("test-secret", time.Hour)

	t.Run("successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, errors.New("user not found"))

		mockUserRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *models.User) error {
				user.ID = uuid.New()
				return nil
			})

		user, err := authService.Register(ctx, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("user already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		email := "existing@example.com"
		password := "password123"

		existingUser := &models.User{
			ID:    uuid.New(),
			Email: email,
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(existingUser, nil)

		user, err := authService.Register(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("user creation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		email := "test2@example.com"
		password := "password123"

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, errors.New("user not found"))

		mockUserRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("database error"))

		user, err := authService.Register(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to create user")
	})
}

func TestAuthService_Login(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService("test-secret", time.Hour)

	t.Run("successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		hashedPassword, _ := cryptoService.HashPassword(password)
		user := &models.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: hashedPassword,
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		mockSessionRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(nil)

		returnedUser, accessToken, refreshToken, err := authService.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, returnedUser)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.Equal(t, user.ID, returnedUser.ID)
	})

	t.Run("invalid email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		email := "nonexistent@example.com"
		password := "password123"

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, errors.New("user not found"))

		user, accessToken, refreshToken, err := authService.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	t.Run("invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		email := "test@example.com"
		password := "wrongpassword"

		hashedPassword, _ := cryptoService.HashPassword("correctpassword")
		user := &models.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: hashedPassword,
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		returnedUser, accessToken, refreshToken, err := authService.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, returnedUser)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Contains(t, err.Error(), "invalid credentials")
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService("test-secret", time.Hour)

	t.Run("successful token refresh", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		refreshToken := "valid-refresh-token"
		userID := uuid.New()

		session := &models.UserSession{
			ID:           uuid.New(),
			UserID:       userID,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}

		user := &models.User{
			ID:    userID,
			Email: "test@example.com",
		}

		mockSessionRepo.EXPECT().
			GetByRefreshToken(ctx, refreshToken).
			Return(session, nil)

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockSessionRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(nil)

		newAccessToken, newRefreshToken, err := authService.RefreshToken(ctx, refreshToken)

		assert.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		assert.NotEqual(t, refreshToken, newRefreshToken)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		refreshToken := "invalid-token"

		mockSessionRepo.EXPECT().
			GetByRefreshToken(ctx, refreshToken).
			Return(nil, errors.New("session not found"))

		accessToken, newRefreshToken, err := authService.RefreshToken(ctx, refreshToken)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, newRefreshToken)
		assert.Contains(t, err.Error(), "invalid refresh token")
	})

	t.Run("expired refresh token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		refreshToken := "expired-token"

		session := &models.UserSession{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(-24 * time.Hour),
		}

		mockSessionRepo.EXPECT().
			GetByRefreshToken(ctx, refreshToken).
			Return(session, nil)

		accessToken, newRefreshToken, err := authService.RefreshToken(ctx, refreshToken)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, newRefreshToken)
		assert.Contains(t, err.Error(), "refresh token expired")
	})
}

func TestAuthService_Logout(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService("test-secret", time.Hour)

	t.Run("successful logout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		refreshToken := "valid-refresh-token"

		session := &models.UserSession{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			RefreshToken: refreshToken,
		}

		mockSessionRepo.EXPECT().
			GetByRefreshToken(ctx, refreshToken).
			Return(session, nil)

		mockSessionRepo.EXPECT().
			Delete(ctx, session.ID).
			Return(nil)

		err := authService.Logout(ctx, refreshToken)

		assert.NoError(t, err)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		refreshToken := "invalid-token"

		mockSessionRepo.EXPECT().
			GetByRefreshToken(ctx, refreshToken).
			Return(nil, errors.New("session not found"))

		err := authService.Logout(ctx, refreshToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid refresh token")
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService("test-secret", time.Hour)

	t.Run("valid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		userID := uuid.New()
		user := &models.User{
			ID:    userID,
			Email: "test@example.com",
		}

		accessToken, _ := jwtService.GenerateToken(user)

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		returnedUser, err := authService.ValidateToken(ctx, accessToken)

		assert.NoError(t, err)
		assert.NotNil(t, returnedUser)
		assert.Equal(t, user.ID, returnedUser.ID)
	})

	t.Run("invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		invalidToken := "invalid-token"

		user, err := authService.ValidateToken(ctx, invalidToken)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
		authService := NewAuthService(mockUserRepo, mockSessionRepo, cryptoService, jwtService, logger.NewMockLogger())

		ctx := context.Background()
		userID := uuid.New()
		user := &models.User{
			ID:    userID,
			Email: "test@example.com",
		}

		accessToken, _ := jwtService.GenerateToken(user)

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, errors.New("user not found"))

		returnedUser, err := authService.ValidateToken(ctx, accessToken)

		assert.Error(t, err)
		assert.Nil(t, returnedUser)
		assert.Contains(t, err.Error(), "user not found")
	})
}
