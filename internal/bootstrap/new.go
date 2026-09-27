package bootstrap

import "go.uber.org/fx"

func New() *fx.App {
	return fx.New(
		Modules,

		fx.Invoke(
			RegisterLifecyle,
		),
	)
}
