package ledger

import "github.com/betsystem/internal/domain/money"

func validateDirection(direction Direction) error {
	switch direction {
	case DirectionDebit, DirectionCredit:
		return nil
	default:
		return ErrInvalidDirection
	}
}

func validateCurrencies(
	amount money.Money,
	balanceBefore money.Money,
	balanceAfter money.Money,
) error {
	if amount.Currency != balanceBefore.Currency ||
		amount.Currency != balanceAfter.Currency {
		return ErrCurrencyMismatch
	}

	return nil
}

func validateBalanceFlow(
	direction Direction,
	amount money.Money,
	balanceBefore money.Money,
	balanceAfter money.Money,
) error {
	var expected money.Money
	var err error

	switch direction {
	case DirectionDebit:
		expected, err = balanceBefore.Sub(amount)

	case DirectionCredit:
		expected, err = balanceBefore.Add(amount)

	default:
		return ErrInvalidDirection
	}

	if err != nil {
		return err
	}

	if expected.AmountMinor != balanceAfter.AmountMinor {
		return ErrInvalidBalanceFlow
	}

	return nil
}
