package repo

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) CreateMessageDelete(
	ctx context.Context,
	value *models.MessageDelete,
) error {
	err := r.db.WithContext(ctx).
		Create(value).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetFirstMessageDeleteUntilTime(
	ctx context.Context,
	untilUnixTime int64,
) (*models.MessageDelete, error) {
	var res models.MessageDelete

	err := r.db.WithContext(ctx).
		Where("execute_at <= ?", untilUnixTime).
		Order("execute_at").
		First(&res).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(models.ErrNotFound)
		}

		return nil, errors.WithStack(err)
	}

	return &res, nil
}

func (r *Repo) GetFirstMessageDeleteAny(
	ctx context.Context,
) (*models.MessageDelete, error) {
	var res models.MessageDelete

	err := r.db.WithContext(ctx).
		Order("execute_at").
		First(&res).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(models.ErrNotFound)
		}

		return nil, errors.WithStack(err)
	}

	return &res, nil
}

func (r *Repo) DeleteMessageDelete(
	ctx context.Context,
	value *models.MessageDelete,
) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", value.ID).
		Delete(&models.MessageDelete{}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
