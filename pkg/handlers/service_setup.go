package handlers

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/views"
	"github.com/opoccomaxao/tg-instrumentation/router"
	"github.com/pkg/errors"
)

func (s *Service) Setup(ctx *router.Context) {
	update := ctx.Update()

	if update.Message.Chat.Type == models.ChatTypePrivate {
		ctx.LogError2(ctx.RespondPrivateMessageText("Ця команда доступна тільки в групах"))

		return
	}

	view := views.Setup{
		ChatID:    update.Message.Chat.ID,
		MessageID: int64(update.Message.ID),
	}

	err := s.fillSetupView(ctx, &view)
	if err != nil {
		ctx.Error(err)

		return
	}

	ctx.LogError2(ctx.SendMessage(view.SendMessageParams()))
}

func (s *Service) fillSetupView(
	ctx *router.Context,
	view *views.Setup,
) error {
	cfg, err := s.domain.GetOrCreateChatByTgID(
		ctx.Context(),
		view.ChatID,
	)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	chat, err := s.tg.GetChatMember(ctx.Context(), &bot.GetChatMemberParams{
		ChatID: view.ChatID,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	view.Config = cfg

	if chat.Administrator != nil {
		view.CanRestrictMembers = chat.Administrator.CanRestrictMembers
	}

	return nil
}
