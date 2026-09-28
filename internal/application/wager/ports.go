package wager

import (
	"context"
	"uuid"

	"github.com/betsystem/internal/application/events"
	"github.com/betsystem/internal/domain/ledger"
	domainwager "github.com/betsystem/internal/domain/wager"
	"github.com/betsystem/internal/domain/wallet"
)

type WalletRepository interface {
	GetForUpdate(
		ctx context.Context,
		walletID uuid.UUID,
	) (wallet.Wallet, error)

	Save(
		ctx context.Context,
		wallet wallet.Wallet,
	) error
}

type WagerTransactionRepository interface {
	Save(
		ctx context.Context,
		transaction domainwager.WagerTransaction,
	) error
}

type LedgerRepository interface {
	Save(
		ctx context.Context,
		entry ledger.LedgerEntry,
	) error
}

type OutboxRepository interface {
	Save(ctx context.Context, event events.Event) error
}
