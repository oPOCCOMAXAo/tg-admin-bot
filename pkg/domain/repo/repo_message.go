package repo

import (
	"context"

	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *Repo) CreateMessageInfo(
	ctx context.Context,
	value *models.MessageInfo,
) error {
	err := r.db.WithContext(ctx).
		Create(value).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetMessageInfoForProcessing(
	ctx context.Context,
) (*models.MessageInfo, error) {
	var res models.MessageInfo

	err := r.db.WithContext(ctx).
		Model(&models.MessageInfo{}).
		Where("is_processed = 0").
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

func (r *Repo) UpdateMessageInfo(
	ctx context.Context,
	value *models.MessageInfo,
) error {
	err := r.db.WithContext(ctx).
		Model(&models.MessageInfo{}).
		Where("id = ?", value.ID).
		Select("*").
		Updates(value).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetMessageInfoPrevious(
	ctx context.Context,
	value *models.MessageInfo,
) (*models.MessageInfo, error) {
	const maxTimeDiff = 60

	var res models.MessageInfo

	err := r.db.WithContext(ctx).
		Model(&models.MessageInfo{}).
		Where("time <= ?", value.Time).
		Where("time >= ?", value.Time-maxTimeDiff).
		Where("chat_id = ?", value.ChatID).
		Where("message_id < ?", value.MessageID).
		Where("user_id = ?", value.UserID).
		Where("sender_chat_id = ?", value.SenderChatID).
		Where("id != ?", value.ID).
		Order("time DESC").
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

func (r *Repo) GetMessagePrevScore(
	ctx context.Context,
	value *models.MessageInfo,
	timeFrom, timeTo int64,
) (int64, error) {
	var res struct {
		Score int64 `gorm:"column:score"`
	}

	err := r.db.WithContext(ctx).
		Select("SUM(score) AS score").
		Table("message_info").
		Where("user_id = ?", value.UserID).
		Where("sender_chat_id = ?", value.SenderChatID).
		Where("time >= ?", timeFrom).
		Where("time <= ?", timeTo).
		Where("id != ?", value.ID).
		Take(&res).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.WithStack(err)
	}

	return res.Score, nil
}
