package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

const (
	ListPianoPostsDefaultLimit = 20
	ListPianoPostsMaxLimit     = 100
)

type ListPianoPostsParentKind int

const (
	ListParentGlobal ListPianoPostsParentKind = iota
	ListParentPiano
	ListParentUser
)

type ListPianoPostsInput struct {
	RequesterID ulid.ULID
	ParentKind  ListPianoPostsParentKind
	PianoID     ulid.ULID // ParentKind=Piano のとき
	UserCustom  string    // ParentKind=User のとき
	PageSize    int32
	AfterID     *ulid.ULID
}

type ListPianoPostsOutput struct {
	Views  []*gateway.PianoPostView
	NextID *ulid.ULID
}

type ListPianoPosts struct {
	gw          gateway.Gateway
	userGateway usergw.Gateway
}

func NewListPianoPosts(gw gateway.Gateway, userGateway usergw.Gateway) *ListPianoPosts {
	return &ListPianoPosts{gw: gw, userGateway: userGateway}
}

func (uc *ListPianoPosts) Execute(ctx context.Context, input ListPianoPostsInput) (*ListPianoPostsOutput, error) {
	limit := input.PageSize
	if limit <= 0 {
		limit = ListPianoPostsDefaultLimit
	} else if limit > ListPianoPostsMaxLimit {
		limit = ListPianoPostsMaxLimit
	}
	queryLimit := limit + 1

	posts, err := uc.fetch(ctx, input, queryLimit)
	if err != nil {
		return nil, err
	}
	next, posts := splitNext(posts, int(limit))
	views, err := uc.gw.BuildListPianoPostViews(ctx, input.RequesterID, posts)
	if err != nil {
		return nil, err
	}
	return &ListPianoPostsOutput{Views: views, NextID: next}, nil
}

func (uc *ListPianoPosts) fetch(ctx context.Context, input ListPianoPostsInput, limit int32) ([]*entity.PianoPost, error) {
	switch input.ParentKind {
	case ListParentPiano:
		if input.PianoID.IsZero() {
			return nil, xerrors.ErrInvalidArgument.WithMessage("piano id is required")
		}
		return uc.gw.ListPianoPostsByPiano(ctx, gateway.ListByPianoParams{
			PianoID: input.PianoID,
			AfterID: input.AfterID,
			Limit:   limit,
		})
	case ListParentUser:
		if input.UserCustom == "" {
			return nil, xerrors.ErrInvalidArgument.WithMessage("user custom_id is required")
		}
		user, err := uc.userGateway.GetUserByCustomID(ctx, input.UserCustom)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, xerrors.ErrUserNotFound
		}
		includePrivate := !input.RequesterID.IsZero() && user.ID == input.RequesterID
		return uc.gw.ListPianoPostsByUser(ctx, gateway.ListByUserParams{
			UserID:         user.ID,
			IncludePrivate: includePrivate,
			AfterID:        input.AfterID,
			Limit:          limit,
		})
	case ListParentGlobal:
		return uc.gw.ListPublicPianoPosts(ctx, gateway.ListPublicParams{
			AfterID: input.AfterID,
			Limit:   limit,
		})
	}
	return nil, xerrors.ErrInvalidArgument.WithMessage("unknown parent kind")
}

// splitNext は queryLimit (= limit+1) で取得した結果から、超過分を切り落として
// next_page_token 用の id を返す。サイズが limit 以下なら next は無し。
func splitNext(items []*entity.PianoPost, limit int) (*ulid.ULID, []*entity.PianoPost) {
	if len(items) <= limit {
		return nil, items
	}
	cut := items[:limit]
	last := cut[len(cut)-1]
	id := last.ID
	return &id, cut
}
