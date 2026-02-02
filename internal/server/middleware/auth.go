// Package middleware содержит HTTP middleware для сервера.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
)

type userKey string

const UserKey userKey = "user"

// GetUserFromContext извлекает пользователя из контекста.
func GetUserFromContext(ctx context.Context) interface{} {
	return ctx.Value(UserKey)
}

// AuthMiddleware предоставляет middleware для аутентификации.
type AuthMiddleware struct {
	authService interfaces.AuthService
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(authService interfaces.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth проверяет аутентификацию пользователя для защищенных маршрутов.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		user, err := m.authService.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
