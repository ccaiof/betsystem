package models

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type WalletLedgerEntryModel struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"primaryKey"`
	WalletID           uuid.UUID
	TransactionID      uuid.UUID
	Direction          string
	AmountMinor        int64
	BalanceBeforeMinor int64
	BalanceAfterMinor  int64
	CreatedAt          time.Time
}

func (WalletLedgerEntryModel) TableName() string {
	return "wallet_ledger_entries"
}
