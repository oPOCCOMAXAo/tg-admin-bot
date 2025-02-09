package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/opoccomaxao/tg-admin-bot/pkg/db"
	"github.com/opoccomaxao/tg-admin-bot/pkg/logger"
	"github.com/opoccomaxao/tg-admin-bot/pkg/server"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("config",
		fx.Provide(New),
	)
}

type Config struct {
	fx.Out

	Logger logger.Config `envPrefix:"LOGGER_"`
	Server server.Config `envPrefix:"SERVER_"`
	DB     db.Config     `envPrefix:"DB_"`
	TG     tg.Config     `envPrefix:"TG_"`
}

func New() (Config, error) {
	var res Config

	err := env.ParseWithOptions(&res, env.Options{
		UseFieldNameByDefault: false,
		RequiredIfNoDef:       false,
	})
	if err != nil {
		return res, errors.WithStack(err)
	}

	return res, nil
}
