package handlers

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/tg"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg/middleware"
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
	pkgrouter "github.com/opoccomaxao/tg-instrumentation/router"
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
	router *pkgrouter.Router,
	service *Service,
) error {
	//nolint:errcheck
	defer router.UpdateCommandsDescription(context.Background())

	router.Text("/start", service.Start).
		WithDescription(apimodels.LCAll, apimodels.CSAllPrivateChats, "Start").
		WithDescription(apimodels.LCUk, apimodels.CSAllPrivateChats, "Почати роботу")

	router.Text("/setup",
		tg.RequireMessageFromAdmin,
		service.Setup,
	).
		WithDescription(apimodels.LCAll, apimodels.CSAllChatAdministrators, "Setup").
		WithDescription(apimodels.LCUk, apimodels.CSAllChatAdministrators, "Налаштування")

	router.Callback("setup_update",
		pkgrouter.AutoAnswerCallbackQuery(),
		middleware.RequireCallbackMessage,
		tg.RequireCallbackFromAdmin,
		service.SetupUpdate,
	)

	router.Callback("setup_set",
		pkgrouter.AutoAnswerCallbackQuery(),
		middleware.RequireCallbackMessage,
		tg.RequireCallbackFromAdmin,
		service.SetupSet,
	)

	router.Callback("setup_apply",
		pkgrouter.AutoAnswerCallbackQuery(),
		middleware.RequireCallbackMessage,
		tg.RequireCallbackFromAdmin,
		service.SetupApply,
	)

	router.Callback("delete_self_admin",
		pkgrouter.AutoAnswerCallbackQuery(),
		middleware.RequireCallbackMessage,
		tg.RequireCallbackFromAdmin,
		service.DeleteSelf,
	)

	router.Callback("delete_self",
		pkgrouter.AutoAnswerCallbackQuery(),
		middleware.RequireCallbackMessage,
		service.DeleteSelf,
	)

	router.Custom(func(update *apimodels.Update) bool {
		return update.Message != nil || update.EditedMessage != nil
	}, service.GroupMessage)

	return nil
}
