package keycloak

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
)

type Verifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewVerifier(
	config Config,
	ctx context.Context,
) (*Verifier, error) {

	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		return nil, fmt.Errorf("create oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return &Verifier{
		verifier: verifier,
	}, nil

}

func (v *Verifier) Verify(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	token, err := v.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	return token, nil
}
