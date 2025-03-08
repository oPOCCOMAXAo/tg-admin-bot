package domain

import (
	"context"
	"log/slog"
	"time"

	"github.com/opoccomaxao/tg-admin-bot/pkg/domain/repo"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"go.uber.org/fx"
)

type Service struct {
	repo       *repo.Repo
	calculator *CalculatorService
	cache      *RuntimeCache
	tg         *tg.Service
	logger     *slog.Logger

	processChan  chan struct{}
	deleteChan   chan struct{}
	restrictChan chan struct{}

	penalties []*models.AntispamPenalty
}

func NewService(
	lc fx.Lifecycle,
	repo *repo.Repo,
	check *CalculatorService,
	tg *tg.Service,
	logger *slog.Logger,
) *Service {
	res := &Service{
		repo:       repo,
		calculator: check,
		tg:         tg,
		cache:      NewRuntimeCache(),
		logger:     logger.WithGroup("domain"),

		processChan:  make(chan struct{}, 1),
		deleteChan:   make(chan struct{}, 1),
		restrictChan: make(chan struct{}, 1),
		penalties:    GetAntispamPenalties(),
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
		s.cache.SetFromChatConfig(config.TgID, config)
	}

	go s.serveProcess()     //nolint:contextcheck
	go s.serveDelete()      //nolint:contextcheck
	go s.serveRestriction() //nolint:contextcheck

	return nil
}

func (s *Service) Now() int64 {
	return time.Now().Unix()
}
