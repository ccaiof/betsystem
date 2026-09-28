package mappers

import (
	"github.com/betsystem/internal/domain/money"
	"github.com/betsystem/internal/domain/wallet"
	"github.com/betsystem/internal/infrastructure/postgres/models"
)

func WalletToModel(entity wallet.Wallet) models.WalletModel {
	return models.WalletModel{
		ID:           entity.ID,
		PlayerID:     entity.PlayerID,
		Currency:     entity.Balance.Currency,
		BalanceMinor: entity.Balance.AmountMinor,
		Version:      entity.Version,
	}
}

func WalletFromModel(
	model models.WalletModel,
) (wallet.Wallet, error) {
	balance, err := money.FromMinor(
		model.BalanceMinor,
		model.Currency,
	)
	if err != nil {
		return wallet.Wallet{}, err
	}

	return wallet.Wallet{
		ID:       model.ID,
		PlayerID: model.PlayerID,
		Balance:  balance,
		Version:  model.Version,
	}, nil
}
