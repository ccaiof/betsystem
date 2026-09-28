package wallet

import (
	"uuid"

	"github.com/betsystem/internal/domain/money"
)

type Wallet struct {
	ID       uuid.UUID
	PlayerID uuid.UUID
	Balance  money.Money
	Version  int64
}

func New(id uuid.UUID, playerID uuid.UUID, currency string) (Wallet, error) {
	balance, err := money.Zero(currency)
	if err != nil {
		return Wallet{}, err
	}

	return Wallet{
		ID:       id,
		PlayerID: playerID,
		Balance:  balance,
		Version:  1,
	}, nil
}

func (w *Wallet) Credit(amount money.Money) error {
	if amount.IsNegative() {
		return money.ErrNegativeNotAllowed
	}

	if amount.IsZero() {
		return nil
	}

	balance, err := w.Balance.Add(amount)
	if err != nil {
		return err
	}

	w.Balance = balance
	w.Version++

	return nil
}

func (w *Wallet) Debit(amount money.Money) error {
	if amount.IsNegative() {
		return money.ErrNegativeNotAllowed
	}

	if amount.IsZero() {
		return nil
	}

	comparison, err := w.Balance.Compare(amount)
	if err != nil {
		return err
	}

	if comparison < 0 {
		return ErrInsufficientFunds
	}

	result, err := w.Balance.Sub(amount)
	if err != nil {
		return err
	}

	if result.IsNegative() {
		return ErrInsufficientFunds
	}

	w.Balance = result
	w.Version++

	return nil
}
