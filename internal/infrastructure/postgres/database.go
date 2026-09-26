package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(config Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
}
