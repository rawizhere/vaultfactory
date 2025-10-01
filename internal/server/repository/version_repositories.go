package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
	"github.com/uptrace/bun"
)

// versionRepository реализует интерфейс VersionRepository для работы с версиями данных.
type versionRepository struct {
	db *bun.DB
}

// NewVersionRepository создает новый экземпляр VersionRepository.
func NewVersionRepository(db *bun.DB) interfaces.VersionRepository {
	return &versionRepository{db: db}
}

// Create создает новую версию данных в базе данных.
func (r *versionRepository) Create(ctx context.Context, version *models.DataVersion) error {
	_, err := r.db.NewInsert().Model(version).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create data version: %w", err)
	}
	return nil
}

// GetByDataID получает все версии данных по ID элемента данных.
func (r *versionRepository) GetByDataID(ctx context.Context, dataID uuid.UUID) ([]*models.DataVersion, error) {
	var versions []*models.DataVersion
	err := r.db.NewSelect().
		Model(&versions).
		Where("data_id = ?", dataID).
		Order("version DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions by data id: %w", err)
	}
	return versions, nil
}

// GetLatestVersion получает последнюю версию данных по ID элемента данных.
func (r *versionRepository) GetLatestVersion(ctx context.Context, dataID uuid.UUID) (*models.DataVersion, error) {
	version := new(models.DataVersion)
	err := r.db.NewSelect().
		Model(version).
		Where("data_id = ?", dataID).
		Order("version DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}
	return version, nil
}
