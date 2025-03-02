package repo

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) GetOrCreateChatByTgID(
	ctx context.Context,
	tgID int64,
) (*models.ChatConfig, error) {
	var res models.ChatConfig

	err := r.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			var existing models.ChatConfig

			err := tx.
				Where("tg_id = ?", tgID).
				Take(&existing).
				Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WithStack(err)
			}

			if existing.ID != 0 {
				res = existing

				return nil
			}

			res.TgID = tgID
			res.CreatedAt = r.Now()
			res.UpdatedAt = res.CreatedAt

			err = tx.
				Create(&res).
				Error
			if err != nil {
				return errors.WithStack(err)
			}

			return nil
		})
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return &res, nil
}

func (r *Repo) UpdateChatConfigInt(
	ctx context.Context,
	tgChatID int64,
	cfgID models.ConfigID,
	value int64,
) error {
	column := r.chatConfig[cfgID]
	if column.Name == "" {
		return nil
	}

	err := r.db.WithContext(ctx).
		Model(&models.ChatConfig{}).
		Where("tg_id = ?", tgChatID).
		UpdateColumn(column.Name, column.ValueInt(value)).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetAllChatConfigs(
	ctx context.Context,
) ([]*models.ChatConfig, error) {
	var res []*models.ChatConfig

	err := r.db.WithContext(ctx).
		Find(&res).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}
