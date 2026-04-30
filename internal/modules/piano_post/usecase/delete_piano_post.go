package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type DeletePianoPostInput struct {
	RequesterID ulid.ULID
	PostID      ulid.ULID
}

type DeletePianoPostOutput struct{}

type DeletePianoPost struct {
	gw gateway.Gateway
}

func NewDeletePianoPost(gw gateway.Gateway) *DeletePianoPost {
	return &DeletePianoPost{gw: gw}
}

func (uc *DeletePianoPost) Execute(ctx context.Context, input DeletePianoPostInput) (*DeletePianoPostOutput, error) {
	if input.RequesterID.IsZero() {
		return nil, xerrors.ErrUnauthenticated
	}
	if input.PostID.IsZero() {
		return nil, xerrors.ErrInvalidArgument.WithMessage("post id is required")
	}
	existing, err := uc.gw.GetPianoPost(ctx, input.PostID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, xerrors.ErrNotFound
	}
	if existing.UserID != input.RequesterID {
		return nil, xerrors.ErrPermissionDenied.WithMessage("not the author")
	}
	if err := uc.gw.DeletePianoPost(ctx, input.PostID); err != nil {
		return nil, err
	}
	return &DeletePianoPostOutput{}, nil
}
