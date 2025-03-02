package handlers

import (
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/views"
	"github.com/opoccomaxao/tg-instrumentation/query"
	"github.com/opoccomaxao/tg-instrumentation/router"
)

type SetupSetRequest struct {
	Set    int64 // bool: 0 - false, other - true; int64: full value.
	RuleID models.ConfigID
}

func (s *Service) SetupSet(ctx *router.Context) {
	update := ctx.Update()

	var req SetupSetRequest

	query := query.Decode(update.CallbackQuery.Data)
	query.GetInt64Into("setup_set", &req.Set)
	query.GetInt64Into("id", req.RuleID.Int64Ref())

	err := s.domain.UpdateChatConfigInt(
		ctx.Context(),
		update.CallbackQuery.Message.Message.Chat.ID,
		req.RuleID,
		req.Set,
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
