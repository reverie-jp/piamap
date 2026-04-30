package repository

import (
	"context"

	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type ListByUserParams struct {
	UserID      ulid.ULID
	AfterPostID *ulid.ULID
	Limit       int32
}

type Repository interface {
	UpsertLike(ctx context.Context, userID, postID ulid.ULID) error
	DeleteLike(ctx context.Context, userID, postID ulid.ULID) error
	ListLikedByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error)
	ListLikedPostIDsAmong(ctx context.Context, userID ulid.ULID, postIDs []ulid.ULID) ([]ulid.ULID, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

func (r *RepositoryImpl) UpsertLike(ctx context.Context, userID, postID ulid.ULID) error {
	return r.q.UpsertPianoPostLike(ctx, sqlc.UpsertPianoPostLikeParams{
		UserID:      userID,
		PianoPostID: postID,
	})
}

func (r *RepositoryImpl) DeleteLike(ctx context.Context, userID, postID ulid.ULID) error {
	return r.q.DeletePianoPostLike(ctx, sqlc.DeletePianoPostLikeParams{
		UserID:      userID,
		PianoPostID: postID,
	})
}

func (r *RepositoryImpl) ListLikedByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error) {
	rows, err := r.q.ListLikedPianoPostIDsByUser(ctx, sqlc.ListLikedPianoPostIDsByUserParams{
		UserID:      params.UserID,
		AfterPostID: params.AfterPostID,
		LimitCount:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]ulid.ULID, len(rows))
	for i, r := range rows {
		out[i] = r.PianoPostID
	}
	return out, nil
}

func (r *RepositoryImpl) ListLikedPostIDsAmong(ctx context.Context, userID ulid.ULID, postIDs []ulid.ULID) ([]ulid.ULID, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}
	idStrs := make([]string, len(postIDs))
	for i, id := range postIDs {
		idStrs[i] = id.String()
	}
	return r.q.ListLikedPostIDsForUserAndPosts(ctx, sqlc.ListLikedPostIDsForUserAndPostsParams{
		UserID:       userID,
		PianoPostIds: idStrs,
	})
}
