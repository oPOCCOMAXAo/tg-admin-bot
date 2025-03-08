package domain

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
)

func (s *Service) HandleMessage(
	ctx context.Context,
	update *apimodels.Update,
) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID

	cfg := s.cache.GetConfig(chatID)
	if cfg == nil || cfg.Enabled == nil {
		return nil
	}

	info := models.MessageInfo{
		Time:         int64(update.Message.Date),
		ChatID:       chatID,
		MessageID:    int64(update.Message.ID),
		UserID:       0,
		SenderChatID: 0,
		IsProcessed:  false,
		Score:        0,
	}

	if update.Message.SenderChat != nil {
		info.SenderChatID = update.Message.SenderChat.ID
	} else {
		info.UserID = update.Message.From.ID
	}

	err := s.calculator.CalculateIntoInfo(ctx, update.Message, &info, cfg)
	if err != nil {
		return err
	}

	err = s.repo.CreateMessageInfo(ctx, &info)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	s.wakeupProcess()

	return nil
}
