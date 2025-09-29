package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*models.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
