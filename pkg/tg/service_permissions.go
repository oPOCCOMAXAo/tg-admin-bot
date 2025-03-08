package tg

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	bmodels "github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

//nolint:gochecknoglobals
var memberTypesAdmin = []bmodels.ChatMemberType{
	bmodels.ChatMemberTypeAdministrator,
	bmodels.ChatMemberTypeOwner,
}

// HiddenAdminNickname is a nickname of sender when admin is hidden.
const HiddenAdminNickname = "groupanonymousbot"

//nolint:cyclop
func (s *Service) hasRequiredAdminPermissions(
	required *bmodels.ChatMemberAdministrator,
	chat *bmodels.ChatMember,
) bool {
	existing := chat.Administrator
	if existing == nil {
		existing = &bmodels.ChatMemberAdministrator{}
	}

	if required.CanChangeInfo && !existing.CanChangeInfo ||
		required.CanDeleteMessages && !existing.CanDeleteMessages ||
		required.CanInviteUsers && !existing.CanInviteUsers ||
		required.CanPinMessages && !existing.CanPinMessages ||
		required.CanRestrictMembers && !existing.CanRestrictMembers ||
		required.CanPromoteMembers && !existing.CanPromoteMembers {
		return false
	}

	return true
}

func (s *Service) GetChatMember(
	ctx context.Context,
	params *bot.GetChatMemberParams,
) (*bmodels.ChatMember, error) {
	if params.UserID == 0 {
		params.UserID = s.client.ID()
	}

	res, err := s.client.GetChatMember(ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

type CheckUserMemberPermissionsParams struct {
	ChatID        int64
	UserID        int64
	Username      string
	RequiredTypes []bmodels.ChatMemberType
	RequiredAdmin *bmodels.ChatMemberAdministrator
}

func (s *Service) CheckUserMemberPermissions(
	ctx context.Context,
	params *CheckUserMemberPermissionsParams,
) error {
	if params.ChatID == 0 {
		return errors.WithStack(models.ErrAuthInvalid)
	}

	var (
		chat *bmodels.ChatMember
		err  error
	)

	if strings.ToLower(params.Username) == HiddenAdminNickname {
		chat = &bmodels.ChatMember{
			Type: bmodels.ChatMemberTypeAdministrator,
		}
	} else {
		chat, err = s.GetChatMember(ctx, &bot.GetChatMemberParams{
			ChatID: params.ChatID,
			UserID: params.UserID,
		})
		if err != nil {
			return err
		}
	}

	if len(params.RequiredTypes) > 0 {
		if lo.IndexOf(params.RequiredTypes, chat.Type) == -1 {
			return errors.WithStack(models.ErrAuthInvalid)
		}
	}

	if params.RequiredAdmin != nil {
		if !s.hasRequiredAdminPermissions(params.RequiredAdmin, chat) {
			return errors.WithStack(models.ErrAuthInvalid)
		}
	}

	return nil
}

type MuteParams struct {
	ChatID       int64
	UserID       int64
	MuteDuration time.Duration
}

func (s *Service) MuteUser(
	ctx context.Context,
	params *MuteParams,
) error {
	_, err := s.client.RestrictChatMember(ctx, &bot.RestrictChatMemberParams{
		ChatID:      params.ChatID,
		UserID:      params.UserID,
		UntilDate:   int(time.Now().Add(params.MuteDuration).Unix()),
		Permissions: &bmodels.ChatPermissions{},
	})
	if err != nil {
		if errors.Is(err, bot.ErrorBadRequest) {
			s.logger.ErrorContext(ctx, "MuteUser",
				slog.Int64("chat_id", params.ChatID),
				slog.Int64("user_id", params.UserID),
				slog.Any("error", err),
			)

			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}

type BanChatParams struct {
	ChatID       int64
	SenderChatID int64
}

func (s *Service) BanSenderChat(
	ctx context.Context,
	params *BanChatParams,
) error {
	_, err := s.client.BanChatSenderChat(ctx, &bot.BanChatSenderChatParams{
		ChatID:       params.ChatID,
		SenderChatID: int(params.SenderChatID),
	})
	if err != nil {
		if errors.Is(err, bot.ErrorBadRequest) {
			s.logger.ErrorContext(ctx, "BanSenderChat",
				slog.Int64("chat_id", params.ChatID),
				slog.Int64("sender_chat_id", params.SenderChatID),
				slog.Any("error", err),
			)

			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) UnbanSenderChat(
	ctx context.Context,
	params *BanChatParams,
) error {
	_, err := s.client.UnbanChatSenderChat(ctx, &bot.UnbanChatSenderChatParams{
		ChatID:       params.ChatID,
		SenderChatID: int(params.SenderChatID),
	})
	if err != nil {
		if errors.Is(err, bot.ErrorBadRequest) {
			s.logger.ErrorContext(ctx, "UnbanSenderChat",
				slog.Int64("chat_id", params.ChatID),
				slog.Int64("sender_chat_id", params.SenderChatID),
				slog.Any("error", err),
			)

			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}
