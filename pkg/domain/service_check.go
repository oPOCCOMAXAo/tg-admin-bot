package domain

import (
	"context"
	"time"

	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
)

const (
	MuteDuration = 5 * time.Minute
)

func (s *Service) ProcessMessageCheck(
	ctx context.Context,
	update *apimodels.Update,
) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID

	rules := s.cache.GetRules(chatID)
	if rules == nil {
		return nil
	}

	valid, err := s.check.IsValid(ctx, update.Message, rules)
	if err != nil {
		return err
	}

	if valid {
		return nil
	}

	//nolint:wrapcheck
	return s.tg.Restrict(ctx, &tg.RestrictParams{
		ChatID:       chatID,
		UserID:       update.Message.From.ID,
		MessageID:    int64(update.Message.ID),
		MuteDuration: MuteDuration,
	})
}
