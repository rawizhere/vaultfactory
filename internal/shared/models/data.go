package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DataType string

const (
	LoginPassword DataType = "login_password"
	TextData      DataType = "text_data"
	BinaryData    DataType = "binary_data"
	BankCard      DataType = "bank_card"
)

type DataItem struct {
	bun.BaseModel `bun:"table:data_items"`

	ID            uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserID        uuid.UUID `json:"user_id" bun:"user_id,type:uuid,notnull"`
	Type          DataType  `json:"type" bun:"type,notnull"`
	Name          string    `json:"name" bun:"name,notnull"`
	Metadata      string    `json:"metadata" bun:"metadata"`
	EncryptedData []byte    `json:"-" bun:"encrypted_data,notnull"`
	EncryptionKey []byte    `json:"-" bun:"encryption_key,notnull"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,default:now()"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,default:now()"`
	Version       int64     `json:"version" bun:"version,default:1"`

	User *User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}

type DataVersion struct {
	bun.BaseModel `bun:"table:data_versions"`

	ID        uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	DataID    uuid.UUID `json:"data_id" bun:"data_id,type:uuid,notnull"`
	Version   int64     `json:"version" bun:"version,notnull"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,default:now()"`

	DataItem *DataItem `json:"data_item,omitempty" bun:"rel:belongs-to,join:data_id=id"`
}


