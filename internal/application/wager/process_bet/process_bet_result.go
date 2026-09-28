package process_bet

import (
	"uuid"

	"github.com/betsystem/internal/domain/money"
)

type ProcessBetResult struct {
	TransactionID uuid.UUID
	Status        string
	Balance       money.Money
}
