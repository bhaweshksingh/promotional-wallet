package dto

import (
	accountInfo "account/pkg/accountinfo/model"
	"account/pkg/ledger/model"
	"github.com/google/uuid"
	"time"
)

const (
	Credit     = "Credit"
	Debit      = "Debit"
	Expiration = "Expiration"
)

func GetAllowedActivityTypes() map[string]bool {
	return map[string]bool{
		Credit:     true,
		Debit:      true,
		Expiration: true,
	}
}

type AccountEvent struct {
	UserID   string `json:"userID"`
	Amount   int64  `json:"amount"`
	Type     string `json:"type"`
	Priority int64  `json:"priority"`
	Expiry   int64  `json:"expiry"`
}

func (e *AccountEvent) GetCreditLedgerEntry(info accountInfo.AccountInfo) *model.Ledger {
	l := &model.Ledger{
		AccountID: info.ID,
		Amount:    e.Amount,
		Activity:  Credit,
		Type:      e.Type,
		Priority:  e.Priority,
		Expiry:    time.Unix(e.Expiry, 0),
	}
	return l
}

func (e *AccountEvent) GetAccountInfo(info accountInfo.AccountInfo, activity string) *accountInfo.AccountInfo {
	l := &accountInfo.AccountInfo{
		ID:      info.ID,
		UserID:  info.UserID,
		Balance: e.Amount,
	}

	if activity == Debit {
		l.Balance = l.Balance * -1
	}

	return l
}

func (e *AccountEvent) GetDebitLedgerEntry(info accountInfo.AccountInfo) *model.Ledger {
	l := &model.Ledger{
		AccountID: info.ID,
		Amount:    e.Amount,
		Activity:  Debit,
		Type:      e.Type,
	}
	return l
}

type AccountQuery struct {
	UserID uuid.UUID
	AccountID uuid.UUID
}

type LogQuery struct {
	ActivityType string
}
