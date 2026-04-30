package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type GetUser struct {
	userGateway gateway.Gateway
}

func NewGetUser(userGateway gateway.Gateway) *GetUser {
	return &GetUser{userGateway: userGateway}
}

func (uc *GetUser) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	target, err := uc.userGateway.GetUserByCustomID(ctx, input.TargetCustomID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, xerrors.ErrUserNotFound
	}
	view, err := uc.userGateway.BuildUserView(ctx, input.RequesterID, target.ID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}
	return &GetUserOutput{View: view}, nil
}
