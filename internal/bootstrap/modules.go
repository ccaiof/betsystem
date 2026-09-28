package bootstrap

import (
	"github.com/betsystem/internal/infrastructure/keycloak"
	"github.com/betsystem/internal/infrastructure/postgres"
	httpapi "github.com/betsystem/internal/presentation/http"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	postgres.Module,
	httpapi.Module,
	keycloak.Module,
)
