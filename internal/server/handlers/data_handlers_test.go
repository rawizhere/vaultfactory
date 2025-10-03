package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/server/middleware"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// MockDataService для тестирования handlers
type MockDataService struct {
	ctrl     *gomock.Controller
	recorder *MockDataServiceMockRecorder
}

type MockDataServiceMockRecorder struct {
	mock *MockDataService
}

func NewMockDataService(ctrl *gomock.Controller) *MockDataService {
	mock := &MockDataService{ctrl: ctrl}
	mock.recorder = &MockDataServiceMockRecorder{mock}
	return mock
}

func (m *MockDataService) EXPECT() *MockDataServiceMockRecorder {
	return m.recorder
}

func (m *MockDataService) CreateData(ctx context.Context, userID uuid.UUID, dataType models.DataType, name, metadata string, data []byte) (*models.DataItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateData", ctx, userID, dataType, name, metadata, data)
	ret0, _ := ret[0].(*models.DataItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDataServiceMockRecorder) CreateData(ctx, userID, dataType, name, metadata, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateData", reflect.TypeOf((*MockDataService)(nil).CreateData), ctx, userID, dataType, name, metadata, data)
}

func (m *MockDataService) GetData(ctx context.Context, userID, dataID uuid.UUID) (*models.DataItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData", ctx, userID, dataID)
	ret0, _ := ret[0].(*models.DataItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDataServiceMockRecorder) GetData(ctx, userID, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockDataService)(nil).GetData), ctx, userID, dataID)
}

func (m *MockDataService) GetUserData(ctx context.Context, userID uuid.UUID) ([]*models.DataItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", ctx, userID)
	ret0, _ := ret[0].([]*models.DataItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDataServiceMockRecorder) GetUserData(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockDataService)(nil).GetUserData), ctx, userID)
}

