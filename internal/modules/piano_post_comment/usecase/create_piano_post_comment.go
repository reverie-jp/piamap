package usecase

import (
	"context"
	"strings"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/gateway"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

const CommentBodyMaxLen = 2000

type CreatePianoPostCommentInput struct {
	RequesterID     ulid.ULID
	PianoPostID     ulid.ULID
	ParentCommentID *ulid.ULID
	Body            string
}

func (i CreatePianoPostCommentInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PianoPostID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("post id is required")
	}
	body := strings.TrimSpace(i.Body)
	if body == "" {
		return xerrors.ErrInvalidArgument.WithMessage("body is required")
	}
	if len(body) > CommentBodyMaxLen {
		return xerrors.ErrInvalidArgument.WithMessage("body too long")
	}
	return nil
}

type CreatePianoPostCommentOutput struct {
	View *gateway.PianoPostCommentView
}

type CreatePianoPostComment struct {
	gw       gateway.Gateway
	postGate postgw.Gateway
}

func NewCreatePianoPostComment(gw gateway.Gateway, postGate postgw.Gateway) *CreatePianoPostComment {
	return &CreatePianoPostComment{gw: gw, postGate: postGate}
}

func (uc *CreatePianoPostComment) Execute(ctx context.Context, input CreatePianoPostCommentInput) (*CreatePianoPostCommentOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	post, err := uc.postGate.GetPianoPost(ctx, input.PianoPostID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, xerrors.ErrNotFound
	}
	// private 投稿には作者本人のみコメント可。
	if post.Visibility == entity.PostVisibilityPrivate && post.UserID != input.RequesterID {
		return nil, xerrors.ErrNotFound
	}

	if input.ParentCommentID != nil && !input.ParentCommentID.IsZero() {
		parent, err := uc.gw.Get(ctx, *input.ParentCommentID)
		if err != nil {
			return nil, err
		}
		if parent == nil || parent.PianoPostID != input.PianoPostID {
			return nil, xerrors.ErrInvalidArgument.WithMessage("parent_comment does not belong to this post")
		}
	}

	commentID := ulid.New()
	body := strings.TrimSpace(input.Body)
	if err := uc.gw.Insert(ctx, gateway.InsertParams{
		ID:              commentID,
		PianoPostID:     input.PianoPostID,
		UserID:          input.RequesterID,
		ParentCommentID: input.ParentCommentID,
		Body:            body,
	}); err != nil {
		return nil, err
	}
	created, err := uc.gw.Get(ctx, commentID)
	if err != nil {
		return nil, err
	}
	view, err := uc.gw.BuildView(ctx, input.RequesterID, created)
	if err != nil {
		return nil, err
	}
	return &CreatePianoPostCommentOutput{View: view}, nil
}
