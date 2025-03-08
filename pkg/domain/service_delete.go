package domain

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
)

const (
	DeleteNotFoundInterval = time.Hour
	DeleteRetryInterval    = 250 * time.Millisecond
)

func (s *Service) ScheduleDelete(
	ctx context.Context,
	data *models.MessageDelete,
) error {
	if data.ExecuteAt == 0 {
		data.ExecuteAt = s.Now()
	}

	defer s.wakeupDelete()

	//nolint:wrapcheck
	return s.repo.CreateMessageDelete(ctx, data)
}

func (s *Service) wakeupDelete() {
	select {
	case s.deleteChan <- struct{}{}:
	default:
	}
}

func (s *Service) serveDelete() {
	s.wakeupDelete()

	ticker := time.NewTicker(DeleteNotFoundInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.deleteChan:
			for {
				untilNext := s.processNextDelete(context.Background())
				if untilNext > 0 {
					ticker.Reset(untilNext)

					break
				}
			}

		case <-ticker.C:
			for {
				untilNext := s.processNextDelete(context.Background())
				if untilNext > 0 {
					ticker.Reset(untilNext)

					break
				}
			}
		}
	}
}

func (s *Service) processNextDelete(
	ctx context.Context,
) time.Duration {
	msg, err := s.repo.GetFirstMessageDeleteUntilTime(ctx, s.Now())
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return s.getNextDeleteInterval(ctx)
		}

		s.logger.Error("processNextDelete",
			slog.Any("error", err),
		)

		return DeleteRetryInterval
	}

	err = s.tg.DeleteMessage(ctx, &tg.DeleteMessageParams{
		ChatID:    msg.ChatID,
		MessageID: msg.MessageID,
	})
	if err != nil {
		s.logger.Error("processNextDelete",
			slog.Int64("id", msg.ID),
			slog.Any("error", err),
		)

		return DeleteRetryInterval
	}

	err = s.repo.DeleteMessageDeleteByID(ctx, msg.ID)
	if err != nil {
		s.logger.Error("processNextDelete",
			slog.Int64("id", msg.ID),
			slog.Any("error", err),
		)
	}

	return 0
}

func (s *Service) getNextDeleteInterval(
	ctx context.Context,
) time.Duration {
	msg, err := s.repo.GetFirstMessageDeleteAny(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return DeleteNotFoundInterval
		}

		s.logger.Error("processNextDelete",
			slog.Int64("id", msg.ID),
			slog.Any("error", err),
		)

		return DeleteRetryInterval
	}

	now := s.Now()
	if msg.ExecuteAt > now {
		return time.Duration(msg.ExecuteAt-now) * time.Second
	}

	return DeleteRetryInterval
}
