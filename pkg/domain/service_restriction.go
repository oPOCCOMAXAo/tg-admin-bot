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
	RestrictNotFoundInterval = time.Hour
	RestrictRetryInterval    = 1 * time.Second
)

func (s *Service) ScheduleRestriction(
	ctx context.Context,
	data *models.Restriction,
) error {
	if data.ExecuteAt == 0 {
		data.ExecuteAt = s.Now()
	}

	defer s.wakeupRestriction()

	//nolint:wrapcheck
	return s.repo.CreateRestriction(ctx, data)
}

func (s *Service) wakeupRestriction() {
	select {
	case s.restrictChan <- struct{}{}:
	default:
	}
}

func (s *Service) serveRestriction() {
	s.wakeupRestriction()

	ticker := time.NewTicker(RestrictNotFoundInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.restrictChan:
			for {
				untilNext := s.processNextRestriction(context.Background())
				if untilNext > 0 {
					ticker.Reset(untilNext)

					break
				}
			}

		case <-ticker.C:
			for {
				untilNext := s.processNextRestriction(context.Background())
				if untilNext > 0 {
					ticker.Reset(untilNext)

					break
				}
			}
		}
	}
}

// processNextRestriction.
// Returns true if there are more restrictions to process.
func (s *Service) processNextRestriction(
	ctx context.Context,
) time.Duration {
	rst, err := s.repo.GetRestrictionForProcessing(ctx, s.Now())
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return s.getNextRestrictInterval(ctx)
		}

		s.logger.Error("processNextRestriction",
			slog.Any("error", err),
		)

		return RestrictRetryInterval
	}

	err = s.processRestriction(ctx, rst)
	if err != nil {
		s.logger.Error("processNextRestriction",
			slog.Int64("id", rst.ID),
			slog.Any("error", err),
		)
	}

	err = s.repo.DeleteRestrictionByID(ctx, rst.ID)
	if err != nil {
		s.logger.Error("processNextRestriction",
			slog.Int64("id", rst.ID),
			slog.Any("error", err),
		)
	}

	return 0
}

func (s *Service) processRestriction(
	ctx context.Context,
	rst *models.Restriction,
) error {
	// we don't care about quantity of restrictions, so apply only last unban.
	if rst.IsUnban {
		last, err := s.repo.GetLastUnbanExecuteAt(ctx, rst)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			//nolint:wrapcheck
			return err
		}

		if last != nil && last.ID != rst.ID {
			return nil
		}
	}

	if rst.SenderChatID != 0 {
		return s.restrictChat(ctx, rst)
	}

	if rst.UserID != 0 {
		return s.restrictUser(ctx, rst)
	}

	return nil
}

func (s *Service) restrictUser(
	ctx context.Context,
	rst *models.Restriction,
) error {
	if rst.IsMute {
		//nolint:wrapcheck
		return s.tg.MuteUser(ctx, &tg.MuteParams{
			ChatID:       rst.ChatID,
			UserID:       rst.UserID,
			MuteDuration: time.Duration(rst.Duration) * time.Second,
		})
	}

	return nil
}

func (s *Service) restrictChat(
	ctx context.Context,
	rst *models.Restriction,
) error {
	if rst.IsMute {
		err := s.repo.CreateRestriction(ctx, &models.Restriction{
			ExecuteAt:    s.Now() + rst.Duration,
			ChatID:       rst.ChatID,
			UserID:       rst.UserID,
			SenderChatID: rst.SenderChatID,
			IsUnban:      true,
		})
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		rst.IsBan = true
	}

	if rst.IsBan {
		//nolint:wrapcheck
		return s.tg.BanSenderChat(ctx, &tg.BanChatParams{
			ChatID:       rst.ChatID,
			SenderChatID: rst.SenderChatID,
		})
	}

	if rst.IsUnban {
		//nolint:wrapcheck
		return s.tg.UnbanSenderChat(ctx, &tg.BanChatParams{
			ChatID:       rst.ChatID,
			SenderChatID: rst.SenderChatID,
		})
	}

	return nil
}

func (s *Service) getNextRestrictInterval(
	ctx context.Context,
) time.Duration {
	msg, err := s.repo.GetFirstRestrictionAny(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return RestrictNotFoundInterval
		}

		s.logger.Error("processNextRestriction",
			slog.Int64("id", msg.ID),
			slog.Any("error", err),
		)

		return RestrictRetryInterval
	}

	now := s.Now()
	if msg.ExecuteAt > now {
		return time.Duration(msg.ExecuteAt-now) * time.Second
	}

	return RestrictRetryInterval
}
