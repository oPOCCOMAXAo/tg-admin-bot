package handlers

import (
	"github.com/go-telegram/bot"
	"github.com/opoccomaxao/tg-instrumentation/router"
)

func (s *Service) SetupApply(ctx *router.Context) {
	update := ctx.Update()

	err := s.domain.CacheChatRuntimeConfig(
		ctx.Context(),
		update.CallbackQuery.Message.Message.Chat.ID,
	)
	if err != nil {
		ctx.Error(err)

		return
	}

	_, _ = ctx.AnswerCallbackQuery(&bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "Правила чату оновлено",
		ShowAlert:       true,
		CacheTime:       1,
	})
}
