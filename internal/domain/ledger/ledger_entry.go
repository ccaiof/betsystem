package ledger

import (
	"time"
	"uuid"

	"github.com/betsystem/internal/domain/money"
)

type LedgerEntry struct {
	ID            uuid.UUID
	WalletID      uuid.UUID
	TransactionID uuid.UUID
	Direction     Direction
	Amount        money.Money
	BalanceBefore money.Money
	BalanceAfter  money.Money
	CreatedAt     time.Time
}

func NewEntry(
	id uuid.UUID,
	walletID uuid.UUID,
	transactionID uuid.UUID,
	direction Direction,
	amount money.Money,
	balanceBefore money.Money,
	balanceAfter money.Money,
	createdAt time.Time,
) (LedgerEntry, error) {
	if err := validateDirection(direction); err != nil {
		return LedgerEntry{}, err
	}

	if amount.IsNegative() {
		return LedgerEntry{}, ErrNegativeAmount
	}

	if amount.IsZero() {
		return LedgerEntry{}, ErrZeroAmount
	}

	if err := validateCurrencies(
		amount,
		balanceBefore,
		balanceAfter,
	); err != nil {
		return LedgerEntry{}, err
	}

	if err := validateBalanceFlow(
		direction,
		amount,
		balanceBefore,
		balanceAfter,
	); err != nil {
		return LedgerEntry{}, err
	}

	return LedgerEntry{
		ID:            id,
		WalletID:      walletID,
		TransactionID: transactionID,
		Direction:     direction,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     createdAt,
	}, nil
}
