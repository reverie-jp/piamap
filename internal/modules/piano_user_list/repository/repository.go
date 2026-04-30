package repository

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type ListByUserParams struct {
	UserID       ulid.ULID
	ListKind     entity.PianoListKind
	AfterPianoID *ulid.ULID
	Limit        int32
}

type Repository interface {
	UpsertList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error
	DeleteList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error
	ListByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error)
	ListKindsForPiano(ctx context.Context, userID, pianoID ulid.ULID) ([]entity.PianoListKind, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

func (r *RepositoryImpl) UpsertList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error {
	return r.q.UpsertPianoUserList(ctx, sqlc.UpsertPianoUserListParams{
		UserID:   userID,
		PianoID:  pianoID,
		ListKind: sqlc.PianoListKind(kind),
	})
}

func (r *RepositoryImpl) DeleteList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error {
	return r.q.DeletePianoUserList(ctx, sqlc.DeletePianoUserListParams{
		UserID:   userID,
		PianoID:  pianoID,
		ListKind: sqlc.PianoListKind(kind),
	})
}

func (r *RepositoryImpl) ListByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error) {
	rows, err := r.q.ListPianoUserListsByUser(ctx, sqlc.ListPianoUserListsByUserParams{
		UserID:       params.UserID,
		ListKind:     sqlc.PianoListKind(params.ListKind),
		AfterPianoID: params.AfterPianoID,
		LimitCount:   params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]ulid.ULID, len(rows))
	for i, r := range rows {
		out[i] = r.PianoID
	}
	return out, nil
}

func (r *RepositoryImpl) ListKindsForPiano(ctx context.Context, userID, pianoID ulid.ULID) ([]entity.PianoListKind, error) {
	rows, err := r.q.ListMyListKindsForPiano(ctx, sqlc.ListMyListKindsForPianoParams{
		UserID:  userID,
		PianoID: pianoID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]entity.PianoListKind, len(rows))
	for i, k := range rows {
		out[i] = entity.PianoListKind(k)
	}
	return out, nil
}
