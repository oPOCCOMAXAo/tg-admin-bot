package domain

import (
	"time"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
)

const (
	ScoreBase         = 1
	ScoreCaps         = 1
	ScoreShort        = 1
	ScoreLong         = 1
	ScoreLinks        = 1
	ScoreEmbed        = 3
	ScoreMedia        = 4
	ScoreMention      = 3
	ScoreMentionExtra = 1
	ScoreFast         = 9
)

type Score uint16

func (s Score) UInt16() uint16 {
	return uint16(s)
}

func (s *Score) Add(score Score) {
	*s += score
}

func (s *Score) AddIf(cond bool, score Score) {
	if cond {
		s.Add(score)
	}
}

func (s *Score) AddGroup(count uint8, score Score) {
	s.Add(score * Score(count))
}

func CalculateScore(
	info *models.MessageInfo,
) uint16 {
	score := Score(ScoreBase)
	score.AddIf(info.HasCaps, ScoreCaps)
	score.AddIf(info.HasShort, ScoreShort)
	score.AddIf(info.HasLong, ScoreLong)
	score.AddIf(info.CountLinks > 0, ScoreLinks)
	score.AddGroup(info.CountEmbeds, ScoreEmbed)
	score.AddGroup(info.CountMedias, ScoreMedia)
	score.AddIf(info.CountMentions > 0, ScoreMention)
	score.AddGroup(info.CountMentions, ScoreMentionExtra)
	score.AddIf(info.IsFast, ScoreFast)

	return score.UInt16()
}

//nolint:mnd
func GetAntispamPenalties() []*models.AntispamPenalty {
	return []*models.AntispamPenalty{
		{
			Name:          "30s",
			CheckInterval: 30 * time.Second,
			MaxScore:      25,
			PenaltyTime:   5 * time.Minute,
		},
		{
			Name:          "1h",
			CheckInterval: 1 * time.Hour,
			MaxScore:      500,
			PenaltyTime:   1 * time.Hour,
		},
		{
			Name:          "24h",
			CheckInterval: 24 * time.Hour,
			MaxScore:      1000,
			PenaltyTime:   24 * time.Hour,
		},
	}
}
