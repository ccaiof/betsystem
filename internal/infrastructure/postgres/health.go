package postgres

import (
	"context"

	"gorm.io/gorm"
)

type HealthChecker struct {
	db *gorm.DB
}

func NewHealthChecker(db *gorm.DB) *HealthChecker {
	return &HealthChecker{
		db: db,
	}
}

func (h *HealthChecker) Check(ctx context.Context) error {
	sqlDb, err := h.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.PingContext(ctx)
}
