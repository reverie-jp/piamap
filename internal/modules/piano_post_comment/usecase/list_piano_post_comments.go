package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/gateway"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

const (
	ListPianoPostCommentsDefaultLimit = 50
	ListPianoPostCommentsMaxLimit     = 200
)

type ListPianoPostCommentsParentKind int

const (
	ListParentPost ListPianoPostCommentsParentKind = iota + 1
	ListParentUser
)

type ListPianoPostCommentsInput struct {
	RequesterID  ulid.ULID
	ParentKind   ListPianoPostCommentsParentKind
	PianoPostID  ulid.ULID
	UserCustomID string
	PageSize     int32
	AfterID      *ulid.ULID
}

type ListPianoPostCommentsOutput struct {
	Views  []*gateway.PianoPostCommentView
	NextID *ulid.ULID
}

type ListPianoPostComments struct {
	gw          gateway.Gateway
	postGate    postgw.Gateway
	userGateway usergw.Gateway
}

func NewListPianoPostComments(gw gateway.Gateway, postGate postgw.Gateway, userGateway usergw.Gateway) *ListPianoPostComments {
	return &ListPianoPostComments{gw: gw, postGate: postGate, userGateway: userGateway}
}

func (uc *ListPianoPostComments) Execute(ctx context.Context, input ListPianoPostCommentsInput) (*ListPianoPostCommentsOutput, error) {
	limit := input.PageSize
	if limit <= 0 {
		limit = ListPianoPostCommentsDefaultLimit
	} else if limit > ListPianoPostCommentsMaxLimit {
		limit = ListPianoPostCommentsMaxLimit
	}
	queryLimit := limit + 1

	comments, err := uc.fetch(ctx, input, queryLimit)
	if err != nil {
		return nil, err
	}
	next, comments := splitNext(comments, int(limit))
	views, err := uc.gw.BuildListViews(ctx, input.RequesterID, comments)
	if err != nil {
		return nil, err
	}
	return &ListPianoPostCommentsOutput{Views: views, NextID: next}, nil
}

func (uc *ListPianoPostComments) fetch(ctx context.Context, input ListPianoPostCommentsInput, limit int32) ([]*entity.PianoPostComment, error) {
	switch input.ParentKind {
	case ListParentPost:
		if input.PianoPostID.IsZero() {
			return nil, xerrors.ErrInvalidArgument.WithMessage("post id is required")
		}
		// private 投稿の場合は作者本人のみ閲覧可。
		post, err := uc.postGate.GetPianoPost(ctx, input.PianoPostID)
		if err != nil {
			return nil, err
		}
		if post == nil {
			return nil, xerrors.ErrNotFound
		}
		if post.Visibility == entity.PostVisibilityPrivate && post.UserID != input.RequesterID {
			return nil, xerrors.ErrNotFound
		}
		return uc.gw.ListByPost(ctx, gateway.ListByPostParams{
			PianoPostID: input.PianoPostID,
			AfterID:     input.AfterID,
			Limit:       limit,
		})
	case ListParentUser:
		if input.UserCustomID == "" {
			return nil, xerrors.ErrInvalidArgument.WithMessage("user custom_id is required")
		}
		user, err := uc.userGateway.GetUserByCustomID(ctx, input.UserCustomID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, xerrors.ErrUserNotFound
		}
		return uc.gw.ListByUser(ctx, gateway.ListByUserParams{
			UserID:  user.ID,
			AfterID: input.AfterID,
			Limit:   limit,
		})
	}
	return nil, xerrors.ErrInvalidArgument.WithMessage("unknown parent kind")
}

func splitNext(items []*entity.PianoPostComment, limit int) (*ulid.ULID, []*entity.PianoPostComment) {
	if len(items) <= limit {
		return nil, items
	}
	cut := items[:limit]
	last := cut[len(cut)-1]
	id := last.ID
	return &id, cut
}
