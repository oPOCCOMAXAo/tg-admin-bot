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
	EnabledMuteLetters bool
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
	res := []bmodels.InlineKeyboardButton{
		{},
	}

	if s.EnabledMuteLetters {
		res[0].Text = "Вимкнути мут літер 'ыэёъЫЭЁЪ'"
		res[0].CallbackData = texts.QueryCommand("setup_set").
			AddParam("id", models.RuleMuteLetters.StringID()).
			AddParam("setup_set", "0").
			Encode()
	} else {
		res[0].Text = "Увімкнути мут літер 'ыэёъЫЭЁЪ'"
		res[0].CallbackData = texts.QueryCommand("setup_set").
			AddParam("id", models.RuleMuteLetters.StringID()).
			AddParam("setup_set", "1").
			Encode()
	}

	return res
}

func (s *Setup) ReplyMarkup() bmodels.ReplyMarkup {
	res := &bmodels.InlineKeyboardMarkup{}
	res.InlineKeyboard = append(res.InlineKeyboard, s.getUpdateButtonRow())
	res.InlineKeyboard = append(res.InlineKeyboard, s.getMuteLettersButtonRow())
	res.InlineKeyboard = append(res.InlineKeyboard, s.getApplyButtonRow())

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
	lines := []string{
		"<b>Налаштування поточного чату:</b>\n",
		"Мут літер 'ыэёъЫЭЁЪ': " + OnOff(s.EnabledMuteLetters),
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
