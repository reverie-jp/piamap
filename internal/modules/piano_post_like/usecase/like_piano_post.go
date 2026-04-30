package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/gateway"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type LikePianoPostInput struct {
	RequesterID ulid.ULID
	PostID      ulid.ULID
}

func (i LikePianoPostInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PostID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("post id is required")
	}
	return nil
}

type LikePianoPostOutput struct{}

type LikePianoPost struct {
	gw       gateway.Gateway
	postGate postgw.Gateway
}

func NewLikePianoPost(gw gateway.Gateway, postGate postgw.Gateway) *LikePianoPost {
	return &LikePianoPost{gw: gw, postGate: postGate}
}

func (uc *LikePianoPost) Execute(ctx context.Context, input LikePianoPostInput) (*LikePianoPostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	post, err := uc.postGate.GetPianoPost(ctx, input.PostID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	// private 投稿は作者本人のみがいいねできる (実用上ないがガード)。
	if post.Visibility == entity.PostVisibilityPrivate && post.UserID != input.RequesterID {
		return nil, xerrors.ErrNotFound
	}
	if err := uc.gw.UpsertLike(ctx, input.RequesterID, input.PostID); err != nil {
		return nil, err
	}
	return &LikePianoPostOutput{}, nil
}
