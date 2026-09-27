package bootstrap

import (
	"context"

	httpapi "github.com/betsystem/internal/adapters/http"
	"go.uber.org/fx"
)

func RegisterLifecyle(lifecycle fx.Lifecycle, server *httpapi.Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(ctx)
		},
	})
}
