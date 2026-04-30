package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/domain/mapper"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type CreateUserParams struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	AvatarURL   *string
}

type UpdateUserProfileParams struct {
	ID                ulid.ULID
	DisplayName       *string
	Biography         *string
	AvatarURL         *string
	Hometown          *string
	PianoStartedYear  *int16
	YearsOfExperience *int16
}

type Repository interface {
	GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error)
	GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error)
	ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	DeleteUser(ctx context.Context, id ulid.ULID) error
	UpdateUserProfile(ctx context.Context, params UpdateUserProfileParams) error
	UpdateUserCustomID(ctx context.Context, id ulid.ULID, customID string) error
	IsUserCurrentlyRestricted(ctx context.Context, id ulid.ULID) (bool, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

func (r *RepositoryImpl) GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToUser(&row), nil
}

func (r *RepositoryImpl) GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error) {
	row, err := r.q.GetUserByCustomID(ctx, customID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToUser(&row), nil
}

func (r *RepositoryImpl) ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error) {
	if len(ids) == 0 {
		return []*entity.User{}, nil
	}
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}
	rows, err := r.q.ListUsersByIDs(ctx, strIDs)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(rows))
	for i := range rows {
		users[i] = mapper.ToUser(&rows[i])
	}
	return users, nil
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, params CreateUserParams) error {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          params.ID,
		CustomID:    params.CustomID,
		DisplayName: params.DisplayName,
		AvatarUrl:   params.AvatarURL,
	})
}

func (r *RepositoryImpl) DeleteUser(ctx context.Context, id ulid.ULID) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *RepositoryImpl) UpdateUserProfile(ctx context.Context, params UpdateUserProfileParams) error {
	return r.q.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:                params.ID,
		DisplayName:       params.DisplayName,
		Biography:         params.Biography,
		AvatarUrl:         params.AvatarURL,
		Hometown:          params.Hometown,
		PianoStartedYear:  params.PianoStartedYear,
		YearsOfExperience: params.YearsOfExperience,
	})
}

func (r *RepositoryImpl) UpdateUserCustomID(ctx context.Context, id ulid.ULID, customID string) error {
	return r.q.UpdateUserCustomID(ctx, sqlc.UpdateUserCustomIDParams{
		ID:       id,
		CustomID: customID,
	})
}

func (r *RepositoryImpl) IsUserCurrentlyRestricted(ctx context.Context, id ulid.ULID) (bool, error) {
	return r.q.IsUserCurrentlyRestricted(ctx, id)
}
