package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type GetPianoPostInput struct {
	RequesterID ulid.ULID
	PostID      ulid.ULID
}

type GetPianoPostOutput struct {
	View *gateway.PianoPostView
}

type GetPianoPost struct {
	gw gateway.Gateway
}

func NewGetPianoPost(gw gateway.Gateway) *GetPianoPost {
	return &GetPianoPost{gw: gw}
}

func (uc *GetPianoPost) Execute(ctx context.Context, input GetPianoPostInput) (*GetPianoPostOutput, error) {
	if input.PostID.IsZero() {
		return nil, xerrors.ErrInvalidArgument.WithMessage("post id is required")
	}
	post, err := uc.gw.GetPianoPost(ctx, input.PostID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	if post.Visibility == entity.PostVisibilityPrivate {
		if input.RequesterID.IsZero() || post.UserID != input.RequesterID {
			return nil, xerrors.ErrNotFound
		}
	}
	view, err := uc.gw.BuildPianoPostView(ctx, input.RequesterID, post)
	if err != nil {
		return nil, err
	}
	return &GetPianoPostOutput{View: view}, nil
}
