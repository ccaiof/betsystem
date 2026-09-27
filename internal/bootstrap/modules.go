package bootstrap

import (
	httpapi "github.com/betsystem/internal/adapters/http"
	"github.com/betsystem/internal/infrastructure/keycloak"
	"github.com/betsystem/internal/infrastructure/postgres"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	postgres.Module,
	httpapi.Module,
	keycloak.Module,
)
