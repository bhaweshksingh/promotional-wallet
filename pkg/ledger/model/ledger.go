package model

import (
	"github.com/google/uuid"
	"time"
)

type Ledger struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid" json:"account_id"`
	Amount    int64 `gorm:"column:amount;type:integer" json:"amount"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	Activity string `gorm:"type:string" json:"activity"`
	Type string `gorm:"type:string" json:"type"`
	Priority int64 `gorm:"type:integer" json:"priority"`
	Expiry time.Time `gorm:"column:expiry;default:now()" json:"expiry"`
}

func (Ledger) TableName() string {
	return "ledger"
}


type AggregateEntry struct {
	Activity string `gorm:"type:string" json:"activity"`
	Priority int64 `gorm:"type:integer" json:"priority"`
	Amount    int64 `gorm:"column:amount;type:integer" json:"amount"`
	Expiry time.Time `gorm:"column:expiry;default:now()" json:"expiry"`
	AccountID uuid.UUID `gorm:"type:uuid" json:"account_id"`
}