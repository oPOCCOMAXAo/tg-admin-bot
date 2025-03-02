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

func (s *Service) UpdateChatConfigInt(
	ctx context.Context,
	tgID int64,
	rule models.ConfigID,
	value int64,
) error {
	//nolint:wrapcheck
	return s.repo.UpdateChatConfigInt(ctx, tgID, rule, value)
}

func (s *Service) CacheChatRuntimeConfig(
	ctx context.Context,
	tgID int64,
) error {
	chat, err := s.repo.GetOrCreateChatByTgID(ctx, tgID)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	s.cache.SetFromChatConfig(chat.TgID, chat)

	return nil
}
