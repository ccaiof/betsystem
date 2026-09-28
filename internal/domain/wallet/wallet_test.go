package wallet

import (
	"testing"
	"uuid"

	"github.com/betsystem/internal/domain/money"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	playerID := uuid.New()
	walletID := uuid.New()

	wallet, err := New(
		walletID,
		playerID,
		"BRL",
	)

	require.NoError(t, err)

	require.Equal(t, walletID, wallet.ID)
	require.Equal(t, playerID, wallet.PlayerID)
	require.Equal(t, "BRL", wallet.Balance.Currency)
	require.Equal(t, int64(0), wallet.Balance.AmountMinor)
	require.Equal(t, int64(1), wallet.Version)
}

func TestWallet_Credit(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	amount, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	err = wallet.Credit(amount)

	require.NoError(t, err)
	require.Equal(t, int64(1000), wallet.Balance.AmountMinor)
	require.Equal(t, int64(2), wallet.Version)
}

func TestWallet_CreditIncrementsVersion(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	first, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	second, err := money.FromMinor(500, "BRL")
	require.NoError(t, err)

	require.NoError(t, wallet.Credit(first))
	require.Equal(t, int64(2), wallet.Version)

	require.NoError(t, wallet.Credit(second))
	require.Equal(t, int64(3), wallet.Version)

	require.Equal(
		t,
		int64(1500),
		wallet.Balance.AmountMinor,
	)
}

func TestWallet_Debit(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	credit, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	debit, err := money.FromMinor(300, "BRL")
	require.NoError(t, err)

	require.NoError(t, wallet.Credit(credit))
	require.NoError(t, wallet.Debit(debit))

	require.Equal(
		t,
		int64(700),
		wallet.Balance.AmountMinor,
	)

	require.Equal(t, int64(3), wallet.Version)
}

func TestWallet_DebitInsufficientFunds(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	credit, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	debit, err := money.FromMinor(1001, "BRL")
	require.NoError(t, err)

	require.NoError(t, wallet.Credit(credit))

	err = wallet.Debit(debit)

	require.ErrorIs(t, err, ErrInsufficientFunds)

	require.Equal(
		t,
		int64(1000),
		wallet.Balance.AmountMinor,
	)

	require.Equal(t, int64(2), wallet.Version)
}

func TestWallet_CreditRejectsDifferentCurrency(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	amount, err := money.FromMinor(1000, "USD")
	require.NoError(t, err)

	err = wallet.Credit(amount)

	require.ErrorIs(t, err, money.ErrCurrencyMismatch)

	require.Equal(
		t,
		int64(0),
		wallet.Balance.AmountMinor,
	)

	require.Equal(t, int64(1), wallet.Version)
}

func TestWallet_DebitRejectsDifferentCurrency(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	amount, err := money.FromMinor(1000, "USD")
	require.NoError(t, err)

	err = wallet.Debit(amount)

	require.ErrorIs(t, err, money.ErrCurrencyMismatch)

	require.Equal(
		t,
		int64(0),
		wallet.Balance.AmountMinor,
	)

	require.Equal(t, int64(1), wallet.Version)
}

func TestWallet_CreditRejectsNegativeAmount(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	amount, err := money.FromMinor(-100, "BRL")
	require.NoError(t, err)

	err = wallet.Credit(amount)

	require.ErrorIs(t, err, money.ErrNegativeNotAllowed)

	require.Equal(t, int64(0), wallet.Balance.AmountMinor)
	require.Equal(t, int64(1), wallet.Version)
}

func TestWallet_DebitRejectsNegativeAmount(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	amount, err := money.FromMinor(-100, "BRL")
	require.NoError(t, err)

	err = wallet.Debit(amount)

	require.ErrorIs(t, err, money.ErrNegativeNotAllowed)

	require.Equal(t, int64(0), wallet.Balance.AmountMinor)
	require.Equal(t, int64(1), wallet.Version)
}

func TestWallet_ZeroOperationDoesNotChangeVersion(t *testing.T) {
	wallet, err := New(
		uuid.New(),
		uuid.New(),
		"BRL",
	)
	require.NoError(t, err)

	zero, err := money.Zero("BRL")
	require.NoError(t, err)

	require.NoError(t, wallet.Credit(zero))
	require.Equal(t, int64(1), wallet.Version)

	require.NoError(t, wallet.Debit(zero))
	require.Equal(t, int64(1), wallet.Version)
}
