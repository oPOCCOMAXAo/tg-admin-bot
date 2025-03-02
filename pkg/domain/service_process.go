package domain

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"github.com/samber/lo"
)

const (
	MuteRuLettersDuration = 5 * time.Minute

	DebugMessageTTLSeconds = 30
)

func (s *Service) wakeupProcess() {
	select {
	case s.processChan <- struct{}{}:
	default:
	}
}

func (s *Service) serveProcess() {
	s.wakeupProcess()

	for range s.processChan {
		for s.processNextMessage(context.Background()) {
		}
	}
}

// processNextMessage.
// Returns true if there are more messages to process.
func (s *Service) processNextMessage(
	ctx context.Context,
) bool {
	info, err := s.repo.GetMessageInfoForProcessing(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return false
		}

		s.logger.Error("processNextMessage",
			slog.Any("error", err),
		)

		return true
	}

	err = s.processMessage(ctx, info)
	if err != nil {
		s.logger.Error("processNextMessage",
			slog.Int64("id", info.ID),
			slog.Any("error", err),
		)

		return true
	}

	info.IsProcessed = true

	err = s.repo.UpdateMessageInfo(ctx, info)
	if err != nil {
		s.logger.Error("processNextMessage",
			slog.Int64("id", info.ID),
			slog.Any("error", err),
		)
	}

	s.logger.InfoContext(ctx, "processNextMessage",
		slog.Int64("id", info.ID),
		slog.Int64("chat_id", info.ChatID),
		slog.Int64("user_id", info.UserID),
		slog.Int64("score", int64(info.Score)),
	)

	return true
}

// processMessage should only update info fields and call tg api.
func (s *Service) processMessage(
	ctx context.Context,
	info *models.MessageInfo,
) error {
	cfg := s.cache.GetConfig(info.ChatID)
	if cfg == nil || cfg.Enabled == nil {
		return nil
	}

	for _, cfgID := range cfg.Enabled {
		switch cfgID {
		case models.CfgEnabledMuteRuLetters:
			err := s.processRuLetters(ctx, info)
			if err != nil {
				return err
			}
		case models.CfgEnabledAntispam:
			err := s.processAntispam(ctx, info, &cfg.Antispam)
			if err != nil {
				return err
			}
		default:
			s.logger.WarnContext(ctx, "processMessage",
				slog.Int64("id", info.ID),
				slog.String("rule", cfgID.StringID()),
				slog.String("error", "unknown rule"),
			)
		}
	}

	return nil
}

func (s *Service) processRuLetters(
	ctx context.Context,
	info *models.MessageInfo,
) error {
	if !info.HasRULetters {
		return nil
	}

	err := s.tg.ReactMessage(ctx, &tg.ReactParams{
		ChatID:        info.ChatID,
		MessageID:     info.MessageID,
		ReactionEmoji: tg.ReactionSwearing,
	})
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	err = s.tg.MuteUser(ctx, &tg.MuteParams{
		ChatID:       info.ChatID,
		UserID:       info.UserID,
		MuteDuration: MuteRuLettersDuration,
	})
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	return nil
}

//nolint:mnd,cyclop,funlen
func (s *Service) processAntispam(
	ctx context.Context,
	info *models.MessageInfo,
	cfg *models.AntispamConfig,
) error {
	prev, err := s.repo.GetMessageInfoPrevious(ctx, info)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		//nolint:wrapcheck
		return err
	}

	if prev != nil {
		info.IsFast = info.Time-prev.Time < 5
	}

	info.Score = CalculateScore(info)

	var (
		banDuration  time.Duration
		warnRequired bool
		debugData    [][]string
	)

	if cfg.Debug {
		debugData = make([][]string, 0, 1+len(s.penalties))
		debugData = append(debugData, []string{"score", strconv.FormatUint(uint64(info.Score), 10)})
	}

	for _, penalty := range s.penalties {
		score, err := s.repo.GetMessagePrevScore(
			ctx,
			info,
			info.Time-int64(penalty.CheckInterval.Seconds()),
			info.Time,
		)
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		score += int64(info.Score)

		if score >= penalty.MaxScore {
			banDuration = max(banDuration, penalty.PenaltyTime)
		}

		if score >= penalty.MaxScore/2 {
			warnRequired = true
		}

		if cfg.Debug {
			debugData = append(debugData, []string{
				penalty.Name,
				strconv.FormatInt(score, 10) + "/" + strconv.FormatInt(penalty.MaxScore, 10),
			})
		}
	}

	if banDuration > 0 {
		err := s.tg.MuteUser(ctx, &tg.MuteParams{
			ChatID:       info.ChatID,
			UserID:       info.UserID,
			MuteDuration: banDuration,
		})
		if err != nil {
			//nolint:wrapcheck
			return err
		}
	}

	if warnRequired {
		err := s.tg.ReactMessage(ctx, &tg.ReactParams{
			ChatID:        info.ChatID,
			MessageID:     info.MessageID,
			ReactionEmoji: lo.Ternary(banDuration > 0, tg.ReactionSwearing, tg.ReactionSee),
		})
		if err != nil {
			//nolint:wrapcheck
			return err
		}
	}

	if cfg.Debug {
		res, err := s.tg.ReplyDebugOrNil(ctx, &tg.ReplyDebugParams{
			ChatID:           info.ChatID,
			ReplyToMessageID: info.MessageID,
			Data:             debugData,
		})
		if err != nil && !errors.Is(err, models.ErrNothingChanged) {
			s.logger.ErrorContext(ctx, "processAntispam",
				slog.Int64("id", info.ID),
				slog.Any("error", err),
			)
		}

		if res != nil {
			err = s.ScheduleDelete(ctx, &models.MessageDelete{
				ChatID:    res.Chat.ID,
				MessageID: int64(res.ID),
				ExecuteAt: s.Now() + DebugMessageTTLSeconds,
			})
			if err != nil {
				s.logger.ErrorContext(ctx, "processAntispam",
					slog.Int64("id", info.ID),
					slog.Any("error", err),
				)
			}
		}
	}

	return nil
}
