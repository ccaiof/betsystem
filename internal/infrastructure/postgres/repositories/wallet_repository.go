package repositories

import (
	"context"
	"errors"
	"uuid"

	"github.com/betsystem/internal/domain/wallet"
	"github.com/betsystem/internal/infrastructure/postgres/mappers"
	"github.com/betsystem/internal/infrastructure/postgres/models"
	"github.com/betsystem/internal/infrastructure/postgres/txmanager"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrWalletNotFound = errors.New("wallet repository: wallet not found")

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) GetForUpdate(
	ctx context.Context,
	walletID uuid.UUID,
) (wallet.Wallet, error) {
	db := txmanager.DBFromContext(ctx, r.db)

	var model models.WalletModel

	err := db.
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		First(&model, "id = ?", walletID).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return wallet.Wallet{}, ErrWalletNotFound
	}

	if err != nil {
		return wallet.Wallet{}, err
	}

	return mappers.WalletFromModel(model)
}

func (r *WalletRepository) Save(
	ctx context.Context,
	entity wallet.Wallet,
) error {
	db := txmanager.DBFromContext(ctx, r.db)

	model := mappers.WalletToModel(entity)

	return db.Save(&model).Error
}
