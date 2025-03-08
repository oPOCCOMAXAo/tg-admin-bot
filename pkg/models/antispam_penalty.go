package models

import "time"

type AntispamPenalty struct {
	Name               string        `gorm:"-"`
	CheckInterval      time.Duration `gorm:"-"`
	MaxScore           int64         `gorm:"-"`
	PenaltyTimeSeconds int64         `gorm:"-"`
}
