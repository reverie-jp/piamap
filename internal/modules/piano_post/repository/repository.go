package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/domain/mapper"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type CreatePianoPostParams struct {
	ID             ulid.ULID
	UserID         ulid.ULID
	PianoID        ulid.ULID
	VisitTime      time.Time
	Rating         int16
	Body           *string
	AmbientNoise   *int16
	FootTraffic    *int16
	Resonance      *int16
	KeyTouchWeight *int16
	TuningQuality  *int16
	Visibility     entity.PostVisibility
}

type UpdatePianoPostParams struct {
	ID                ulid.ULID
	SetVisitTime      bool
	VisitTime         *time.Time
	SetRating         bool
	Rating            *int16
	SetBody           bool
	Body              *string
	SetAmbientNoise   bool
	AmbientNoise      *int16
	SetFootTraffic    bool
	FootTraffic       *int16
	SetResonance      bool
	Resonance         *int16
	SetKeyTouchWeight bool
	KeyTouchWeight    *int16
	SetTuningQuality  bool
	TuningQuality     *int16
	SetVisibility     bool
	Visibility        *entity.PostVisibility
}

type ListByPianoParams struct {
	PianoID ulid.ULID
	AfterID *ulid.ULID
	Limit   int32
}

type ListByUserParams struct {
	UserID         ulid.ULID
	IncludePrivate bool
	AfterID        *ulid.ULID
	Limit          int32
}

type ListPublicParams struct {
	AfterID *ulid.ULID
	Limit   int32
}

type Repository interface {
	GetPianoPostByID(ctx context.Context, id ulid.ULID) (*entity.PianoPost, error)
	ListPianoPostsByPiano(ctx context.Context, params ListByPianoParams) ([]*entity.PianoPost, error)
	ListPianoPostsByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPost, error)
	ListPublicPianoPosts(ctx context.Context, params ListPublicParams) ([]*entity.PianoPost, error)
	CreatePianoPost(ctx context.Context, params CreatePianoPostParams) error
	UpdatePianoPost(ctx context.Context, params UpdatePianoPostParams) error
	DeletePianoPost(ctx context.Context, id ulid.ULID) error
	UpsertPianoUserListVisited(ctx context.Context, userID, pianoID ulid.ULID) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

func (r *RepositoryImpl) GetPianoPostByID(ctx context.Context, id ulid.ULID) (*entity.PianoPost, error) {
	row, err := r.q.GetPianoPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToPianoPostFromGetRow(&row), nil
}

func (r *RepositoryImpl) ListPianoPostsByPiano(ctx context.Context, params ListByPianoParams) ([]*entity.PianoPost, error) {
	rows, err := r.q.ListPianoPostsByPiano(ctx, sqlc.ListPianoPostsByPianoParams{
		PianoID:    params.PianoID,
		AfterID:    params.AfterID,
		LimitCount: params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.PianoPost, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoPostFromListByPianoRow(&rows[i])
	}
	return out, nil
}

func (r *RepositoryImpl) ListPianoPostsByUser(ctx context.Context, params ListByUserParams) ([]*entity.PianoPost, error) {
	rows, err := r.q.ListPianoPostsByUser(ctx, sqlc.ListPianoPostsByUserParams{
		UserID:         params.UserID,
		IncludePrivate: params.IncludePrivate,
		AfterID:        params.AfterID,
		LimitCount:     params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.PianoPost, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoPostFromListByUserRow(&rows[i])
	}
	return out, nil
}

func (r *RepositoryImpl) ListPublicPianoPosts(ctx context.Context, params ListPublicParams) ([]*entity.PianoPost, error) {
	rows, err := r.q.ListPublicPianoPosts(ctx, sqlc.ListPublicPianoPostsParams{
		AfterID:    params.AfterID,
		LimitCount: params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.PianoPost, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoPostFromListPublicRow(&rows[i])
	}
	return out, nil
}

func (r *RepositoryImpl) CreatePianoPost(ctx context.Context, params CreatePianoPostParams) error {
	return r.q.CreatePianoPost(ctx, sqlc.CreatePianoPostParams{
		ID:             params.ID,
		UserID:         params.UserID,
		PianoID:        params.PianoID,
		VisitTime:      params.VisitTime,
		Rating:         params.Rating,
		Body:           params.Body,
		AmbientNoise:   params.AmbientNoise,
		FootTraffic:    params.FootTraffic,
		Resonance:      params.Resonance,
		KeyTouchWeight: params.KeyTouchWeight,
		TuningQuality:  params.TuningQuality,
		Visibility:     sqlc.PostVisibility(params.Visibility),
	})
}

func (r *RepositoryImpl) UpdatePianoPost(ctx context.Context, params UpdatePianoPostParams) error {
	return r.q.UpdatePianoPost(ctx, sqlc.UpdatePianoPostParams{
		ID:                params.ID,
		SetVisitTime:      params.SetVisitTime,
		VisitTime:         params.VisitTime,
		SetRating:         params.SetRating,
		Rating:            params.Rating,
		SetBody:           params.SetBody,
		Body:              params.Body,
		SetAmbientNoise:   params.SetAmbientNoise,
		AmbientNoise:      params.AmbientNoise,
		SetFootTraffic:    params.SetFootTraffic,
		FootTraffic:       params.FootTraffic,
		SetResonance:      params.SetResonance,
		Resonance:         params.Resonance,
		SetKeyTouchWeight: params.SetKeyTouchWeight,
		KeyTouchWeight:    params.KeyTouchWeight,
		SetTuningQuality:  params.SetTuningQuality,
		TuningQuality:     params.TuningQuality,
		SetVisibility:     params.SetVisibility,
		Visibility:        toNullPostVisibility(params.Visibility),
	})
}

func (r *RepositoryImpl) DeletePianoPost(ctx context.Context, id ulid.ULID) error {
	return r.q.DeletePianoPost(ctx, id)
}

func (r *RepositoryImpl) UpsertPianoUserListVisited(ctx context.Context, userID, pianoID ulid.ULID) error {
	return r.q.UpsertPianoUserListVisited(ctx, sqlc.UpsertPianoUserListVisitedParams{
		UserID:  userID,
		PianoID: pianoID,
	})
}

func toNullPostVisibility(v *entity.PostVisibility) sqlc.NullPostVisibility {
	if v == nil {
		return sqlc.NullPostVisibility{}
	}
	return sqlc.NullPostVisibility{PostVisibility: sqlc.PostVisibility(*v), Valid: true}
}
