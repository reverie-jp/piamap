package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type UnlikePianoPostInput struct {
	RequesterID ulid.ULID
	PostID      ulid.ULID
}

func (i UnlikePianoPostInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PostID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("post id is required")
	}
	return nil
}

type UnlikePianoPostOutput struct{}

type UnlikePianoPost struct {
	gw gateway.Gateway
}

func NewUnlikePianoPost(gw gateway.Gateway) *UnlikePianoPost {
	return &UnlikePianoPost{gw: gw}
}

func (uc *UnlikePianoPost) Execute(ctx context.Context, input UnlikePianoPostInput) (*UnlikePianoPostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := uc.gw.DeleteLike(ctx, input.RequesterID, input.PostID); err != nil {
		return nil, err
	}
	return &UnlikePianoPostOutput{}, nil
}
