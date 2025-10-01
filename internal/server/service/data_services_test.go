package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tempizhere/vaultfactory/internal/server/service/mocks"
	"github.com/tempizhere/vaultfactory/internal/shared/crypto"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

func TestDataService_CreateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	dataType := models.LoginPassword
	name := "test data"
	metadata := "test metadata"
	data := []byte("test data content")

	mockDataRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockVersionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	result, err := service.CreateData(ctx, userID, dataType, name, metadata, data)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, dataType, result.Type)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, metadata, result.Metadata)
	assert.Equal(t, int64(1), result.Version)
	assert.NotEmpty(t, result.EncryptedData)
	assert.NotEmpty(t, result.EncryptionKey)
}

func TestDataService_CreateData_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	dataType := models.LoginPassword
	name := "test data"
	metadata := "test metadata"
	data := []byte("test data content")

	mockDataRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("repository error"))

	result, err := service.CreateData(ctx, userID, dataType, name, metadata, data)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create data item")
}

func TestDataService_GetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	dataID := uuid.New()

	dataItem := &models.DataItem{
		ID:     dataID,
		UserID: userID,
		Type:   models.LoginPassword,
		Name:   "test data",
	}

	mockDataRepo.EXPECT().GetByID(ctx, dataID).Return(dataItem, nil)

	result, err := service.GetData(ctx, userID, dataID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dataID, result.ID)
	assert.Equal(t, userID, result.UserID)
}

func TestDataService_GetData_AccessDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	dataID := uuid.New()

	dataItem := &models.DataItem{
		ID:     dataID,
		UserID: otherUserID,
		Type:   models.LoginPassword,
		Name:   "test data",
	}

	mockDataRepo.EXPECT().GetByID(ctx, dataID).Return(dataItem, nil)

	result, err := service.GetData(ctx, userID, dataID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "access denied")
}

func TestDataService_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()

	dataItems := []*models.DataItem{
		{
			ID:            uuid.New(),
			UserID:        userID,
			Type:          models.LoginPassword,
			Name:          "test data 1",
			EncryptedData: []byte("encrypted1"),
			EncryptionKey: []byte("key1"),
		},
		{
			ID:            uuid.New(),
			UserID:        userID,
			Type:          models.TextData,
			Name:          "test data 2",
			EncryptedData: []byte("encrypted2"),
			EncryptionKey: []byte("key2"),
		},
	}

	mockDataRepo.EXPECT().GetByUserID(ctx, userID).Return(dataItems, nil)

	result, err := service.GetUserData(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	for _, item := range result {
		assert.Nil(t, item.EncryptedData)
		assert.Nil(t, item.EncryptionKey)
	}
}

func TestDataService_UpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	dataID := uuid.New()

	encryptionKey, _ := cryptoService.GenerateKey()
	encryptedData, _ := cryptoService.Encrypt([]byte("old data"), encryptionKey)

	dataItem := &models.DataItem{
		ID:            dataID,
		UserID:        userID,
		Type:          models.LoginPassword,
		Name:          "old name",
		Metadata:      "old metadata",
		EncryptedData: encryptedData,
		EncryptionKey: encryptionKey,
		Version:       1,
	}

	mockDataRepo.EXPECT().GetByID(ctx, dataID).Return(dataItem, nil)
	mockDataRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	mockVersionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	newName := "new name"
	newMetadata := "new metadata"
	newData := []byte("new data content")

	result, err := service.UpdateData(ctx, userID, dataID, newName, newMetadata, newData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newMetadata, result.Metadata)
	assert.Equal(t, int64(2), result.Version)
}

func TestDataService_UpdateData_AccessDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	dataID := uuid.New()

	dataItem := &models.DataItem{
		ID:     dataID,
		UserID: otherUserID,
		Type:   models.LoginPassword,
		Name:   "test data",
	}

	mockDataRepo.EXPECT().GetByID(ctx, dataID).Return(dataItem, nil)

	result, err := service.UpdateData(ctx, userID, dataID, "new name", "new metadata", []byte("new data"))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "access denied")
}

func TestDataService_DeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	dataID := uuid.New()

	dataItem := &models.DataItem{
		ID:     dataID,
		UserID: userID,
		Type:   models.LoginPassword,
		Name:   "test data",
	}

	mockDataRepo.EXPECT().GetByID(ctx, dataID).Return(dataItem, nil)
	mockDataRepo.EXPECT().Delete(ctx, dataID).Return(nil)

	err := service.DeleteData(ctx, userID, dataID)

	assert.NoError(t, err)
}

func TestDataService_DeleteData_AccessDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	dataID := uuid.New()

	dataItem := &models.DataItem{
		ID:     dataID,
		UserID: otherUserID,
		Type:   models.LoginPassword,
		Name:   "test data",
	}

	mockDataRepo.EXPECT().GetByID(ctx, dataID).Return(dataItem, nil)

	err := service.DeleteData(ctx, userID, dataID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestDataService_SyncData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataRepo := mocks.NewMockDataRepository(ctrl)
	mockVersionRepo := mocks.NewMockVersionRepository(ctrl)
	cryptoService := crypto.NewCryptoService()

	service := NewDataService(mockDataRepo, mockVersionRepo, cryptoService)

	ctx := context.Background()
	userID := uuid.New()
	lastSync := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	dataItems := []*models.DataItem{
		{
			ID:            uuid.New(),
			UserID:        userID,
			Type:          models.LoginPassword,
			Name:          "test data 1",
			EncryptedData: []byte("encrypted1"),
			EncryptionKey: []byte("key1"),
		},
	}

	mockDataRepo.EXPECT().GetUpdatedSince(ctx, userID, lastSync).Return(dataItems, nil)

	result, err := service.SyncData(ctx, userID, lastSync)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Nil(t, result[0].EncryptedData)
	assert.Nil(t, result[0].EncryptionKey)
}
