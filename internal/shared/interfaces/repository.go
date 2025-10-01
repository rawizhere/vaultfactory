// Package interfaces определяет интерфейсы для работы с репозиториями данных.
package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// UserRepository определяет интерфейс для работы с пользователями.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SessionRepository определяет интерфейс для работы с сессиями пользователей.
type SessionRepository interface {
	Create(ctx context.Context, session *models.UserSession) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.UserSession, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserSession, error)
	Update(ctx context.Context, session *models.UserSession) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// DataRepository определяет интерфейс для работы с данными пользователей.
type DataRepository interface {
	Create(ctx context.Context, data *models.DataItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.DataItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.DataItem, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, dataType models.DataType) ([]*models.DataItem, error)
	Update(ctx context.Context, data *models.DataItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUpdatedSince(ctx context.Context, userID uuid.UUID, since time.Time) ([]*models.DataItem, error)
}

// VersionRepository определяет интерфейс для работы с версиями данных.
type VersionRepository interface {
	Create(ctx context.Context, version *models.DataVersion) error
	GetByDataID(ctx context.Context, dataID uuid.UUID) ([]*models.DataVersion, error)
	GetLatestVersion(ctx context.Context, dataID uuid.UUID) (*models.DataVersion, error)
}
