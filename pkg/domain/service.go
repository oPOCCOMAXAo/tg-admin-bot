package domain

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/domain/repo"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"go.uber.org/fx"
)

type Service struct {
	repo  *repo.Repo
	check *CheckService
	cache *RulesCache
	tg    *tg.Service
}

func NewService(
	lc fx.Lifecycle,
	repo *repo.Repo,
	check *CheckService,
	tg *tg.Service,
) *Service {
	res := &Service{
		repo:  repo,
		check: check,
		tg:    tg,
		cache: NewRulesCache(),
	}

	lc.Append(fx.Hook{
		OnStart: res.OnStart,
	})

	return res
}

func (s *Service) OnStart(
	ctx context.Context,
) error {
	allConfigs, err := s.repo.GetAllChatConfigs(ctx)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	for _, config := range allConfigs {
		list := config.RulesList()
		if list.IsEmpty() {
			continue
		}

		s.cache.SetRules(config.TgID, list)
	}

	return nil
}
