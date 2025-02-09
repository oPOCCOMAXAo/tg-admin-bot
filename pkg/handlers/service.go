package handlers

import (
	"log/slog"

	"github.com/opoccomaxao/tg-admin-bot/pkg/domain"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
)

type Service struct {
	logger *slog.Logger
	domain *domain.Service
	tg     *tg.Service
}

func NewService(
	logger *slog.Logger,
	domain *domain.Service,
	tg *tg.Service,
) *Service {
	return &Service{
		logger: logger,
		domain: domain,
		tg:     tg,
	}
}
