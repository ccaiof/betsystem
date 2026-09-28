package repositories

import (
	"context"

	domainwager "github.com/betsystem/internal/domain/wager"
	"github.com/betsystem/internal/infrastructure/postgres/mappers"
	"github.com/betsystem/internal/infrastructure/postgres/txmanager"
	"gorm.io/gorm"
)

type WagerTransactionRepository struct {
	db *gorm.DB
}

func NewWagerTransactionRepository(
	db *gorm.DB,
) *WagerTransactionRepository {
	return &WagerTransactionRepository{
		db: db,
	}
}

func (r *WagerTransactionRepository) Save(
	ctx context.Context,
	transaction domainwager.WagerTransaction,
) error {
	db := txmanager.DBFromContext(ctx, r.db)

	model := mappers.WagerTransactionToModel(transaction)

	return db.Create(&model).Error
}
