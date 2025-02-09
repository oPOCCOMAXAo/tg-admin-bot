package domain

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
)

func (s *Service) GetOrCreateChatByTgID(
	ctx context.Context,
	tgID int64,
) (*models.ChatConfig, error) {
	//nolint:wrapcheck
	return s.repo.GetOrCreateChatByTgID(ctx, tgID)
}

func (s *Service) UpdateChatRule(
	ctx context.Context,
	tgID int64,
	rule models.Rule,
	enabled bool,
) error {
	//nolint:wrapcheck
	return s.repo.UpdateChatRule(ctx, tgID, rule, enabled)
}

func (s *Service) CacheChatRules(
	ctx context.Context,
	tgID int64,
) error {
	chat, err := s.repo.GetOrCreateChatByTgID(ctx, tgID)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	s.cache.SetRules(chat.TgID, chat.RulesList())

	return nil
}
