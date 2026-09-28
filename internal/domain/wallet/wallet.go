package wallet

import (
	"uuid"

	"github.com/betsystem/internal/domain/money"
)

type Wallet struct {
	id       uuid.UUID
	playerID uuid.UUID
	balance  money.Money
	version  int64
}

func New(id uuid.UUID, playerID uuid.UUID, currency string) (Wallet, error) {
	balance, err := money.Zero(currency)
	if err != nil {
		return Wallet{}, err
	}

	return Wallet{
		id:       id,
		playerID: playerID,
		balance:  balance,
		version:  0,
	}, nil
}

func (w Wallet) ID() uuid.UUID {
	return w.id
}

func (w Wallet) PlayerID() uuid.UUID {
	return w.playerID
}

func (w Wallet) Balance() money.Money {
	return w.balance
}

func (w Wallet) Currency() string {
	return w.balance.Currency()
}

func (w Wallet) Version() int64 {
	return w.version
}

func (w *Wallet) Credit(amount money.Money) error {
	if amount.IsNegative() {
		return money.ErrNegativeNotAllowed
	}

	if amount.IsZero() {
		return nil
	}

	balance, err := w.balance.Add(amount)
	if err != nil {
		return err
	}

	w.balance = balance
	w.version++

	return nil
}

func (w *Wallet) Debit(amount money.Money) error {
	if amount.IsNegative() {
		return money.ErrNegativeNotAllowed
	}

	if amount.IsZero() {
		return nil
	}

	comparison, err := w.balance.Compare(amount)
	if err != nil {
		return err
	}

	if comparison < 0 {
		return ErrInsufficientFunds
	}

	result, err := w.balance.Sub(amount)
	if err != nil {
		return err
	}

	if result.IsNegative() {
		return ErrInsufficientFunds
	}

	w.balance = result
	w.version++

	return nil
}
