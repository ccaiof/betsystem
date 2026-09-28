package mappers

import (
	domainwager "github.com/betsystem/internal/domain/wager"
	"github.com/betsystem/internal/infrastructure/postgres/models"
)

func WagerTransactionToModel(
	transaction domainwager.WagerTransaction,
) models.WagerTransactionModel {
	var financialResultBalanceMinor *int64

	if transaction.FinancialResultBalance != nil {
		value := transaction.FinancialResultBalance.AmountMinor
		financialResultBalanceMinor = &value
	}

	return models.WagerTransactionModel{
		ID:                             transaction.ID,
		ProviderID:                     transaction.ProviderID,
		ExternalTransactionID:          transaction.ExternalTransactionID,
		IdempotencyKey:                 transaction.IdempotencyKey,
		PayloadHash:                    transaction.PayloadHash,
		WalletID:                       transaction.WalletID,
		PlayerID:                       transaction.PlayerID,
		RoundID:                        transaction.RoundID,
		GameID:                         transaction.GameID,
		Kind:                           string(transaction.Kind),
		AmountMinor:                    transaction.Amount.AmountMinor,
		Currency:                       transaction.Amount.Currency,
		ReferenceExternalTransactionID: transaction.ReferenceExternalID,
		ReferenceTransactionID:         transaction.ReferenceTransactionID,
		Status:                         string(transaction.Status),
		FailureCode:                    transaction.FailureCode,
		RetryCount:                     transaction.RetryCount,
		NextRetryAt:                    transaction.NextRetryAt,
		FinancialResultBalanceMinor:    financialResultBalanceMinor,
		CreatedAt:                      transaction.CreatedAt,
		UpdatedAt:                      transaction.UpdatedAt,
		ProcessedAt:                    transaction.ProcessedAt,
	}
}
