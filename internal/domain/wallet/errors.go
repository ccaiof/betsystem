package wallet

import "errors"

var (
	ErrInvalidPlayerID   = errors.New("wallet: invalid player id")
	ErrInsufficientFunds = errors.New("wallet: insufficient funds")
)
