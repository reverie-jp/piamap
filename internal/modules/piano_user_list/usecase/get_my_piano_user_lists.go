package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type GetMyPianoUserListsInput struct {
	RequesterID ulid.ULID
	PianoID     ulid.ULID
}

func (i GetMyPianoUserListsInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PianoID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("piano id is required")
	}
	return nil
}

type GetMyPianoUserListsOutput struct {
	ListKinds []entity.PianoListKind
}

type GetMyPianoUserLists struct {
	gw gateway.Gateway
}

func NewGetMyPianoUserLists(gw gateway.Gateway) *GetMyPianoUserLists {
	return &GetMyPianoUserLists{gw: gw}
}

func (uc *GetMyPianoUserLists) Execute(ctx context.Context, input GetMyPianoUserListsInput) (*GetMyPianoUserListsOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	kinds, err := uc.gw.ListKindsForPiano(ctx, input.RequesterID, input.PianoID)
	if err != nil {
		return nil, err
	}
	return &GetMyPianoUserListsOutput{ListKinds: kinds}, nil
}
