package dto

import (
	"github.com/betsystem/internal/domain/money"
)

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func NewMoneyResponse(value money.Money) Money {
	return Money{
		Amount:   value.String(),
		Currency: value.Currency(),
	}
}

func (m Money) ToDomain() (money.Money, error) {
	return money.FromDecimalString(
		m.Amount,
		m.Currency,
		false,
	)
}
