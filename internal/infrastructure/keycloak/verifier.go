package keycloak

import "github.com/coreos/go-oidc"

func NewVerifier(
	config Config,
	provider *oidc.Provider,
) *oidc.IDTokenVerifier {

	return provider.Verifier(
		&oidc.Config{
			ClientID: config.ClientID,
		},
	)

}
