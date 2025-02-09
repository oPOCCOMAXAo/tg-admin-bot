package handlers

import (
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/texts"
	"github.com/opoccomaxao/tg-admin-bot/pkg/views"
	"github.com/opoccomaxao/tg-instrumentation/router"
)

type SetupSetRequest struct {
	Set    int64 // 0 - off, 1 - on
	RuleID models.Rule
}

func (s *Service) SetupSet(ctx *router.Context) {
	update := ctx.Update()

	var req SetupSetRequest

	query := texts.DecodeQuery(update.CallbackQuery.Data)
	query.GetInt64Into("setup_set", &req.Set)
	query.GetInt64Into("id", req.RuleID.Int64Ref())

	err := s.domain.UpdateChatRule(
		ctx.Context(),
		update.CallbackQuery.Message.Message.Chat.ID,
		req.RuleID,
		req.Set != 0,
	)
	if err != nil {
		ctx.Error(err)

		return
	}

	view := views.Setup{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: int64(update.CallbackQuery.Message.Message.ID),
	}

	err = s.fillSetupView(ctx, &view)
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
