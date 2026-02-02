package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// User представляет пользователя системы.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Email        string    `json:"email" bun:"email,unique,notnull"`
	PasswordHash string    `json:"-" bun:"password_hash,notnull"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at,default:now()"`
	UpdatedAt    time.Time `json:"updated_at" bun:"updated_at,default:now()"`
}

// UserSession представляет сессию пользователя.
type UserSession struct {
	bun.BaseModel `bun:"table:user_sessions"`

	ID           uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" bun:"user_id,type:uuid,notnull"`
	RefreshToken string    `json:"-" bun:"refresh_token,unique,notnull"`
	ExpiresAt    time.Time `json:"expires_at" bun:"expires_at,notnull"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at,default:now()"`
	UpdatedAt    time.Time `json:"updated_at" bun:"updated_at,default:now()"`

	User *User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}
