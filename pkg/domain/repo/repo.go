package repo

import (
	"time"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB

	chatConfig map[models.ConfigID]models.ColumnConfig
}

func NewRepo(
	db *gorm.DB,
) *Repo {
	return &Repo{
		db: db,

		chatConfig: models.ChatConfig{}.Columns(),
	}
}

func (r *Repo) Now() int64 {
	return time.Now().Unix()
}
