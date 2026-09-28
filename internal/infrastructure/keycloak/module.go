package keycloak

import (
	"github.com/betsystem/internal/application/auth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"keycloak",

	fx.Provide(
		NewConfig,
		NewProvider,
		NewVerifier,
		func(verifier *Verifier) auth.Verifier {
			return verifier
		},
	),
)
