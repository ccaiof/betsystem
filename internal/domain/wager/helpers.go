package wager

func validateExternalKind(kind Kind) error {
	switch kind {
	case KindBet,
		KindWin,
		KindLoss,
		KindRefund,
		KindRollback:
		return nil

	case KindOpening:
		return ErrInvalidKind

	default:
		return ErrInvalidKind
	}
}

func isValidStatus(status Status) bool {
	switch status {
	case StatusPending,
		StatusPendingReference,
		StatusProcessed,
		StatusRejected,
		StatusFailed:
		return true

	default:
		return false
	}
}

func canTransition(
	from Status,
	to Status,
) bool {
	switch from {
	case StatusPending:
		return to == StatusProcessed ||
			to == StatusPendingReference ||
			to == StatusRejected ||
			to == StatusFailed

	case StatusPendingReference:
		return to == StatusProcessed ||
			to == StatusRejected ||
			to == StatusFailed

	case StatusProcessed,
		StatusRejected,
		StatusFailed:
		return false

	default:
		return false
	}
}

func (w *WagerTransaction) transitionTo(
	target Status,
) error {
	if !isValidStatus(target) {
		return ErrInvalidStatus
	}

	if w.Status.IsTerminal() {
		return ErrTerminalState
	}

	if !canTransition(w.Status, target) {
		return ErrInvalidTransition
	}

	w.Status = target

	return nil
}
