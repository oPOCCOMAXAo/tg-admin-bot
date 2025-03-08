package repo

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) CreateRestriction(
	ctx context.Context,
	value *models.Restriction,
) error {
	err := r.db.WithContext(ctx).
		Create(value).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetRestrictionForProcessing(
	ctx context.Context,
	now int64,
) (*models.Restriction, error) {
	var res models.Restriction

	err := r.db.WithContext(ctx).
		Model(&models.Restriction{}).
		Where("execute_at <= ?", now).
		Order("execute_at ASC").
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

func (r *Repo) GetFirstRestrictionAny(
	ctx context.Context,
) (*models.Restriction, error) {
	var res models.Restriction

	err := r.db.WithContext(ctx).
		Order("execute_at ASC").
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

func (r *Repo) DeleteRestrictionByID(
	ctx context.Context,
	id int64,
) error {
	err := r.db.WithContext(ctx).
		Delete(&models.Restriction{ID: id}).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetLastUnbanExecuteAt(
	ctx context.Context,
	value *models.Restriction,
) (*models.Restriction, error) {
	var res models.Restriction

	err := r.db.WithContext(ctx).
		Model(&models.Restriction{}).
		Where("chat_id = ?", value.ChatID).
		Where("user_id = ?", value.UserID).
		Where("sender_chat_id = ?", value.SenderChatID).
		Where("is_unban = 1").
		Order("execute_at DESC").
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
