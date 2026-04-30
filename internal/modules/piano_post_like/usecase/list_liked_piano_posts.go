package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/gateway"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

const (
	ListLikedPianoPostsDefaultLimit = 20
	ListLikedPianoPostsMaxLimit     = 100
)

type ListLikedPianoPostsInput struct {
	RequesterID  ulid.ULID
	UserCustomID string
	PageSize     int32
	AfterPostID  *ulid.ULID
}

func (i ListLikedPianoPostsInput) Validate() error {
	if i.UserCustomID == "" {
		return xerrors.ErrInvalidArgument.WithMessage("user custom_id is required")
	}
	return nil
}

type ListLikedPianoPostsOutput struct {
	Views      []*postgw.PianoPostView
	NextPostID *ulid.ULID
}

type ListLikedPianoPosts struct {
	gw          gateway.Gateway
	userGateway usergw.Gateway
	postGate    postgw.Gateway
}

func NewListLikedPianoPosts(gw gateway.Gateway, userGateway usergw.Gateway, postGate postgw.Gateway) *ListLikedPianoPosts {
	return &ListLikedPianoPosts{gw: gw, userGateway: userGateway, postGate: postGate}
}

func (uc *ListLikedPianoPosts) Execute(ctx context.Context, input ListLikedPianoPostsInput) (*ListLikedPianoPostsOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	user, err := uc.userGateway.GetUserByCustomID(ctx, input.UserCustomID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}
	limit := input.PageSize
	if limit <= 0 {
		limit = ListLikedPianoPostsDefaultLimit
	} else if limit > ListLikedPianoPostsMaxLimit {
		limit = ListLikedPianoPostsMaxLimit
	}
	queryLimit := limit + 1

	postIDs, err := uc.gw.ListLikedByUser(ctx, gateway.ListByUserParams{
		UserID:      user.ID,
		AfterPostID: input.AfterPostID,
		Limit:       queryLimit,
	})
	if err != nil {
		return nil, err
	}
	var nextID *ulid.ULID
	if len(postIDs) > int(limit) {
		postIDs = postIDs[:limit]
		last := postIDs[len(postIDs)-1]
		nextID = &last
	}

	posts := make([]*entity.PianoPost, 0, len(postIDs))
	for _, id := range postIDs {
		p, err := uc.postGate.GetPianoPost(ctx, id)
		if err != nil {
			return nil, err
		}
		if p == nil {
			continue
		}
		// private 投稿は作者本人のみ閲覧可。
		if p.Visibility == entity.PostVisibilityPrivate && p.UserID != input.RequesterID {
			continue
		}
		posts = append(posts, p)
	}
	views, err := uc.postGate.BuildListPianoPostViews(ctx, input.RequesterID, posts)
	if err != nil {
		return nil, err
	}
	return &ListLikedPianoPostsOutput{Views: views, NextPostID: nextID}, nil
}
