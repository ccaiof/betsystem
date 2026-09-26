package bootstrap

import (
	"github.com/betsystem/internal/infrastructure/postgres"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	postgres.Module,
)
