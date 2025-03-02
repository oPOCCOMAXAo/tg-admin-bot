package views

import (
	"strings"

	"github.com/go-telegram/bot"
	bmodels "github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/texts"
)

type Setup struct {
	ChatID             int64
	MessageID          int64
	Config             *models.ChatConfig
	CanRestrictMembers bool
}

func (s *Setup) getUpdateButtonRow() []bmodels.InlineKeyboardButton {
	return []bmodels.InlineKeyboardButton{
		{
			Text:         "Оновити",
			CallbackData: "setup_update",
		},
	}
}

func (s *Setup) getApplyButtonRow() []bmodels.InlineKeyboardButton {
	return []bmodels.InlineKeyboardButton{
		{
			Text:         "Застосувати",
			CallbackData: "setup_apply",
		},
	}
}

func (s *Setup) getMuteLettersButtonRow() []bmodels.InlineKeyboardButton {
	if s.Config == nil {
		return nil
	}

	res := []bmodels.InlineKeyboardButton{}

	if s.Config.EnabledMuteLetters {
		res = append(res, btnDisableMuteLetters)
	} else {
		res = append(res, btnEnableMuteLetters)
	}

	return res
}

func (s *Setup) getAntispamButtonRow() []bmodels.InlineKeyboardButton {
	if s.Config == nil {
		return nil
	}

	res := []bmodels.InlineKeyboardButton{}

	if s.Config.EnabledAntispam {
		res = append(res, btnDisableAntispam)

		if s.Config.AntispamDebug {
			res = append(res, btnAntispamDebugOff)
		} else {
			res = append(res, btnAntispamDebugOn)
		}
	} else {
		res = append(res, btnEnableAntispam)
	}

	return res
}

func (s *Setup) ReplyMarkup() bmodels.ReplyMarkup {
	res := &bmodels.InlineKeyboardMarkup{}
	res.InlineKeyboard = append(res.InlineKeyboard, s.getUpdateButtonRow())
	res.InlineKeyboard = append(res.InlineKeyboard, s.getMuteLettersButtonRow())
	res.InlineKeyboard = append(res.InlineKeyboard, s.getAntispamButtonRow())
	res.InlineKeyboard = append(res.InlineKeyboard, s.getApplyButtonRow())
	res.InlineKeyboard = append(res.InlineKeyboard, []bmodels.InlineKeyboardButton{btnCloseAdmin})

	return res
}

func (s *Setup) getRequiredPermissions() []string {
	res := []string{}

	if !s.CanRestrictMembers {
		res = append(res, "can_restrict_members")
	}

	return res
}

func (s *Setup) getText() string {
	const cmdDescPrefix = "  - "

	lines := []string{
		"<b>Налаштування поточного чату:</b>\n",
	}

	if s.Config == nil {
		lines = append(lines, "<i>Налаштування не знайдені.</i>")
	} else {
		lines = append(lines, "> Мут літер 'ыэёъЫЭЁЪ': "+OnOff(s.Config.EnabledMuteLetters))
		lines = append(lines, "<i>"+texts.JoinListLinesWithPrefix([]string{
			"при наявності цих літер в повідомленні видає мут на 5 хвилин",
		}, cmdDescPrefix)+"</i>")

		lines = append(lines, "")

		lines = append(lines, "> Антиспам: "+OnOff(s.Config.EnabledAntispam))
		lines = append(lines, "<i>"+texts.JoinListLinesWithPrefix([]string{
			"рахує рейтинг і видає мут за кількома правилами",
		}, cmdDescPrefix)+"</i>")

		if s.Config.EnabledAntispam {
			lines = append(lines, ">> debug: "+OnOff(s.Config.AntispamDebug))

			lines = append(lines, "<i>"+texts.JoinListLinesWithPrefix([]string{
				"виводить рейтинг повідомлень у debug-повідомленні",
				"автовидаляє ці повідомлення через 30 секунд",
				"дає можливість кожному видалити debug-повідомлення",
			}, cmdDescPrefix)+"</i>")
		}
	}

	lines = append(lines, "\nАдміністратори чату можуть змінити ці налаштування кнопками нижче.")

	perms := s.getRequiredPermissions()
	if len(perms) > 0 {
		lines = append(lines, "\n<b>Для зміни налаштувань потрібні додаткові права:</b>")
		for _, perm := range perms {
			lines = append(lines, "- "+perm)
		}
	}

	return strings.Join(lines, "\n")
}

func (s *Setup) SendMessageParams() *bot.SendMessageParams {
	res := &bot.SendMessageParams{
		ChatID:         s.ChatID,
		ReplyMarkup:    s.ReplyMarkup(),
		Text:           s.getText(),
		ParseMode:      bmodels.ParseModeHTML,
		ProtectContent: true,
	}

	if s.MessageID != 0 {
		res.ReplyParameters = &bmodels.ReplyParameters{
			MessageID: int(s.MessageID),
			ChatID:    s.ChatID,
		}
	}

	return res
}

func (s *Setup) EditMessageParams() *bot.EditMessageTextParams {
	res := &bot.EditMessageTextParams{
		ChatID:      s.ChatID,
		ReplyMarkup: s.ReplyMarkup(),
		Text:        s.getText(),
		ParseMode:   bmodels.ParseModeHTML,
		MessageID:   int(s.MessageID),
	}

	return res
}
