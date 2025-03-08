package domain

import (
	"context"
	"errors"
	"log/slog"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
)

const DebugMessageDefaultTTLSeconds = 30

type Params struct {
	tg.ReplyDebugParams
	DeleteAfterSeconds int64 // 0 means no deletion
}

func (s *Service) ReplyDebug(
	ctx context.Context,
	params *Params,
) error {
	msg, err := s.tg.ReplyDebugOrNil(ctx, &params.ReplyDebugParams)
	if err != nil && !errors.Is(err, models.ErrNothingChanged) {
		s.logger.ErrorContext(ctx, "ReplyDebug",
			slog.Int64("message_id", int64(msg.ID)),
			slog.Int64("chat_id", params.ChatID),
			slog.Any("error", err),
		)
	}

	if msg == nil || params.DeleteAfterSeconds <= 0 {
		return nil
	}

	err = s.ScheduleDelete(ctx, &models.MessageDelete{
		ChatID:    params.ChatID,
		MessageID: int64(msg.ID),
		ExecuteAt: s.Now() + params.DeleteAfterSeconds,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "ReplyDebug",
			slog.Int64("message_id", int64(msg.ID)),
			slog.Int64("chat_id", params.ChatID),
			slog.Any("error", err),
		)
	}

	return nil
}
