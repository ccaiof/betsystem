package wager

import "errors"

var (
	ErrInvalidKind = errors.New("wager: invalid kind")

	ErrInvalidStatus = errors.New("wager: invalid status")

	ErrInvalidTransition = errors.New("wager: invalid status transition")

	ErrTerminalState = errors.New("wager: transaction is in a terminal state")

	ErrFailureCodeRequired = errors.New("wager: failure code is required")
)
