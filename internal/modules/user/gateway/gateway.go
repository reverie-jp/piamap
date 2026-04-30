package gateway

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/user/repository"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// UserView は User entity に requester 視点のフラグを足した表示用型。
type UserView struct {
	User *entity.User
	IsMe bool
}

type CreateUserParams = repository.CreateUserParams
type UpdateUserProfileParams = repository.UpdateUserProfileParams

type Gateway interface {
	GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error)
	GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error)
	ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	DeleteUser(ctx context.Context, id ulid.ULID) error
	UpdateUserProfile(ctx context.Context, params UpdateUserProfileParams) error
	UpdateUserCustomID(ctx context.Context, id ulid.ULID, customID string) error
	IsUserCurrentlyRestricted(ctx context.Context, id ulid.ULID) (bool, error)
	BuildUserView(ctx context.Context, requesterID, id ulid.ULID) (*UserView, error)
	BuildListUserViews(ctx context.Context, requesterID ulid.ULID, ids []ulid.ULID) ([]*UserView, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}

func (g *gatewayImpl) GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error) {
	return g.repo.GetUserByID(ctx, id)
}

func (g *gatewayImpl) GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error) {
	return g.repo.GetUserByCustomID(ctx, customID)
}

func (g *gatewayImpl) ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error) {
	return g.repo.ListUsersByIDs(ctx, ids)
}

func (g *gatewayImpl) CreateUser(ctx context.Context, params CreateUserParams) error {
	return g.repo.CreateUser(ctx, params)
}

func (g *gatewayImpl) DeleteUser(ctx context.Context, id ulid.ULID) error {
	return g.repo.DeleteUser(ctx, id)
}

func (g *gatewayImpl) UpdateUserProfile(ctx context.Context, params UpdateUserProfileParams) error {
	return g.repo.UpdateUserProfile(ctx, params)
}

func (g *gatewayImpl) UpdateUserCustomID(ctx context.Context, id ulid.ULID, customID string) error {
	return g.repo.UpdateUserCustomID(ctx, id, customID)
}

func (g *gatewayImpl) IsUserCurrentlyRestricted(ctx context.Context, id ulid.ULID) (bool, error) {
	return g.repo.IsUserCurrentlyRestricted(ctx, id)
}

func (g *gatewayImpl) BuildUserView(ctx context.Context, requesterID, id ulid.ULID) (*UserView, error) {
	views, err := g.BuildListUserViews(ctx, requesterID, []ulid.ULID{id})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

func (g *gatewayImpl) BuildListUserViews(ctx context.Context, requesterID ulid.ULID, ids []ulid.ULID) ([]*UserView, error) {
	if len(ids) == 0 {
		return []*UserView{}, nil
	}
	users, err := g.repo.ListUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	views := make([]*UserView, len(users))
	for i, u := range users {
		views[i] = &UserView{
			User: u,
			IsMe: !requesterID.IsZero() && u.ID == requesterID,
		}
	}
	return views, nil
}
