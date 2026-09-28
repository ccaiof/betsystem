package keycloak

import (
	"context"
	"fmt"

	"github.com/betsystem/internal/application/auth"
	"github.com/coreos/go-oidc"
)

type Verifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewVerifier(
	config *Config,
	provider *oidc.Provider,
) (*Verifier, error) {
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return &Verifier{
		verifier: verifier,
	}, nil

}

type tokenClaims struct {
	Subject         string `json:"sub"`
	AuthorizedParty string `json:"azp"`
	RealmAccess     struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func (v *Verifier) Verify(ctx context.Context, rawToken string) (auth.Identity, error) {
	token, err := v.verifier.Verify(ctx, rawToken)
	if err != nil {
		return auth.Identity{}, fmt.Errorf("verify token: %w", err)
	}

	var claims tokenClaims

	if err := token.Claims(&claims); err != nil {
		return auth.Identity{}, fmt.Errorf("decode token claims: %w", err)
	}

	return auth.Identity{
		Subject:    claims.Subject,
		ProviderID: claims.AuthorizedParty,
		Roles:      claims.RealmAccess.Roles,
	}, nil
}
