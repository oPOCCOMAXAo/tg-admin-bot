package tg

import (
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/opoccomaxao/tg-instrumentation/router"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("tg",
		fx.Provide(newModule),
	)
}

type moduleParams struct {
	fx.In
	fx.Lifecycle

	Config Config
	Logger *slog.Logger
}

type moduleResult struct {
	fx.Out

	Service *Service
	Router  *router.Router
	Client  *bot.Bot
}

func newModule(
	params moduleParams,
) (moduleResult, error) {
	var (
		res moduleResult
		err error
	)

	res.Service, err = New(
		params.Config,
		params.Logger,
	)
	if err != nil {
		return res, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: res.Service.OnStart,
	})

	res.Router = res.Service.Router()
	res.Client = res.Service.Client()

	return res, nil
}
