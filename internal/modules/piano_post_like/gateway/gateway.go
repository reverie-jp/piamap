package gateway

import (
	"context"

	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/repository"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type ListByUserParams = repository.ListByUserParams

type Gateway interface {
	UpsertLike(ctx context.Context, userID, postID ulid.ULID) error
	DeleteLike(ctx context.Context, userID, postID ulid.ULID) error
	ListLikedByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error)
	ListLikedPostIDsAmong(ctx context.Context, userID ulid.ULID, postIDs []ulid.ULID) ([]ulid.ULID, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}

func (g *gatewayImpl) UpsertLike(ctx context.Context, userID, postID ulid.ULID) error {
	return g.repo.UpsertLike(ctx, userID, postID)
}

func (g *gatewayImpl) DeleteLike(ctx context.Context, userID, postID ulid.ULID) error {
	return g.repo.DeleteLike(ctx, userID, postID)
}

func (g *gatewayImpl) ListLikedByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error) {
	return g.repo.ListLikedByUser(ctx, params)
}

func (g *gatewayImpl) ListLikedPostIDsAmong(ctx context.Context, userID ulid.ULID, postIDs []ulid.ULID) ([]ulid.ULID, error) {
	return g.repo.ListLikedPostIDsAmong(ctx, userID, postIDs)
}
