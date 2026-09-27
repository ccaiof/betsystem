package models

import (
	"time"
	"uuid"
)

type WagerTransactionModel struct {
	ID                             uuid.UUID `gorm:"primaryKey"`
	ProviderID                     *string
	ExternalTransactionID          *string
	IdempotencyKey                 *string
	PayloadHash                    *string
	WalletID                       uuid.UUID
	PlayerID                       uuid.UUID
	RoundID                        *string
	GameID                         *string
	Kind                           string
	AmountMinor                    int64
	Currency                       string
	ReferenceExternalTransactionID *string
	ReferenceTransactionID         *uuid.UUID
	Status                         string
	FailureCode                    *string
	RetryCount                     int
	NextRetryAt                    *time.Time
	FinancialResultBalanceMinor    *int64
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	ProcessedAt                    *time.Time
}

func (WagerTransactionModel) TableName() string {
	return "wager_transactions"
}
