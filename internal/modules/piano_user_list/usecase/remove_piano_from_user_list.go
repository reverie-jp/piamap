package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type RemovePianoFromUserListInput struct {
	RequesterID ulid.ULID
	PianoID     ulid.ULID
	ListKind    entity.PianoListKind
}

func (i RemovePianoFromUserListInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PianoID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("piano id is required")
	}
	if !i.ListKind.Valid() {
		return xerrors.ErrInvalidArgument.WithMessage("invalid list_kind")
	}
	return nil
}

type RemovePianoFromUserListOutput struct{}

type RemovePianoFromUserList struct {
	gw gateway.Gateway
}

func NewRemovePianoFromUserList(gw gateway.Gateway) *RemovePianoFromUserList {
	return &RemovePianoFromUserList{gw: gw}
}

func (uc *RemovePianoFromUserList) Execute(ctx context.Context, input RemovePianoFromUserListInput) (*RemovePianoFromUserListOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := uc.gw.DeleteList(ctx, input.RequesterID, input.PianoID, input.ListKind); err != nil {
		return nil, err
	}
	return &RemovePianoFromUserListOutput{}, nil
}
