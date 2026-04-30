package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/gateway"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type AddPianoToUserListInput struct {
	RequesterID ulid.ULID
	PianoID     ulid.ULID
	ListKind    entity.PianoListKind
}

func (i AddPianoToUserListInput) Validate() error {
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

type AddPianoToUserListOutput struct{}

type AddPianoToUserList struct {
	gw           gateway.Gateway
	pianoGateway pianogw.Gateway
}

func NewAddPianoToUserList(gw gateway.Gateway, pianoGateway pianogw.Gateway) *AddPianoToUserList {
	return &AddPianoToUserList{gw: gw, pianoGateway: pianoGateway}
}

func (uc *AddPianoToUserList) Execute(ctx context.Context, input AddPianoToUserListInput) (*AddPianoToUserListOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	piano, err := uc.pianoGateway.GetPiano(ctx, input.PianoID)
	if err != nil {
		return nil, err
	}
	if piano == nil || piano.Status == entity.PianoStatusRemoved {
		return nil, xerrors.ErrPianoNotFound
	}
	if err := uc.gw.UpsertList(ctx, input.RequesterID, input.PianoID, input.ListKind); err != nil {
		return nil, err
	}
	return &AddPianoToUserListOutput{}, nil
}
