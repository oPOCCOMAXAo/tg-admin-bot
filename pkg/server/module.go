package server

import (
	"context"
	"net/http"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("server",
		fx.Provide(
			newModule,
		),
	)
}

type moduleParams struct {
	fx.In
	fx.Lifecycle

	CancelCauseFunc context.CancelCauseFunc
	Config          Config
}

type moduleResult struct {
	fx.Out

	Server *Server
	Router *http.ServeMux
}

func newModule(params moduleParams) moduleResult {
	res := moduleResult{}

	res.Server = New(params.Config)

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return res.Server.OnStart(ctx, params.CancelCauseFunc)
		},
		OnStop: res.Server.OnStop,
	})

	res.Router = res.Server.router

	return res
}
