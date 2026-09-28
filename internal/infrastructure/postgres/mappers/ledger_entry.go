package mappers

import (
	"github.com/betsystem/internal/domain/ledger"
	"github.com/betsystem/internal/infrastructure/postgres/models"
)

func LedgerEntryToModel(
	entry ledger.LedgerEntry,
) models.WalletLedgerEntryModel {
	return models.WalletLedgerEntryModel{
		ID:                 entry.ID,
		WalletID:           entry.WalletID,
		TransactionID:      entry.TransactionID,
		Direction:          string(entry.Direction),
		AmountMinor:        entry.Amount.AmountMinor,
		BalanceBeforeMinor: entry.BalanceBefore.AmountMinor,
		BalanceAfterMinor:  entry.BalanceAfter.AmountMinor,
		CreatedAt:          entry.CreatedAt,
	}
}
