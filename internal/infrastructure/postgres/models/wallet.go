package models

import (
	"time"
	"uuid"
)

type WalletModel struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	PlayerID     uuid.UUID
	Currency     string
	BalanceMinor int64
	Version      int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (WalletModel) TableName() string {
	return "wallets"
}
