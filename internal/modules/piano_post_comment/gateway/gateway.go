package gateway

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/repository"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type InsertParams = repository.InsertParams
type ListByPostParams = repository.ListByPostParams
type ListByUserParams = repository.ListByUserParams

// PianoPostCommentView はコメントに表示用 author メタと piano_id を付けた型。
// resource name (pianos/.../posts/.../comments/...) を組み立てるために PianoID が必要。
type PianoPostCommentView struct {
	Comment           *entity.PianoPostComment
	PianoID           ulid.ULID
	AuthorCustomID    string
	AuthorDisplayName string
	AuthorAvatarURL   *string
	IsAuthor          bool
}

type Gateway interface {
	Insert(ctx context.Context, params InsertParams) error
	Get(ctx context.Context, id ulid.ULID) (*entity.PianoPostComment, error)
	Delete(ctx context.Context, id ulid.ULID) error
	ListByPost(ctx context.Context, params ListByPostParams) ([]*entity.PianoPostComment, error)
	ListByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPostComment, error)

	BuildView(ctx context.Context, requesterID ulid.ULID, comment *entity.PianoPostComment) (*PianoPostCommentView, error)
	BuildListViews(ctx context.Context, requesterID ulid.ULID, comments []*entity.PianoPostComment) ([]*PianoPostCommentView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
	postGateway postgw.Gateway
}

func New(q sqlc.Querier, userGateway usergw.Gateway, postGateway postgw.Gateway) Gateway {
	return &gatewayImpl{repo: repository.New(q), userGateway: userGateway, postGateway: postGateway}
}

func (g *gatewayImpl) Insert(ctx context.Context, params InsertParams) error {
	return g.repo.Insert(ctx, params)
}

func (g *gatewayImpl) Get(ctx context.Context, id ulid.ULID) (*entity.PianoPostComment, error) {
	return g.repo.Get(ctx, id)
}

func (g *gatewayImpl) Delete(ctx context.Context, id ulid.ULID) error {
	return g.repo.Delete(ctx, id)
}

func (g *gatewayImpl) ListByPost(ctx context.Context, params ListByPostParams) ([]*entity.PianoPostComment, error) {
	return g.repo.ListByPost(ctx, params)
}

func (g *gatewayImpl) ListByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPostComment, error) {
	return g.repo.ListByUser(ctx, params)
}

func (g *gatewayImpl) BuildView(ctx context.Context, requesterID ulid.ULID, comment *entity.PianoPostComment) (*PianoPostCommentView, error) {
	views, err := g.BuildListViews(ctx, requesterID, []*entity.PianoPostComment{comment})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

func (g *gatewayImpl) BuildListViews(ctx context.Context, requesterID ulid.ULID, comments []*entity.PianoPostComment) ([]*PianoPostCommentView, error) {
	if len(comments) == 0 {
		return []*PianoPostCommentView{}, nil
	}
	userIDs := make([]ulid.ULID, 0, len(comments))
	postIDs := make([]ulid.ULID, 0, len(comments))
	seenUser := make(map[string]bool, len(comments))
	seenPost := make(map[string]bool, len(comments))
	for _, c := range comments {
		if c == nil {
			continue
		}
		if k := c.UserID.String(); !seenUser[k] {
			seenUser[k] = true
			userIDs = append(userIDs, c.UserID)
		}
		if k := c.PianoPostID.String(); !seenPost[k] {
			seenPost[k] = true
			postIDs = append(postIDs, c.PianoPostID)
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
	pianoByPostID := make(map[string]ulid.ULID, len(postIDs))
	for _, id := range postIDs {
		post, err := g.postGateway.GetPianoPost(ctx, id)
		if err != nil {
			return nil, err
		}
		if post != nil {
			pianoByPostID[id.String()] = post.PianoID
		}
	}
	views := make([]*PianoPostCommentView, len(comments))
	for i, c := range comments {
		if c == nil {
			continue
		}
		view := &PianoPostCommentView{
			Comment:  c,
			IsAuthor: !requesterID.IsZero() && c.UserID == requesterID,
		}
		if pid, ok := pianoByPostID[c.PianoPostID.String()]; ok {
			view.PianoID = pid
		}
		if u, ok := userByID[c.UserID.String()]; ok {
			view.AuthorCustomID = u.CustomID
			view.AuthorDisplayName = u.DisplayName
			view.AuthorAvatarURL = u.AvatarURL
		}
		views[i] = view
	}
	return views, nil
}
