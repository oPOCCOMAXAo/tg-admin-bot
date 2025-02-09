package db

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/migrations"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(newModule),
	)
}

type moduleParams struct {
	fx.In
	fx.Lifecycle

	Config Config
}

type moduleResults struct {
	fx.Out

	DB *gorm.DB
}

func newModule(
	params moduleParams,
) (moduleResults, error) {
	var res moduleResults

	var err error

	res.DB, err = NewMySQL(params.Config)
	if err != nil {
		return res, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return migrations.Migrate(ctx, res.DB)
		},
	})

	return res, nil
}
