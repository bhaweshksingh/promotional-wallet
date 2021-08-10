package model

import (
	"github.com/google/uuid"
	"time"
)

type AccountInfo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	UserID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"user_id"`
	Balance   int64     `gorm:"type:integer;" json:"balance"`
}

func (i AccountInfo) TableName() string {
	return "accounts"
}

func (i *AccountInfo) UpdateBalance(balance int64) {
	i.Balance += balance
}
