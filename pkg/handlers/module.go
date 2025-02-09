package handlers

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
	"github.com/opoccomaxao/tg-instrumentation/router"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("handlers",
		fx.Provide(NewService, fx.Private),
		fx.Invoke(RegisterHandlers),
	)
}

func RegisterHandlers(
	tg *tg.Service,
	router *router.Router,
	service *Service,
) error {
	defer tg.SetupCommands(context.Background())

	router.Text("/start", service.Start).
		WithDescription(apimodels.LCAll, apimodels.CSAllPrivateChats, "Start").
		WithDescription(apimodels.LCUk, apimodels.CSAllPrivateChats, "Почати роботу")

	router.Text("/setup", service.Setup).
		WithDescription(apimodels.LCAll, apimodels.CSAllChatAdministrators, "Setup").
		WithDescription(apimodels.LCUk, apimodels.CSAllChatAdministrators, "Налаштування")

	router.Custom(func(update *apimodels.Update) bool {
		return update.Message != nil || update.EditedMessage != nil
	}, service.GroupMessage)

	return nil
}
