package postgres

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/betsystem/internal/domain/money"
	"github.com/betsystem/internal/domain/wallet"
	"github.com/betsystem/internal/infrastructure/postgres/repositories"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestWalletRepository_SaveAndGetForUpdate(t *testing.T) {
	db := openTestDB(t)

	repository := repositories.NewWalletRepository(db)

	ctx := context.Background()

	walletID := uuid.New()
	playerID := uuid.New()

	entity, err := wallet.New(
		walletID,
		playerID,
		"BRL",
	)
	require.NoError(t, err)

	amount, err := money.FromDecimalString("100.00", "BRL", false)
	require.NoError(t, err)

	require.NoError(t, entity.Credit(amount))

	require.NoError(t, repository.Save(ctx, entity))

	loaded, err := repository.GetForUpdate(ctx, walletID)
	require.NoError(t, err)

	require.Equal(t, entity.ID, loaded.ID)
	require.Equal(t, entity.PlayerID, loaded.PlayerID)
	require.Equal(t, entity.Balance.Currency, loaded.Balance.Currency)
	require.Equal(t, entity.Balance.AmountMinor, loaded.Balance.AmountMinor)
	require.Equal(t, entity.Version, loaded.Version)
}

func TestWalletRepository_GetForUpdate_LocksWallet(t *testing.T) {
	db := openTestDB(t)

	repository := repositories.NewWalletRepository(db)

	walletID := uuid.New()
	playerID := uuid.New()

	entity, err := wallet.New(walletID, playerID, "BRL")
	require.NoError(t, err)

	initialAmount, err := money.FromDecimalString("100.00", "BRL", false)
	require.NoError(t, err)

	require.NoError(t, entity.Credit(initialAmount))
	require.NoError(t, repository.Save(context.Background(), entity))

	firstLocked := make(chan struct{})
	secondAttempting := make(chan struct{})
	releaseFirst := make(chan struct{})

	firstDone := make(chan error, 1)
	secondDone := make(chan error, 1)

	go func() {
		err := db.Transaction(func(tx *gorm.DB) error {
			repo := repositories.NewWalletRepository(tx)

			loaded, err := repo.GetForUpdate(context.Background(), walletID)
			if err != nil {
				return err
			}

			close(firstLocked)

			amount, err := money.FromDecimalString("30.00", "BRL", false)
			if err != nil {
				return err
			}

			require.NoError(t, loaded.Debit(amount))
			require.NoError(t, repo.Save(context.Background(), loaded))

			<-releaseFirst

			return nil
		})

		firstDone <- err
	}()

	<-firstLocked

	go func() {
		close(secondAttempting)

		err := db.Transaction(func(tx *gorm.DB) error {
			repo := repositories.NewWalletRepository(tx)

			loaded, err := repo.GetForUpdate(context.Background(), walletID)
			if err != nil {
				return err
			}

			require.Equal(t, int64(7000), loaded.Balance.AmountMinor)

			return nil
		})

		secondDone <- err
	}()

	<-secondAttempting

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-secondDone:
		require.NoError(t, err)
		t.Fatal("second transaction should still be waiting for the row lock")
	default:
	}

	close(releaseFirst)

	require.NoError(t, <-firstDone)
	require.NoError(t, <-secondDone)
}
