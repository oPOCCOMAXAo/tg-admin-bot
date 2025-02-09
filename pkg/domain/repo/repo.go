package repo

import (
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(
	db *gorm.DB,
) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Now() int64 {
	return time.Now().Unix()
}
