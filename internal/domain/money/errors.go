package money

import "errors"

var (
	ErrInvalidAmount      = errors.New("money: invalid amount")
	ErrInvalidCurrency    = errors.New("money: invalid currency")
	ErrCurrencyMismatch   = errors.New("money: currency mismatch")
	ErrNegativeNotAllowed = errors.New("money: negative amount not allowed")
	ErrOverflow           = errors.New("money: overflow")
)
