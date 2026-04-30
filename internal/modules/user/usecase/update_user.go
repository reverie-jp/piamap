package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type UpdateUser struct {
	userGateway gateway.Gateway
}

func NewUpdateUser(userGateway gateway.Gateway) *UpdateUser {
	return &UpdateUser{userGateway: userGateway}
}

func (uc *UpdateUser) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	if input.SetCustomID {
		existing, err := uc.userGateway.GetUserByCustomID(ctx, input.CustomID)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != input.RequesterID {
			return nil, xerrors.ErrUserCustomIDInUse
		}
		if existing == nil || existing.ID != input.RequesterID {
			if err := uc.userGateway.UpdateUserCustomID(ctx, input.RequesterID, input.CustomID); err != nil {
				return nil, err
			}
		}
	}

	if input.SetDisplayName || input.SetBiography || input.SetAvatarURL ||
		input.SetHometown || input.SetPianoStartedYear || input.SetYearsOfExperience {
		params := gateway.UpdateUserProfileParams{ID: input.RequesterID}
		if input.SetDisplayName {
			dn := input.DisplayName
			params.DisplayName = &dn
		}
		if input.SetBiography {
			params.Biography = input.Biography
		}
		if input.SetAvatarURL {
			params.AvatarURL = input.AvatarURL
		}
		if input.SetHometown {
			params.Hometown = input.Hometown
		}
		if input.SetPianoStartedYear {
			params.PianoStartedYear = input.PianoStartedYear
		}
		if input.SetYearsOfExperience {
			params.YearsOfExperience = input.YearsOfExperience
		}
		if err := uc.userGateway.UpdateUserProfile(ctx, params); err != nil {
			return nil, err
		}
	}

	view, err := uc.userGateway.BuildUserView(ctx, input.RequesterID, input.RequesterID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}
	return &UpdateUserOutput{View: view}, nil
}
