package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tempizhere/vaultfactory/internal/server/auth"
	"github.com/tempizhere/vaultfactory/internal/shared/crypto"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// authService реализует интерфейс AuthService для аутентификации пользователей.
type authService struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.SessionRepository
	crypto      *crypto.CryptoService
	jwt         *auth.JWTService
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(
	userRepo interfaces.UserRepository,
	sessionRepo interfaces.SessionRepository,
	crypto *crypto.CryptoService,
	jwt *auth.JWTService,
) interfaces.AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		crypto:      crypto,
		jwt:         jwt,
	}
}

// Register регистрирует нового пользователя.
func (s *authService) Register(ctx context.Context, email, password string) (*models.User, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	passwordHash, err := s.crypto.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login выполняет аутентификацию пользователя и возвращает токены.
func (s *authService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	if !s.crypto.VerifyPassword(password, user.PasswordHash) {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := &models.UserSession{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, "", "", fmt.Errorf("failed to create session: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// RefreshToken обновляет access токен с помощью refresh токена.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(session.ExpiresAt) {
		return "", "", fmt.Errorf("refresh token expired")
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	newAccessToken, err := s.jwt.GenerateToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return "", "", fmt.Errorf("failed to update session: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout завершает сессию пользователя.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	if err := s.sessionRepo.Delete(ctx, session.ID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// ValidateToken проверяет валидность access токена и возвращает пользователя.
func (s *authService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.jwt.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
