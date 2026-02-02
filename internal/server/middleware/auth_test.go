package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// MockAuthService для тестирования middleware
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

func (m *MockAuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, email, password)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceMockRecorder) Register(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), ctx, email, password)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

func (mr *MockAuthServiceMockRecorder) Login(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), ctx, email, password)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", ctx, refreshToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockAuthServiceMockRecorder) RefreshToken(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuthService)(nil).RefreshToken), ctx, refreshToken)
}

func (m *MockAuthService) Logout(ctx context.Context, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockAuthServiceMockRecorder) Logout(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthService)(nil).Logout), ctx, refreshToken)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, token)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceMockRecorder) ValidateToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuthService)(nil).ValidateToken), ctx, token)
}

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		middleware := NewAuthMiddleware(mockAuthService)

		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		mockAuthService.EXPECT().
			ValidateToken(gomock.Any(), "valid-token").
			Return(user, nil)

		handler := middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userFromCtx := GetUserFromContext(r.Context())
			assert.Equal(t, user, userFromCtx)
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		middleware := NewAuthMiddleware(mockAuthService)

		handler := middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called")
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header required")
	})

	t.Run("invalid authorization header format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		middleware := NewAuthMiddleware(mockAuthService)

		handler := middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called")
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		middleware := NewAuthMiddleware(mockAuthService)

		mockAuthService.EXPECT().
			ValidateToken(gomock.Any(), "invalid-token").
			Return(nil, assert.AnError)

		handler := middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called")
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("user in context", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		ctx := context.WithValue(context.Background(), UserKey, user)
		result := GetUserFromContext(ctx)

		assert.Equal(t, user, result)
	})

	t.Run("no user in context", func(t *testing.T) {
		ctx := context.Background()
		result := GetUserFromContext(ctx)

		assert.Nil(t, result)
	})
}
