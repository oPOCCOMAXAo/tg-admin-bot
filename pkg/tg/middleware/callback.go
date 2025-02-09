package middleware

import (
	"github.com/opoccomaxao/tg-instrumentation/router"
)

func RequireCallbackMessage(ctx *router.Context) {
	update := ctx.Update()
	if update.CallbackQuery == nil ||
		update.CallbackQuery.Message.Message == nil {
		ctx.Abort()
	}
}
