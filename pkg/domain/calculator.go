package domain

import (
	"context"
	"slices"
	"unicode"
	"unicode/utf8"

	bmodels "github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/texts"
)

const BannedLetters = "ыэёъЫЭЁЪ"

type CalculateFunc func(
	context.Context,
	*bmodels.Message,
	*models.MessageInfo,
	*models.RuntimeConfig,
) error

type CalculatorService struct {
	letters map[rune]struct{}
	funcMap map[models.ConfigID]CalculateFunc
}

func NewCalculatorService() *CalculatorService {
	res := &CalculatorService{
		letters: texts.RuneSetFromString(BannedLetters),
	}

	res.funcMap = res.makeMap()

	return res
}

func (s *CalculatorService) makeMap() map[models.ConfigID]CalculateFunc {
	return map[models.ConfigID]CalculateFunc{
		models.CfgEnabledMuteRuLetters: s.calculateRuLetters,
		models.CfgEnabledAntispam:      s.calculateAntispam,
	}
}

func (*CalculatorService) getText(
	message *bmodels.Message,
) string {
	switch {
	case message.Text != "":
		return message.Text
	case message.Caption != "":
		return message.Caption
	default:
		return ""
	}
}

//nolint:cyclop
func (*CalculatorService) countMedias(
	message *bmodels.Message,
) uint8 {
	var count uint8

	if message.Photo != nil {
		count++
	}

	if message.Sticker != nil {
		count++
	}

	if message.Video != nil {
		count++
	}

	if message.Audio != nil {
		count++
	}

	if message.Voice != nil {
		count++
	}

	if message.Document != nil ||
		message.Animation != nil {
		count++
	}

	if message.VideoNote != nil {
		count++
	}

	if message.Contact != nil {
		count++
	}

	if message.Location != nil {
		count++
	}

	if message.Venue != nil {
		count++
	}

	if message.Game != nil {
		count++
	}

	return count
}

func (s *CalculatorService) IsTextContainsRuLetters(text string) bool {
	for _, r := range text {
		if _, ok := s.letters[r]; ok {
			return true
		}
	}

	return false
}

func (s *CalculatorService) countTextUpper(text string) int {
	upper := 0

	for _, r := range text {
		if unicode.IsUpper(r) {
			upper++
		}
	}

	return upper
}

func (s *CalculatorService) CalculateIntoInfo(
	ctx context.Context,
	message *bmodels.Message,
	info *models.MessageInfo,
	cfg *models.RuntimeConfig,
) error {
	for _, cfgID := range cfg.Enabled {
		calcFunc := s.funcMap[cfgID]
		if calcFunc == nil {
			continue
		}

		err := calcFunc(ctx, message, info, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CalculatorService) calculateRuLetters(
	_ context.Context,
	message *bmodels.Message,
	info *models.MessageInfo,
	_ *models.RuntimeConfig,
) error {
	info.HasRULetters = s.IsTextContainsRuLetters(s.getText(message))

	return nil
}

//nolint:mnd
func (s *CalculatorService) calculateAntispam(
	_ context.Context,
	message *bmodels.Message,
	info *models.MessageInfo,
	_ *models.RuntimeConfig,
) error {
	text := s.getText(message)
	length := utf8.RuneCountInString(text)

	if length > 10 {
		upper := s.countTextUpper(text)
		info.HasCaps = float64(upper)/float64(length) > 0.5
	}

	for _, me := range slices.Concat(
		message.Entities,
		message.CaptionEntities,
	) {
		switch me.Type {
		case bmodels.MessageEntityTypeURL:
			info.CountLinks++
		case bmodels.MessageEntityTypeMention:
			info.CountMentions++
		}
	}

	info.CountMedias += s.countMedias(message)

	if info.CountMedias == 0 {
		info.HasShort = length < 5
	}

	info.HasLong = length > 1000

	if info.CountLinks > 0 {
		if message.LinkPreviewOptions != nil {
			if message.LinkPreviewOptions.IsDisabled == nil || !*message.LinkPreviewOptions.IsDisabled {
				info.CountEmbeds++
			}
		}
	}

	return nil
}
