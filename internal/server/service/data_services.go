package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/crypto"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

type dataService struct {
	dataRepo    interfaces.DataRepository
	versionRepo interfaces.VersionRepository
	crypto      *crypto.CryptoService
}

func NewDataService(
	dataRepo interfaces.DataRepository,
	versionRepo interfaces.VersionRepository,
	crypto *crypto.CryptoService,
) interfaces.DataService {
	return &dataService{
		dataRepo:    dataRepo,
		versionRepo: versionRepo,
		crypto:      crypto,
	}
}

func (s *dataService) CreateData(ctx context.Context, userID uuid.UUID, dataType models.DataType, name, metadata string, data []byte) (*models.DataItem, error) {
	encryptionKey, err := s.crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	encryptedData, err := s.crypto.Encrypt(data, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	dataItem := &models.DataItem{
		UserID:        userID,
		Type:          dataType,
		Name:          name,
		Metadata:      metadata,
		EncryptedData: encryptedData,
		EncryptionKey: encryptionKey,
		Version:       1,
	}

	if err := s.dataRepo.Create(ctx, dataItem); err != nil {
		return nil, fmt.Errorf("failed to create data item: %w", err)
	}

	version := &models.DataVersion{
		DataID:  dataItem.ID,
		Version: 1,
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create data version: %w", err)
	}

	return dataItem, nil
}

func (s *dataService) GetData(ctx context.Context, userID, dataID uuid.UUID) (*models.DataItem, error) {
	dataItem, err := s.dataRepo.GetByID(ctx, dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data item: %w", err)
	}

	if dataItem.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return dataItem, nil
}

func (s *dataService) GetUserData(ctx context.Context, userID uuid.UUID) ([]*models.DataItem, error) {
	items, err := s.dataRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	for _, item := range items {
		item.EncryptedData = nil
		item.EncryptionKey = nil
	}

	return items, nil
}

func (s *dataService) GetUserDataByType(ctx context.Context, userID uuid.UUID, dataType models.DataType) ([]*models.DataItem, error) {
	items, err := s.dataRepo.GetByUserIDAndType(ctx, userID, dataType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data by type: %w", err)
	}

	for _, item := range items {
		item.EncryptedData = nil
		item.EncryptionKey = nil
	}

	return items, nil
}

func (s *dataService) UpdateData(ctx context.Context, userID, dataID uuid.UUID, name, metadata string, data []byte) (*models.DataItem, error) {
	dataItem, err := s.dataRepo.GetByID(ctx, dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data item: %w", err)
	}

	if dataItem.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	encryptedData, err := s.crypto.Encrypt(data, dataItem.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	dataItem.Name = name
	dataItem.Metadata = metadata
	dataItem.EncryptedData = encryptedData
	dataItem.Version++
	dataItem.UpdatedAt = time.Now()

	if err := s.dataRepo.Update(ctx, dataItem); err != nil {
		return nil, fmt.Errorf("failed to update data item: %w", err)
	}

	version := &models.DataVersion{
		DataID:  dataItem.ID,
		Version: dataItem.Version,
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create data version: %w", err)
	}

	return dataItem, nil
}

func (s *dataService) DeleteData(ctx context.Context, userID, dataID uuid.UUID) error {
	dataItem, err := s.dataRepo.GetByID(ctx, dataID)
	if err != nil {
		return fmt.Errorf("failed to get data item: %w", err)
	}

	if dataItem.UserID != userID {
		return fmt.Errorf("access denied")
	}

	if err := s.dataRepo.Delete(ctx, dataID); err != nil {
		return fmt.Errorf("failed to delete data item: %w", err)
	}

	return nil
}

func (s *dataService) SyncData(ctx context.Context, userID uuid.UUID, lastSync time.Time) ([]*models.DataItem, error) {
	items, err := s.dataRepo.GetUpdatedSince(ctx, userID, lastSync)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated data: %w", err)
	}

	for _, item := range items {
		item.EncryptedData = nil
		item.EncryptionKey = nil
	}

	return items, nil
}
