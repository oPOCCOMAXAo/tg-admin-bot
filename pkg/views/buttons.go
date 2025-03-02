package views

import (
	bmodels "github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-instrumentation/query"
)

//nolint:gochecknoglobals
var (
	btnCloseAny = bmodels.InlineKeyboardButton{
		Text: "Закрити",
		CallbackData: query.Command("delete_self").
			Encode(),
	}

	btnCloseAdmin = bmodels.InlineKeyboardButton{
		Text: "Закрити",
		CallbackData: query.Command("delete_self_admin").
			Encode(),
	}

	btnEnableMuteLetters = bmodels.InlineKeyboardButton{
		Text: "Увімкнути мут літер 'ыэёъЫЭЁЪ'",
		CallbackData: query.Command("setup_set").
			WithParam("id", models.CfgEnabledMuteRuLetters.StringID()).
			WithParam("setup_set", "1").
			Encode(),
	}

	btnDisableMuteLetters = bmodels.InlineKeyboardButton{
		Text: "Вимкнути мут літер 'ыэёъЫЭЁЪ'",
		CallbackData: query.Command("setup_set").
			WithParam("id", models.CfgEnabledMuteRuLetters.StringID()).
			WithParam("setup_set", "0").
			Encode(),
	}

	btnEnableAntispam = bmodels.InlineKeyboardButton{
		Text: "Увімкнути антиспам",
		CallbackData: query.Command("setup_set").
			WithParam("id", models.CfgEnabledAntispam.StringID()).
			WithParam("setup_set", "1").
			Encode(),
	}

	btnDisableAntispam = bmodels.InlineKeyboardButton{
		Text: "Вимкнути антиспам",
		CallbackData: query.Command("setup_set").
			WithParam("id", models.CfgEnabledAntispam.StringID()).
			WithParam("setup_set", "0").
			Encode(),
	}

	btnAntispamDebugOn = bmodels.InlineKeyboardButton{
		Text: "debug",
		CallbackData: query.Command("setup_set").
			WithParam("id", models.CfgAntispamDebug.StringID()).
			WithParam("setup_set", "1").
			Encode(),
	}

	btnAntispamDebugOff = bmodels.InlineKeyboardButton{
		Text: "debug",
		CallbackData: query.Command("setup_set").
			WithParam("id", models.CfgAntispamDebug.StringID()).
			WithParam("setup_set", "0").
			Encode(),
	}
)
