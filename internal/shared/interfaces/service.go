package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

// AuthService определяет интерфейс для аутентификации пользователей.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}

// DataService определяет интерфейс для работы с данными пользователей.
type DataService interface {
	CreateData(ctx context.Context, userID uuid.UUID, dataType models.DataType, name, metadata string, data []byte) (*models.DataItem, error)
	GetData(ctx context.Context, userID, dataID uuid.UUID) (*models.DataItem, error)
	GetUserData(ctx context.Context, userID uuid.UUID) ([]*models.DataItem, error)
	GetUserDataByType(ctx context.Context, userID uuid.UUID, dataType models.DataType) ([]*models.DataItem, error)
	UpdateData(ctx context.Context, userID, dataID uuid.UUID, name, metadata string, data []byte) (*models.DataItem, error)
	DeleteData(ctx context.Context, userID, dataID uuid.UUID) error
	SyncData(ctx context.Context, userID uuid.UUID, lastSync time.Time) ([]*models.DataItem, error)
}

// CryptoService определяет интерфейс для криптографических операций.
type CryptoService interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(encryptedData []byte, key []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
}
