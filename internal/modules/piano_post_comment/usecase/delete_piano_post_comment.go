package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type DeletePianoPostCommentInput struct {
	RequesterID ulid.ULID
	CommentID   ulid.ULID
}

type DeletePianoPostCommentOutput struct{}

type DeletePianoPostComment struct {
	gw gateway.Gateway
}

func NewDeletePianoPostComment(gw gateway.Gateway) *DeletePianoPostComment {
	return &DeletePianoPostComment{gw: gw}
}

func (uc *DeletePianoPostComment) Execute(ctx context.Context, input DeletePianoPostCommentInput) (*DeletePianoPostCommentOutput, error) {
	if input.RequesterID.IsZero() {
		return nil, xerrors.ErrUnauthenticated
	}
	if input.CommentID.IsZero() {
		return nil, xerrors.ErrInvalidArgument.WithMessage("comment id is required")
	}
	existing, err := uc.gw.Get(ctx, input.CommentID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, xerrors.ErrNotFound
	}
	if existing.UserID != input.RequesterID {
		return nil, xerrors.ErrPermissionDenied.WithMessage("not the author")
	}
	if err := uc.gw.Delete(ctx, input.CommentID); err != nil {
		return nil, err
	}
	return &DeletePianoPostCommentOutput{}, nil
}
