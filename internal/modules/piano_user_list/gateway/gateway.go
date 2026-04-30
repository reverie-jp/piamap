package gateway

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/repository"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type ListByUserParams = repository.ListByUserParams

type Gateway interface {
	UpsertList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error
	DeleteList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error
	ListByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error)
	ListKindsForPiano(ctx context.Context, userID, pianoID ulid.ULID) ([]entity.PianoListKind, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}

func (g *gatewayImpl) UpsertList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error {
	return g.repo.UpsertList(ctx, userID, pianoID, kind)
}
func (g *gatewayImpl) DeleteList(ctx context.Context, userID, pianoID ulid.ULID, kind entity.PianoListKind) error {
	return g.repo.DeleteList(ctx, userID, pianoID, kind)
}
func (g *gatewayImpl) ListByUser(ctx context.Context, params ListByUserParams) ([]ulid.ULID, error) {
	return g.repo.ListByUser(ctx, params)
}
func (g *gatewayImpl) ListKindsForPiano(ctx context.Context, userID, pianoID ulid.ULID) ([]entity.PianoListKind, error) {
	return g.repo.ListKindsForPiano(ctx, userID, pianoID)
}
