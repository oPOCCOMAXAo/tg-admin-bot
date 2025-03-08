package handlers

import (
	_ "embed"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-instrumentation/router"
	"github.com/samber/lo"
)

//go:embed raw/help_uk.tgmd
var textHelpUK string

func (s *Service) Help(ctx *router.Context) {
	ctx.LogError2(ctx.RespondPrivateMessage(&bot.SendMessageParams{
		Text:      textHelpUK,
		ParseMode: models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: lo.ToPtr(true),
		},
	}))
}
