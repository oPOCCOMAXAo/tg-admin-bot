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
		&models.Restriction{},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	err = dropIndexes(
		db.Migrator(),
		&models.MessageInfo{},
		"search_idx",
		"search_idx_v2",
	)
	if err != nil {
		return err
	}

	return nil
}

func dropIndexes(
	migrator gorm.Migrator,
	dst any,
	indexes ...string,
) error {
	for _, index := range indexes {
		if !migrator.HasIndex(dst, index) {
			continue
		}

		err := migrator.DropIndex(dst, index)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
