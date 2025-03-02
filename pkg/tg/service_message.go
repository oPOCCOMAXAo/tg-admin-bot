package tg

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	bmodels "github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-instrumentation/texts"
	"github.com/pkg/errors"
)

type ReactParams struct {
	ChatID        int64
	MessageID     int64
	ReactionEmoji string
}

func (s *Service) ReactMessage(
	ctx context.Context,
	params *ReactParams,
) error {
	_, err := s.client.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    params.ChatID,
		MessageID: int(params.MessageID),
		Reaction: []bmodels.ReactionType{
			{
				Type: bmodels.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &bmodels.ReactionTypeEmoji{
					Emoji: params.ReactionEmoji,
				},
			},
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type ReplyDebugParams struct {
	ChatID           int64
	ReplyToMessageID int64 // optional
	Text             string
	Data             [][]string
}

// ReplyDebugOrNil sends a debug message to the chat.
// If the message is empty, returns ErrNothingChanged.
func (s *Service) ReplyDebugOrNil(
	ctx context.Context,
	params *ReplyDebugParams,
) (*bmodels.Message, error) {
	msg := bot.SendMessageParams{
		ChatID:    params.ChatID,
		ParseMode: bmodels.ParseModeHTML,
		ReplyMarkup: bmodels.InlineKeyboardMarkup{
			InlineKeyboard: [][]bmodels.InlineKeyboardButton{
				{{Text: "close", CallbackData: "delete_self"}},
			},
		},
		DisableNotification: true,
	}

	if params.ReplyToMessageID != 0 {
		msg.ReplyParameters = &bmodels.ReplyParameters{
			MessageID: int(params.ReplyToMessageID),
			ChatID:    params.ChatID,
		}
	}

	text := []string{}

	if params.Text != "" {
		text = append(text,
			"<b>Debug text</b>:",
			texts.EscapeHTML(params.Text),
		)
	}

	if len(params.Data) > 0 {
		text = append(text, "<b>Debug data</b>:")

		for _, row := range params.Data {
			if len(row) == 0 {
				continue
			}

			title := texts.EscapeHTML(row[0])
			value := strings.Join(row[1:], " ")
			text = append(text, texts.EscapeHTML(title)+": "+texts.EscapeHTML(value))
		}
	}

	msg.Text = strings.Join(text, "\n")
	if len(msg.Text) == 0 {
		return nil, errors.WithStack(models.ErrNothingChanged)
	}

	res, err := s.client.SendMessage(ctx, &msg)
	if err != nil {
		return res, errors.WithStack(err)
	}

	return res, nil
}

type DeleteMessageParams struct {
	ChatID    int64
	MessageID int64
}

func (s *Service) DeleteMessage(
	ctx context.Context,
	params *DeleteMessageParams,
) error {
	_, err := s.client.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    params.ChatID,
		MessageID: int(params.MessageID),
	})
	if err != nil {
		if errors.Is(err, bot.ErrorBadRequest) {
			s.logger.ErrorContext(ctx, "DeleteMessage",
				slog.Int64("chat_id", params.ChatID),
				slog.Int64("message_id", params.MessageID),
				slog.Any("error", err),
			)

			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}
