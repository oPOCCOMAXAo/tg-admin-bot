package app

import (
	"github.com/opoccomaxao/tg-admin-bot/pkg/config"
	"github.com/opoccomaxao/tg-admin-bot/pkg/db"
	"github.com/opoccomaxao/tg-admin-bot/pkg/domain"
	"github.com/opoccomaxao/tg-admin-bot/pkg/endpoints"
	"github.com/opoccomaxao/tg-admin-bot/pkg/handlers"
	"github.com/opoccomaxao/tg-admin-bot/pkg/logger"
	"github.com/opoccomaxao/tg-admin-bot/pkg/server"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"go.uber.org/fx"
)

func Run() error {
	fx.New(
		fx.Provide(NewCancelCause),
		fx.WithLogger(NewFxLogger),
		config.Module(),
		logger.Module(),
		db.Module(),
		tg.Module(),
		server.Module(),
		domain.Module(),
		handlers.Invoke(),
		endpoints.Invoke(),
	).Run()

	return nil
}
