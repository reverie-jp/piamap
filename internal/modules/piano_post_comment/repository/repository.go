package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type InsertParams struct {
	ID              ulid.ULID
	PianoPostID     ulid.ULID
	UserID          ulid.ULID
	ParentCommentID *ulid.ULID
	Body            string
}

type ListByPostParams struct {
	PianoPostID ulid.ULID
	AfterID     *ulid.ULID
	Limit       int32
}

type ListByUserParams struct {
	UserID  ulid.ULID
	AfterID *ulid.ULID
	Limit   int32
}

type Repository interface {
	Insert(ctx context.Context, params InsertParams) error
	Get(ctx context.Context, id ulid.ULID) (*entity.PianoPostComment, error)
	Delete(ctx context.Context, id ulid.ULID) error
	ListByPost(ctx context.Context, params ListByPostParams) ([]*entity.PianoPostComment, error)
	ListByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPostComment, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

func (r *RepositoryImpl) Insert(ctx context.Context, params InsertParams) error {
	return r.q.InsertPianoPostComment(ctx, sqlc.InsertPianoPostCommentParams{
		ID:              params.ID,
		PianoPostID:     params.PianoPostID,
		UserID:          params.UserID,
		ParentCommentID: params.ParentCommentID,
		Body:            params.Body,
	})
}

func (r *RepositoryImpl) Get(ctx context.Context, id ulid.ULID) (*entity.PianoPostComment, error) {
	row, err := r.q.GetPianoPostComment(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rowToEntity(row), nil
}

func (r *RepositoryImpl) Delete(ctx context.Context, id ulid.ULID) error {
	return r.q.DeletePianoPostComment(ctx, id)
}

func (r *RepositoryImpl) ListByPost(ctx context.Context, params ListByPostParams) ([]*entity.PianoPostComment, error) {
	rows, err := r.q.ListPianoPostCommentsByPost(ctx, sqlc.ListPianoPostCommentsByPostParams{
		PianoPostID: params.PianoPostID,
		AfterID:     params.AfterID,
		LimitCount:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.PianoPostComment, len(rows))
	for i, row := range rows {
		out[i] = rowToEntity(row)
	}
	return out, nil
}

func (r *RepositoryImpl) ListByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPostComment, error) {
	rows, err := r.q.ListPianoPostCommentsByUser(ctx, sqlc.ListPianoPostCommentsByUserParams{
		UserID:     params.UserID,
		AfterID:    params.AfterID,
		LimitCount: params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.PianoPostComment, len(rows))
	for i, row := range rows {
		out[i] = rowToEntity(row)
	}
	return out, nil
}

func rowToEntity(row sqlc.PianoPostComment) *entity.PianoPostComment {
	return &entity.PianoPostComment{
		ID:              row.ID,
		PianoPostID:     row.PianoPostID,
		UserID:          row.UserID,
		ParentCommentID: row.ParentCommentID,
		Body:            row.Body,
		CreateTime:      row.CreateTime,
		UpdateTime:      row.UpdateTime,
	}
}
