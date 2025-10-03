package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// MockAuthService для тестирования handlers
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

func TestAuthHandler_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		mockAuthService.EXPECT().
			Register(gomock.Any(), "test@example.com", "password123").
			Return(user, nil)

		mockAuthService.EXPECT().
			Login(gomock.Any(), "test@example.com", "password123").
			Return(user, "access-token", "refresh-token", nil)

		reqBody := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, response.User.(map[string]interface{})["email"])
		assert.Equal(t, "access-token", response.AccessToken)
	})

	t.Run("invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("registration error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.EXPECT().
			Register(gomock.Any(), "test@example.com", "password123").
			Return(nil, assert.AnError)

		reqBody := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "assert.AnError")
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		mockAuthService.EXPECT().
			Login(gomock.Any(), "test@example.com", "password123").
			Return(user, "access-token", "refresh-token", nil)

		reqBody := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, response.User.(map[string]interface{})["email"])
		assert.Equal(t, "access-token", response.AccessToken)
		assert.Equal(t, "refresh-token", response.RefreshToken)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.EXPECT().
			Login(gomock.Any(), "test@example.com", "wrongpassword").
			Return(nil, "", "", assert.AnError)

		reqBody := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "assert.AnError")
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	t.Run("successful token refresh", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.EXPECT().
			RefreshToken(gomock.Any(), "refresh-token").
			Return("new-access-token", "new-refresh-token", nil)

		reqBody := RefreshRequest{
			RefreshToken: "refresh-token",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Refresh(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "new-access-token", response.AccessToken)
		assert.Equal(t, "new-refresh-token", response.RefreshToken)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.EXPECT().
			RefreshToken(gomock.Any(), "invalid-token").
			Return("", "", assert.AnError)

		reqBody := RefreshRequest{
			RefreshToken: "invalid-token",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Refresh(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "assert.AnError")
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	t.Run("successful logout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.EXPECT().
			Logout(gomock.Any(), "refresh-token").
			Return(nil)

		reqBody := RefreshRequest{
			RefreshToken: "refresh-token",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("logout error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.EXPECT().
			Logout(gomock.Any(), "invalid-token").
			Return(assert.AnError)

		reqBody := RefreshRequest{
			RefreshToken: "invalid-token",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "assert.AnError")
	})
}
