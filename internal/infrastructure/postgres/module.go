package postgres

import (
	"github.com/betsystem/internal/infrastructure/postgres/txmanager"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"postgres",
	fx.Provide(NewConfig, NewDatabase, txmanager.New),
)
