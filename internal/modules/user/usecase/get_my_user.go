package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type GetMyUser struct {
	userGateway gateway.Gateway
}

func NewGetMyUser(userGateway gateway.Gateway) *GetMyUser {
	return &GetMyUser{userGateway: userGateway}
}

func (uc *GetMyUser) Execute(ctx context.Context, input GetMyUserInput) (*GetMyUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	view, err := uc.userGateway.BuildUserView(ctx, input.RequesterID, input.RequesterID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}
	return &GetMyUserOutput{View: view}, nil
}
