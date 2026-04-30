package gateway

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano/repository"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// PianoView は Piano entity に requester 視点 + 表示用 creator info を足したもの。
type PianoView struct {
	Piano           *entity.Piano
	CreatorCustomID string
	IsCreator       bool
}

// PianoEditView は PianoEdit に表示用の editor 情報を付与したもの。
type PianoEditView struct {
	Edit              *entity.PianoEdit
	EditorCustomID    string
	EditorDisplayName string
}

type CreatePianoParams = repository.CreatePianoParams
type UpdatePianoParams = repository.UpdatePianoParams
type SearchInBBoxParams = repository.SearchInBBoxParams
type SearchNearbyParams = repository.SearchNearbyParams
type SearchByTextParams = repository.SearchByTextParams
type AttributeFilters = repository.AttributeFilters
type CreatePianoEditParams = repository.CreatePianoEditParams
type ListPianoEditsParams = repository.ListPianoEditsParams

type Gateway interface {
	GetPiano(ctx context.Context, id ulid.ULID) (*entity.Piano, error)
	SearchInBBox(ctx context.Context, params SearchInBBoxParams) ([]*entity.Piano, error)
	SearchNearby(ctx context.Context, params SearchNearbyParams) ([]*entity.Piano, error)
	SearchByText(ctx context.Context, params SearchByTextParams) ([]*entity.Piano, error)
	CreatePiano(ctx context.Context, params CreatePianoParams) error
	UpdatePiano(ctx context.Context, params UpdatePianoParams) error
	UpdatePianoLocation(ctx context.Context, id ulid.ULID, loc entity.LatLng) error
	CreatePianoEdit(ctx context.Context, params CreatePianoEditParams) error
	ListPianoEditsByPiano(ctx context.Context, params ListPianoEditsParams) ([]*entity.PianoEdit, error)
	BuildPianoView(ctx context.Context, requesterID ulid.ULID, piano *entity.Piano) (*PianoView, error)
	BuildListPianoViews(ctx context.Context, requesterID ulid.ULID, pianos []*entity.Piano) ([]*PianoView, error)
	BuildListPianoEditViews(ctx context.Context, edits []*entity.PianoEdit) ([]*PianoEditView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
}

func New(q sqlc.Querier, userGateway usergw.Gateway) Gateway {
	return &gatewayImpl{
		repo:        repository.New(q),
		userGateway: userGateway,
	}
}

func (g *gatewayImpl) GetPiano(ctx context.Context, id ulid.ULID) (*entity.Piano, error) {
	return g.repo.GetPianoByID(ctx, id)
}

func (g *gatewayImpl) SearchInBBox(ctx context.Context, params SearchInBBoxParams) ([]*entity.Piano, error) {
	return g.repo.SearchInBBox(ctx, params)
}

func (g *gatewayImpl) SearchNearby(ctx context.Context, params SearchNearbyParams) ([]*entity.Piano, error) {
	return g.repo.SearchNearby(ctx, params)
}

func (g *gatewayImpl) SearchByText(ctx context.Context, params SearchByTextParams) ([]*entity.Piano, error) {
	return g.repo.SearchByText(ctx, params)
}

func (g *gatewayImpl) CreatePiano(ctx context.Context, params CreatePianoParams) error {
	return g.repo.CreatePiano(ctx, params)
}

func (g *gatewayImpl) UpdatePiano(ctx context.Context, params UpdatePianoParams) error {
	return g.repo.UpdatePiano(ctx, params)
}

func (g *gatewayImpl) UpdatePianoLocation(ctx context.Context, id ulid.ULID, loc entity.LatLng) error {
	return g.repo.UpdatePianoLocation(ctx, id, loc)
}

func (g *gatewayImpl) CreatePianoEdit(ctx context.Context, params CreatePianoEditParams) error {
	return g.repo.CreatePianoEdit(ctx, params)
}

func (g *gatewayImpl) ListPianoEditsByPiano(ctx context.Context, params ListPianoEditsParams) ([]*entity.PianoEdit, error) {
	return g.repo.ListPianoEditsByPiano(ctx, params)
}

func (g *gatewayImpl) BuildListPianoEditViews(ctx context.Context, edits []*entity.PianoEdit) ([]*PianoEditView, error) {
	if len(edits) == 0 {
		return []*PianoEditView{}, nil
	}
	editorIDs := make([]ulid.ULID, 0, len(edits))
	seen := make(map[string]bool, len(edits))
	for _, e := range edits {
		if e == nil || e.EditorUserID == nil {
			continue
		}
		k := e.EditorUserID.String()
		if seen[k] {
			continue
		}
		seen[k] = true
		editorIDs = append(editorIDs, *e.EditorUserID)
	}
	editorByID := make(map[string]*entity.User, len(editorIDs))
	if len(editorIDs) > 0 {
		users, err := g.userGateway.ListUsersByIDs(ctx, editorIDs)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			editorByID[u.ID.String()] = u
		}
	}
	views := make([]*PianoEditView, len(edits))
	for i, e := range edits {
		if e == nil {
			continue
		}
		v := &PianoEditView{Edit: e}
		if e.EditorUserID != nil {
			if u, ok := editorByID[e.EditorUserID.String()]; ok {
				v.EditorCustomID = u.CustomID
				v.EditorDisplayName = u.DisplayName
			}
		}
		views[i] = v
	}
	return views, nil
}

func (g *gatewayImpl) BuildPianoView(ctx context.Context, requesterID ulid.ULID, piano *entity.Piano) (*PianoView, error) {
	views, err := g.BuildListPianoViews(ctx, requesterID, []*entity.Piano{piano})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

// BuildListPianoViews は creator user を 1 回の ListUsersByIDs で取得して N+1 を避ける。
func (g *gatewayImpl) BuildListPianoViews(ctx context.Context, requesterID ulid.ULID, pianos []*entity.Piano) ([]*PianoView, error) {
	if len(pianos) == 0 {
		return []*PianoView{}, nil
	}

	creatorIDs := make([]ulid.ULID, 0, len(pianos))
	seen := make(map[string]bool, len(pianos))
	for _, p := range pianos {
		if p == nil || p.CreatorUserID == nil {
			continue
		}
		key := p.CreatorUserID.String()
		if seen[key] {
			continue
		}
		seen[key] = true
		creatorIDs = append(creatorIDs, *p.CreatorUserID)
	}

	creatorByID := make(map[string]string, len(creatorIDs))
	if len(creatorIDs) > 0 {
		users, err := g.userGateway.ListUsersByIDs(ctx, creatorIDs)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			creatorByID[u.ID.String()] = u.CustomID
		}
	}

	views := make([]*PianoView, len(pianos))
	for i, p := range pianos {
		if p == nil {
			continue
		}
		view := &PianoView{Piano: p}
		if p.CreatorUserID != nil {
			view.CreatorCustomID = creatorByID[p.CreatorUserID.String()]
			view.IsCreator = !requesterID.IsZero() && *p.CreatorUserID == requesterID
		}
		views[i] = view
	}
	return views, nil
}
