package migrations

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Migrate(
	ctx context.Context,
	dbOrig *gorm.DB,
) error {
	db := dbOrig.WithContext(ctx)

	err := db.AutoMigrate(
		&models.ChatConfig{},
		&models.MessageInfo{},
		&models.MessageDelete{},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
