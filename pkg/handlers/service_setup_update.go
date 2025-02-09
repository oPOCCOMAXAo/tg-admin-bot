package handlers

import (
	"github.com/opoccomaxao/tg-admin-bot/pkg/views"
	"github.com/opoccomaxao/tg-instrumentation/router"
)

func (s *Service) SetupUpdate(ctx *router.Context) {
	update := ctx.Update()

	view := views.Setup{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: int64(update.CallbackQuery.Message.Message.ID),
	}

	err := s.fillSetupView(ctx, &view)
	if err != nil {
		ctx.Error(err)

		return
	}

	_, err = ctx.EditMessageText(view.EditMessageParams())
	if err != nil {
		ctx.Error(err)

		return
	}
}
