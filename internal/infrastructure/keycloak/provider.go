package keycloak

import (
	"context"

	"github.com/coreos/go-oidc"
)

func NewProvider(config Config) (*oidc.Provider, error) {
	return oidc.NewProvider(
		context.Background(),
		config.Issuer,
	)
}
