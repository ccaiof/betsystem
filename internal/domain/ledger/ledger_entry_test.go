package ledger

import (
	"testing"
	"time"
	"uuid"

	"github.com/betsystem/internal/domain/money"
	"github.com/stretchr/testify/require"
)

func moneyBRL(t *testing.T, amount string) money.Money {
	t.Helper()

	value, err := money.FromDecimalString(amount, "BRL", false)
	require.NoError(t, err)

	return value
}

func TestNewEntry_Debit(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "30.00")
	after := moneyBRL(t, "70.00")

	entry, err := NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		DirectionDebit,
		amount,
		before,
		after,
		time.Now(),
	)

	require.NoError(t, err)

	require.Equal(t, DirectionDebit, entry.Direction)
	require.Equal(t, amount, entry.Amount)
	require.Equal(t, before, entry.BalanceBefore)
	require.Equal(t, after, entry.BalanceAfter)
}

func TestNewEntry_Credit(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "30.00")
	after := moneyBRL(t, "130.00")

	entry, err := NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		DirectionCredit,
		amount,
		before,
		after,
		time.Now(),
	)

	require.NoError(t, err)

	require.Equal(t, DirectionCredit, entry.Direction)
	require.Equal(t, amount, entry.Amount)
	require.Equal(t, before, entry.BalanceBefore)
	require.Equal(t, after, entry.BalanceAfter)
}

func TestNewEntry_RejectsInvalidDebitBalance(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "30.00")
	after := moneyBRL(t, "80.00")

	_, err := NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		DirectionDebit,
		amount,
		before,
		after,
		time.Now(),
	)

	require.ErrorIs(t, err, ErrInvalidBalanceFlow)
}

func TestNewEntry_RejectsInvalidCreditBalance(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "30.00")
	after := moneyBRL(t, "120.00")

	_, err := NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		DirectionCredit,
		amount,
		before,
		after,
		time.Now(),
	)

	require.ErrorIs(t, err, ErrInvalidBalanceFlow)
}

func TestNewEntry_RejectsCurrencyMismatch(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "30.00")

	after, err := money.FromDecimalString("70.00", "USD", false)
	require.NoError(t, err)

	_, err = NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		DirectionDebit,
		amount,
		before,
		after,
		time.Now(),
	)

	require.ErrorIs(t, err, ErrCurrencyMismatch)
}

func TestNewEntry_RejectsZeroAmount(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "0.00")
	after := moneyBRL(t, "100.00")

	_, err := NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		DirectionCredit,
		amount,
		before,
		after,
		time.Now(),
	)

	require.ErrorIs(t, err, ErrZeroAmount)
}

func TestNewEntry_RejectsInvalidDirection(t *testing.T) {
	before := moneyBRL(t, "100.00")
	amount := moneyBRL(t, "30.00")
	after := moneyBRL(t, "130.00")

	_, err := NewEntry(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		Direction("INVALID"),
		amount,
		before,
		after,
		time.Now(),
	)

	require.ErrorIs(t, err, ErrInvalidDirection)
}
