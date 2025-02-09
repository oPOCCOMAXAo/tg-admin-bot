package migrations

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func Migrate(
	ctx context.Context,
	dbOrig *gorm.DB,
) error {
	db := dbOrig.WithContext(ctx)

	err := db.AutoMigrate()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
