package domain

import (
	"context"

	bmodels "github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/texts"
)

const BannedLetters = "ыэёъЫЭЁЪ"

type IsValidFunc func(context.Context, *bmodels.Message) (bool, error)

type CheckService struct {
	letters  map[rune]struct{}
	checkMap map[models.Rule]IsValidFunc
}

func NewCheckService() *CheckService {
	res := &CheckService{
		letters: texts.RuneSetFromString(BannedLetters),
	}

	res.checkMap = res.MakeMap()

	return res
}

func (s *CheckService) MakeMap() map[models.Rule]IsValidFunc {
	return map[models.Rule]IsValidFunc{
		models.RuleMuteLetters: s.VerifyBannedLetters,
	}
}

// IsValid checks if the message is valid according to the rules.
func (s *CheckService) IsValid(
	ctx context.Context,
	message *bmodels.Message,
	rules []models.Rule,
) (bool, error) {
	for _, rule := range rules {
		isValidFunc := s.checkMap[rule]
		if isValidFunc == nil {
			continue
		}

		valid, err := isValidFunc(ctx, message)
		if err != nil || !valid {
			return false, err
		}
	}

	return true, nil
}

func (s *CheckService) IsTextContainsBannedLetters(text string) bool {
	for _, r := range text {
		if _, ok := s.letters[r]; ok {
			return true
		}
	}

	return false
}

func (s *CheckService) VerifyBannedLetters(
	_ context.Context,
	message *bmodels.Message,
) (bool, error) {
	if message.Text != "" && s.IsTextContainsBannedLetters(message.Text) {
		return false, nil
	}

	if message.Caption != "" && s.IsTextContainsBannedLetters(message.Caption) {
		return false, nil
	}

	return true, nil
}