func (m *MockDataService) GetUserDataByType(ctx context.Context, userID uuid.UUID, dataType models.DataType) ([]*models.DataItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDataByType", ctx, userID, dataType)
	ret0, _ := ret[0].([]*models.DataItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDataServiceMockRecorder) GetUserDataByType(ctx, userID, dataType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDataByType", reflect.TypeOf((*MockDataService)(nil).GetUserDataByType), ctx, userID, dataType)
}

func (m *MockDataService) UpdateData(ctx context.Context, userID, dataID uuid.UUID, name, metadata string, data []byte) (*models.DataItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateData", ctx, userID, dataID, name, metadata, data)
	ret0, _ := ret[0].(*models.DataItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDataServiceMockRecorder) UpdateData(ctx, userID, dataID, name, metadata, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateData", reflect.TypeOf((*MockDataService)(nil).UpdateData), ctx, userID, dataID, name, metadata, data)
}

func (m *MockDataService) DeleteData(ctx context.Context, userID, dataID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteData", ctx, userID, dataID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockDataServiceMockRecorder) DeleteData(ctx, userID, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteData", reflect.TypeOf((*MockDataService)(nil).DeleteData), ctx, userID, dataID)
}

func (m *MockDataService) SyncData(ctx context.Context, userID uuid.UUID, since time.Time) ([]*models.DataItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncData", ctx, userID, since)
	ret0, _ := ret[0].([]*models.DataItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDataServiceMockRecorder) SyncData(ctx, userID, since interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncData", reflect.TypeOf((*MockDataService)(nil).SyncData), ctx, userID, since)
}

func TestDataHandler_CreateData(t *testing.T) {
	t.Run("successful data creation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		dataID := uuid.New()
		user := &models.User{ID: userID}

		createdItem := &models.DataItem{
			ID:        dataID,
			UserID:    userID,
			Type:      models.LoginPassword,
			Name:      "test-password",
			Metadata:  "test-metadata",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		}

		mockDataService.EXPECT().
			CreateData(gomock.Any(), userID, models.LoginPassword, "test-password", "test-metadata", gomock.Any()).
			Return(createdItem, nil)

		reqBody := CreateDataRequest{
			Type:     models.LoginPassword,
			Name:     "test-password",
			Metadata: "test-metadata",
			Data:     json.RawMessage(`{"username": "test", "password": "secret"}`),
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/data", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		w := httptest.NewRecorder()

		handler.CreateData(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DataResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, dataID.String(), response.ID)
		assert.Equal(t, models.LoginPassword, response.Type)
		assert.Equal(t, "test-password", response.Name)
	})

	t.Run("invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		user := &models.User{ID: userID}

		req := httptest.NewRequest("POST", "/data", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		w := httptest.NewRecorder()

		handler.CreateData(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("missing user in context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		reqBody := CreateDataRequest{
			Type:     models.LoginPassword,
			Name:     "test-password",
			Metadata: "test-metadata",
			Data:     json.RawMessage(`{"username": "test", "password": "secret"}`),
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/data", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Ожидаем панику, так как пользователь не найден в контексте
		assert.Panics(t, func() {
			handler.CreateData(w, req)
		})
	})
}

func TestDataHandler_GetData(t *testing.T) {
	t.Run("successful data retrieval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		dataID := uuid.New()
		user := &models.User{ID: userID}

		dataItem := &models.DataItem{
			ID:        dataID,
			UserID:    userID,
			Type:      models.LoginPassword,
			Name:      "test-password",
			Metadata:  "test-metadata",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		}

		mockDataService.EXPECT().
			GetData(gomock.Any(), userID, dataID).
			Return(dataItem, nil)

		req := httptest.NewRequest("GET", "/data/"+dataID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		req = mux.SetURLVars(req, map[string]string{"id": dataID.String()})
		w := httptest.NewRecorder()

		handler.GetData(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DataResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, dataID.String(), response.ID)
		assert.Equal(t, models.LoginPassword, response.Type)
		assert.Equal(t, "test-password", response.Name)
	})

	t.Run("data not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		dataID := uuid.New()
		user := &models.User{ID: userID}

		mockDataService.EXPECT().
			GetData(gomock.Any(), userID, dataID).
			Return(nil, assert.AnError)

		req := httptest.NewRequest("GET", "/data/"+dataID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		req = mux.SetURLVars(req, map[string]string{"id": dataID.String()})
		w := httptest.NewRecorder()

		handler.GetData(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "assert.AnError")
	})
}

func TestDataHandler_GetUserData(t *testing.T) {
	t.Run("successful user data retrieval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		user := &models.User{ID: userID}

		dataItems := []*models.DataItem{
			{
				ID:        uuid.New(),
				UserID:    userID,
				Type:      models.LoginPassword,
				Name:      "test-password",
				Metadata:  "test-metadata",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Version:   1,
			},
		}

		mockDataService.EXPECT().
			GetUserData(gomock.Any(), userID).
			Return(dataItems, nil)

		req := httptest.NewRequest("GET", "/data", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		w := httptest.NewRecorder()

		handler.GetUserData(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []DataResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, models.LoginPassword, response[0].Type)
		assert.Equal(t, "test-password", response[0].Name)
	})

	t.Run("no data found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		user := &models.User{ID: userID}

		mockDataService.EXPECT().
			GetUserData(gomock.Any(), userID).
			Return([]*models.DataItem{}, nil)

		req := httptest.NewRequest("GET", "/data", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		w := httptest.NewRecorder()

		handler.GetUserData(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []DataResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
	})
}

func TestDataHandler_DeleteData(t *testing.T) {
	t.Run("successful data deletion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		dataID := uuid.New()
		user := &models.User{ID: userID}

		mockDataService.EXPECT().
			DeleteData(gomock.Any(), userID, dataID).
			Return(nil)

		req := httptest.NewRequest("DELETE", "/data/"+dataID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		req = mux.SetURLVars(req, map[string]string{"id": dataID.String()})
		w := httptest.NewRecorder()

		handler.DeleteData(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("data not found for deletion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDataService := NewMockDataService(ctrl)
		handler := NewDataHandler(mockDataService)

		userID := uuid.New()
		dataID := uuid.New()
		user := &models.User{ID: userID}

		mockDataService.EXPECT().
			DeleteData(gomock.Any(), userID, dataID).
			Return(assert.AnError)

		req := httptest.NewRequest("DELETE", "/data/"+dataID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserKey, user))
		req = mux.SetURLVars(req, map[string]string{"id": dataID.String()})
		w := httptest.NewRecorder()

		handler.DeleteData(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "assert.AnError")
	})
}
