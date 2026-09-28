package process_bet

import (
	"time"
	"uuid"
)

type ProcessBetCommand struct {
	TransactionID         uuid.UUID
	ProviderID            string
	ExternalTransactionID string
	IdempotencyKey        string
	PayloadHash           string

	WalletID uuid.UUID
	PlayerID uuid.UUID

	RoundID string
	GameID  string

	Amount   string
	Currency string

	Now time.Time
}
