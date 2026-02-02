// Package repository содержит реализации репозиториев для работы с базой данных.
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
	"github.com/uptrace/bun"
)

// dataRepository реализует интерфейс DataRepository для работы с данными пользователей.
type dataRepository struct {
	db *bun.DB
}

// NewDataRepository создает новый экземпляр DataRepository.
func NewDataRepository(db *bun.DB) interfaces.DataRepository {
	return &dataRepository{db: db}
}

// Create создает новый элемент данных в базе данных.
func (r *dataRepository) Create(ctx context.Context, data *models.DataItem) error {
	_, err := r.db.NewInsert().Model(data).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create data item: %w", err)
	}
	return nil
}

// GetByID получает элемент данных по ID.
func (r *dataRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.DataItem, error) {
	data := new(models.DataItem)
	err := r.db.NewSelect().
		Model(data).
		Relation("User").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get data item by id: %w", err)
	}
	return data, nil
}

// GetByUserID получает все данные пользователя.
func (r *dataRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.DataItem, error) {
	var items []*models.DataItem
	err := r.db.NewSelect().
		Model(&items).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get data items by user id: %w", err)
	}
	return items, nil
}

// GetByUserIDAndType получает данные пользователя определенного типа.
func (r *dataRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, dataType models.DataType) ([]*models.DataItem, error) {
	var items []*models.DataItem
	err := r.db.NewSelect().
		Model(&items).
		Where("user_id = ? AND type = ?", userID, dataType).
		Order("updated_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get data items by user id and type: %w", err)
	}
	return items, nil
}

// Update обновляет элемент данных в базе данных.
func (r *dataRepository) Update(ctx context.Context, data *models.DataItem) error {
	_, err := r.db.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update data item: %w", err)
	}
	return nil
}

// Delete удаляет элемент данных из базы данных.
func (r *dataRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*models.DataItem)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete data item: %w", err)
	}
	return nil
}

// GetUpdatedSince получает данные, измененные после указанного времени.
func (r *dataRepository) GetUpdatedSince(ctx context.Context, userID uuid.UUID, since time.Time) ([]*models.DataItem, error) {
	var items []*models.DataItem
	err := r.db.NewSelect().
		Model(&items).
		Where("user_id = ? AND updated_at > ?", userID, since).
		Order("updated_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated data items: %w", err)
	}
	return items, nil
}
