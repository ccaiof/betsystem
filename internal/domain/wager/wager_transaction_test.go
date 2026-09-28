package wager

import (
	"testing"
	"time"
	"uuid"

	"github.com/betsystem/internal/domain/money"
	"github.com/stretchr/testify/require"
)

func newExternalWager(t *testing.T) WagerTransaction {
	t.Helper()

	amount, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	wager, err := NewExternal(
		uuid.New(),
		"provider-acme",
		"external-123",
		"idempotency-123",
		"hash-123",
		uuid.New(),
		uuid.New(),
		"round-123",
		"game-123",
		KindBet,
		amount,
		time.Now(),
	)

	require.NoError(t, err)

	return wager
}

func TestNewExternal_StartsPending(t *testing.T) {
	wager := newExternalWager(t)

	require.Equal(t, StatusPending, wager.Status)
	require.Equal(t, KindBet, wager.Kind)
}

func TestNewExternal_RejectsOpening(t *testing.T) {
	amount, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	_, err = NewExternal(
		uuid.New(),
		"provider-acme",
		"external-123",
		"idempotency-123",
		"hash-123",
		uuid.New(),
		uuid.New(),
		"round-123",
		"game-123",
		KindOpening,
		amount,
		time.Now(),
	)

	require.ErrorIs(t, err, ErrInvalidKind)
}

func TestNewOpening(t *testing.T) {
	amount, err := money.FromMinor(1000, "BRL")
	require.NoError(t, err)

	walletID := uuid.New()
	playerID := uuid.New()

	wager, err := NewOpening(
		uuid.New(),
		walletID,
		playerID,
		amount,
		time.Now(),
	)

	require.NoError(t, err)

	require.Equal(t, KindOpening, wager.Kind)
	require.Equal(t, StatusPending, wager.Status)
	require.Equal(t, walletID, wager.WalletID)
	require.Equal(t, playerID, wager.PlayerID)
	require.Equal(t, amount, wager.Amount)

	require.Nil(t, wager.ProviderID)
	require.Nil(t, wager.ExternalTransactionID)
	require.Nil(t, wager.IdempotencyKey)
	require.Nil(t, wager.PayloadHash)
	require.Nil(t, wager.RoundID)
	require.Nil(t, wager.GameID)
	require.Nil(t, wager.ReferenceExternalID)
	require.Nil(t, wager.ReferenceTransactionID)
}

func TestWagerTransaction_MarkProcessed(t *testing.T) {
	wager := newExternalWager(t)
	now := time.Now()

	err := wager.MarkProcessed(now)

	require.NoError(t, err)
	require.Equal(t, StatusProcessed, wager.Status)
	require.Equal(t, now, *wager.ProcessedAt)
	require.Equal(t, now, wager.UpdatedAt)
}

func TestWagerTransaction_MarkPendingReference(t *testing.T) {
	wager := newExternalWager(t)
	now := time.Now()

	err := wager.MarkPendingReference(now)

	require.NoError(t, err)
	require.Equal(t, StatusPendingReference, wager.Status)
	require.Equal(t, now, wager.UpdatedAt)
}

func TestWagerTransaction_Reject(t *testing.T) {
	wager := newExternalWager(t)
	now := time.Now()

	err := wager.Reject(
		"INSUFFICIENT_FUNDS",
		now,
	)

	require.NoError(t, err)
	require.Equal(t, StatusRejected, wager.Status)
	require.NotNil(t, wager.FailureCode)
	require.Equal(t, "INSUFFICIENT_FUNDS", *wager.FailureCode)
}

func TestWagerTransaction_Fail(t *testing.T) {
	wager := newExternalWager(t)
	now := time.Now()

	err := wager.Fail(
		"PERSISTENCE_FAILURE",
		now,
	)

	require.NoError(t, err)
	require.Equal(t, StatusFailed, wager.Status)
	require.NotNil(t, wager.FailureCode)
	require.Equal(t, "PERSISTENCE_FAILURE", *wager.FailureCode)
}

func TestWagerTransaction_TerminalStateCannotTransition(t *testing.T) {
	t.Run("processed", func(t *testing.T) {
		wager := newExternalWager(t)

		require.NoError(
			t,
			wager.MarkProcessed(time.Now()),
		)

		err := wager.Reject(
			"TEST",
			time.Now(),
		)

		require.ErrorIs(t, err, ErrTerminalState)
		require.Equal(t, StatusProcessed, wager.Status)
	})

	t.Run("rejected", func(t *testing.T) {
		wager := newExternalWager(t)

		require.NoError(
			t,
			wager.Reject("TEST", time.Now()),
		)

		err := wager.MarkProcessed(time.Now())

		require.ErrorIs(t, err, ErrTerminalState)
		require.Equal(t, StatusRejected, wager.Status)
	})

	t.Run("failed", func(t *testing.T) {
		wager := newExternalWager(t)

		require.NoError(
			t,
			wager.Fail("TEST", time.Now()),
		)

		err := wager.MarkProcessed(time.Now())

		require.ErrorIs(t, err, ErrTerminalState)
		require.Equal(t, StatusFailed, wager.Status)
	})
}

func TestWagerTransaction_RejectsInvalidTransition(t *testing.T) {
	wager := newExternalWager(t)

	require.NoError(
		t,
		wager.MarkPendingReference(time.Now()),
	)

	err := wager.transitionTo(StatusPending)

	require.ErrorIs(t, err, ErrInvalidTransition)
	require.Equal(t, StatusPendingReference, wager.Status)
}

func TestWagerTransaction_RejectRequiresFailureCode(t *testing.T) {
	wager := newExternalWager(t)

	err := wager.Reject("", time.Now())

	require.ErrorIs(t, err, ErrFailureCodeRequired)
	require.Equal(t, StatusPending, wager.Status)
}

func TestWagerTransaction_FailRequiresFailureCode(t *testing.T) {
	wager := newExternalWager(t)

	err := wager.Fail("", time.Now())

	require.ErrorIs(t, err, ErrFailureCodeRequired)
	require.Equal(t, StatusPending, wager.Status)
}
