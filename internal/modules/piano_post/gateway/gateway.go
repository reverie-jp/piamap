package gateway

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/repository"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// PianoPostView は PianoPost に表示用 author / piano メタを付与した型。
// タイムラインの N+1 を避けるため List 経由で hydrate する。
type PianoPostView struct {
	Post              *entity.PianoPost
	AuthorCustomID    string
	AuthorDisplayName string
	PianoDisplayName  string
	IsAuthor          bool
	ViewerLiked       bool
}

// LikeLookup は piano_post_like モジュールへの依存を最小限にする interface。
// requesterID が [postIDs] のうちどれをいいねしているか返す。
type LikeLookup interface {
	ListLikedPostIDsAmong(ctx context.Context, userID ulid.ULID, postIDs []ulid.ULID) ([]ulid.ULID, error)
}

type CreatePianoPostParams = repository.CreatePianoPostParams
type UpdatePianoPostParams = repository.UpdatePianoPostParams
type ListByPianoParams = repository.ListByPianoParams
type ListByUserParams = repository.ListByUserParams
type ListPublicParams = repository.ListPublicParams

type Gateway interface {
	GetPianoPost(ctx context.Context, id ulid.ULID) (*entity.PianoPost, error)
	ListPianoPostsByPiano(ctx context.Context, params ListByPianoParams) ([]*entity.PianoPost, error)
	ListPianoPostsByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPost, error)
	ListPublicPianoPosts(ctx context.Context, params ListPublicParams) ([]*entity.PianoPost, error)
	CreatePianoPost(ctx context.Context, params CreatePianoPostParams) error
	UpdatePianoPost(ctx context.Context, params UpdatePianoPostParams) error
	DeletePianoPost(ctx context.Context, id ulid.ULID) error
	UpsertPianoUserListVisited(ctx context.Context, userID, pianoID ulid.ULID) error

	BuildPianoPostView(ctx context.Context, requesterID ulid.ULID, post *entity.PianoPost) (*PianoPostView, error)
	BuildListPianoPostViews(ctx context.Context, requesterID ulid.ULID, posts []*entity.PianoPost) ([]*PianoPostView, error)
}

type gatewayImpl struct {
	repo         repository.Repository
	userGateway  usergw.Gateway
	pianoGateway pianogw.Gateway
	likeLookup   LikeLookup
}

func New(q sqlc.Querier, userGateway usergw.Gateway, pianoGateway pianogw.Gateway, likeLookup LikeLookup) Gateway {
	return &gatewayImpl{
		repo:         repository.New(q),
		userGateway:  userGateway,
		pianoGateway: pianoGateway,
		likeLookup:   likeLookup,
	}
}

func (g *gatewayImpl) GetPianoPost(ctx context.Context, id ulid.ULID) (*entity.PianoPost, error) {
	return g.repo.GetPianoPostByID(ctx, id)
}

func (g *gatewayImpl) ListPianoPostsByPiano(ctx context.Context, params ListByPianoParams) ([]*entity.PianoPost, error) {
	return g.repo.ListPianoPostsByPiano(ctx, params)
}

func (g *gatewayImpl) ListPianoPostsByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPost, error) {
	return g.repo.ListPianoPostsByUser(ctx, params)
}

func (g *gatewayImpl) ListPublicPianoPosts(ctx context.Context, params ListPublicParams) ([]*entity.PianoPost, error) {
	return g.repo.ListPublicPianoPosts(ctx, params)
}

func (g *gatewayImpl) CreatePianoPost(ctx context.Context, params CreatePianoPostParams) error {
	return g.repo.CreatePianoPost(ctx, params)
}

func (g *gatewayImpl) UpdatePianoPost(ctx context.Context, params UpdatePianoPostParams) error {
	return g.repo.UpdatePianoPost(ctx, params)
}

func (g *gatewayImpl) DeletePianoPost(ctx context.Context, id ulid.ULID) error {
	return g.repo.DeletePianoPost(ctx, id)
}

func (g *gatewayImpl) UpsertPianoUserListVisited(ctx context.Context, userID, pianoID ulid.ULID) error {
	return g.repo.UpsertPianoUserListVisited(ctx, userID, pianoID)
}

func (g *gatewayImpl) BuildPianoPostView(ctx context.Context, requesterID ulid.ULID, post *entity.PianoPost) (*PianoPostView, error) {
	views, err := g.BuildListPianoPostViews(ctx, requesterID, []*entity.PianoPost{post})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

// BuildListPianoPostViews は author / piano を 1 度ずつまとめて引いて N+1 を避ける。
func (g *gatewayImpl) BuildListPianoPostViews(ctx context.Context, requesterID ulid.ULID, posts []*entity.PianoPost) ([]*PianoPostView, error) {
	if len(posts) == 0 {
		return []*PianoPostView{}, nil
	}

	userIDs := make([]ulid.ULID, 0, len(posts))
	pianoIDs := make([]ulid.ULID, 0, len(posts))
	seenUser := make(map[string]bool, len(posts))
	seenPiano := make(map[string]bool, len(posts))
	for _, p := range posts {
		if p == nil {
			continue
		}
		if k := p.UserID.String(); !seenUser[k] {
			seenUser[k] = true
			userIDs = append(userIDs, p.UserID)
		}
		if k := p.PianoID.String(); !seenPiano[k] {
			seenPiano[k] = true
			pianoIDs = append(pianoIDs, p.PianoID)
		}
	}

	users, err := g.userGateway.ListUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	userByID := make(map[string]*entity.User, len(users))
	for _, u := range users {
		userByID[u.ID.String()] = u
	}

	pianoByID := make(map[string]*entity.Piano, len(pianoIDs))
	for _, id := range pianoIDs {
		piano, err := g.pianoGateway.GetPiano(ctx, id)
		if err != nil {
			return nil, err
		}
		if piano != nil {
			pianoByID[id.String()] = piano
		}
	}

	// 認証ユーザーのいいね済み post id を一括取得 (N+1 回避)。
	likedSet := map[string]bool{}
	if !requesterID.IsZero() && g.likeLookup != nil {
		postIDs := make([]ulid.ULID, 0, len(posts))
		for _, p := range posts {
			if p != nil {
				postIDs = append(postIDs, p.ID)
			}
		}
		liked, err := g.likeLookup.ListLikedPostIDsAmong(ctx, requesterID, postIDs)
		if err != nil {
			return nil, err
		}
		for _, id := range liked {
			likedSet[id.String()] = true
		}
	}

	views := make([]*PianoPostView, len(posts))
	for i, p := range posts {
		if p == nil {
			continue
		}
		view := &PianoPostView{
			Post:        p,
			IsAuthor:    !requesterID.IsZero() && p.UserID == requesterID,
			ViewerLiked: likedSet[p.ID.String()],
		}
		if u, ok := userByID[p.UserID.String()]; ok {
			view.AuthorCustomID = u.CustomID
			view.AuthorDisplayName = u.DisplayName
		}
		if pn, ok := pianoByID[p.PianoID.String()]; ok {
			view.PianoDisplayName = pn.Name
		}
		views[i] = view
	}
	return views, nil
}
