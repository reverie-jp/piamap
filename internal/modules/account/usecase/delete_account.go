package usecase

import (
	"context"

	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type DeleteAccount struct {
	userGateway usergw.Gateway
}

func NewDeleteAccount(userGateway usergw.Gateway) *DeleteAccount {
	return &DeleteAccount{userGateway: userGateway}
}

func (uc *DeleteAccount) Execute(ctx context.Context, input DeleteAccountInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	user, err := uc.userGateway.GetUserByID(ctx, input.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return xerrors.ErrAccountNotFound
	}
	if user.CustomID != input.ConfirmCustomID {
		return xerrors.ErrCustomIDMismatch
	}
	return uc.userGateway.DeleteUser(ctx, input.UserID)
}
