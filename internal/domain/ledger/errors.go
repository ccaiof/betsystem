package ledger

import "errors"

var (
	ErrInvalidDirection   = errors.New("ledger: invalid direction")
	ErrNegativeAmount     = errors.New("ledger: negative amount")
	ErrZeroAmount         = errors.New("ledger: zero amount")
	ErrCurrencyMismatch   = errors.New("ledger: currency mismatch")
	ErrInvalidBalanceFlow = errors.New("ledger: invalid balance flow")
)
