package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/domain/mapper"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type CreateAuthProviderParams struct {
	UserID         ulid.ULID
	Provider       string
	ProviderUserID string
}

type CreateRefreshTokenParams struct {
	UserID     ulid.ULID
	RawToken   string
	ExpireTime time.Time
}

type AuthProvider struct {
	UserID ulid.ULID
}

type Repository interface {
	GetAuthProviderByProvider(ctx context.Context, provider, providerUserID string) (*AuthProvider, error)
	CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error
	CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) error
	GetRefreshTokenByRaw(ctx context.Context, raw string) (*entity.RefreshToken, error)
	DeleteRefreshTokenByRaw(ctx context.Context, raw string, userID ulid.ULID) error
	DeleteExpiredRefreshTokensByUserID(ctx context.Context, userID ulid.ULID) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

// hashRefreshToken: 生 refresh token は DB に保存しない (digest のみ保存)。
func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (r *RepositoryImpl) GetAuthProviderByProvider(ctx context.Context, provider, providerUserID string) (*AuthProvider, error) {
	row, err := r.q.GetAuthProviderByProvider(ctx, sqlc.GetAuthProviderByProviderParams{
		Provider:       sqlc.AuthProvider(provider),
		ProviderUserID: providerUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &AuthProvider{UserID: row.UserID}, nil
}

func (r *RepositoryImpl) CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error {
	return r.q.CreateAuthProvider(ctx, sqlc.CreateAuthProviderParams{
		ID:             ulid.New(),
		UserID:         params.UserID,
		Provider:       sqlc.AuthProvider(params.Provider),
		ProviderUserID: params.ProviderUserID,
	})
}

func (r *RepositoryImpl) CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) error {
	return r.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID:         ulid.New(),
		UserID:     params.UserID,
		TokenHash:  hashRefreshToken(params.RawToken),
		ExpireTime: params.ExpireTime,
	})
}

func (r *RepositoryImpl) GetRefreshTokenByRaw(ctx context.Context, raw string) (*entity.RefreshToken, error) {
	row, err := r.q.GetRefreshTokenByHash(ctx, hashRefreshToken(raw))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToRefreshToken(&row), nil
}

func (r *RepositoryImpl) DeleteRefreshTokenByRaw(ctx context.Context, raw string, userID ulid.ULID) error {
	return r.q.DeleteRefreshTokenByHash(ctx, sqlc.DeleteRefreshTokenByHashParams{
		TokenHash: hashRefreshToken(raw),
		UserID:    userID,
	})
}

func (r *RepositoryImpl) DeleteExpiredRefreshTokensByUserID(ctx context.Context, userID ulid.ULID) error {
	return r.q.DeleteExpiredRefreshTokensByUserID(ctx, userID)
}
