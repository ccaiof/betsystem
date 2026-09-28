package auth

import "context"

type Verifier interface {
	Verify(ctx context.Context, rawToken string) (Identity, error)
}
