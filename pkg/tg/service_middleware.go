package tg

import (
	"github.com/opoccomaxao/tg-instrumentation/router"
)

func (s *Service) RequireCallbackFromAdmin(ctx *router.Context) {
	update := ctx.Update()
	if update.CallbackQuery == nil {
		ctx.Abort()

		return
	}

	err := s.CheckUserMemberPermissions(ctx.Context(), &CheckUserMemberPermissionsParams{
		ChatID:        update.CallbackQuery.Message.Message.Chat.ID,
		UserID:        update.CallbackQuery.From.ID,
		Username:      update.CallbackQuery.From.Username,
		RequiredTypes: memberTypesAdmin,
	})
	if err != nil {
		ctx.Error(err)
		ctx.Abort()

		return
	}
}

func (s *Service) RequireMessageFromAdmin(ctx *router.Context) {
	update := ctx.Update()
	if update.Message == nil {
		ctx.Abort()

		return
	}

	err := s.CheckUserMemberPermissions(ctx.Context(), &CheckUserMemberPermissionsParams{
		ChatID:        update.Message.Chat.ID,
		UserID:        update.Message.From.ID,
		Username:      update.Message.From.Username,
		RequiredTypes: memberTypesAdmin,
	})
	if err != nil {
		ctx.Error(err)
		ctx.Abort()

		return
	}
}
