package tg

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
)

type RestrictParams struct {
	ChatID       int64
	UserID       int64
	MessageID    int64
	MuteDuration time.Duration
}

func (s *Service) Restrict(
	ctx context.Context,
	params *RestrictParams,
) error {
	var err error

	_, err = s.client.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    params.ChatID,
		MessageID: int(params.MessageID),
		Reaction: []models.ReactionType{
			{
				Type: models.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &models.ReactionTypeEmoji{
					Emoji: "🤬",
				},
			},
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = s.client.RestrictChatMember(ctx, &bot.RestrictChatMemberParams{
		ChatID:      params.ChatID,
		UserID:      params.UserID,
		UntilDate:   int(time.Now().Add(params.MuteDuration).Unix()),
		Permissions: &models.ChatPermissions{},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
