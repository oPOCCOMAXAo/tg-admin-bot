package domain

import (
	"github.com/opoccomaxao/tg-admin-bot/pkg/domain/repo"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("domain",
		fx.Provide(repo.NewRepo, fx.Private),
		fx.Provide(NewCalculatorService, fx.Private),
		fx.Provide(NewService),
	)
}
