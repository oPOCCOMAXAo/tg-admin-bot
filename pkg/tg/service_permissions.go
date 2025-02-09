package tg

import (
	"context"

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

	chat, err := s.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: params.ChatID,
		UserID: params.UserID,
	})
	if err != nil {
		return err
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
