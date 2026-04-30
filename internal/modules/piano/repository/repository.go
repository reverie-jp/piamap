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

type CreatePianoParams struct {
	ID               ulid.ULID
	Name             string
	Description      *string
	Location         entity.LatLng
	Address          *string
	Prefecture       *string
	City             *string
	Kind             entity.PianoKind
	VenueType        *string
	PianoType        entity.PianoType
	PianoBrand       string
	PianoModel       *string
	ManufactureYear  *int16
	Hours            *string
	Availability     entity.PianoAvailability
	AvailabilityNote *string
	InstallTime      *time.Time
	CreatorUserID    *ulid.ULID
}

type UpdatePianoParams struct {
	ID                  ulid.ULID
	SetName             bool
	Name                *string
	SetDescription      bool
	Description         *string
	SetAddress          bool
	Address             *string
	SetPrefecture       bool
	Prefecture          *string
	SetCity             bool
	City                *string
	SetKind             bool
	Kind                *entity.PianoKind
	SetVenueType        bool
	VenueType           *string
	SetPianoType        bool
	PianoType           *entity.PianoType
	SetPianoBrand       bool
	PianoBrand          *string
	SetPianoModel       bool
	PianoModel          *string
	SetManufactureYear  bool
	ManufactureYear     *int16
	SetHours            bool
	Hours               *string
	SetStatus           bool
	Status              *entity.PianoStatus
	SetAvailability     bool
	Availability        *entity.PianoAvailability
	SetAvailabilityNote bool
	AvailabilityNote    *string
	SetInstallTime      bool
	InstallTime         *time.Time
	SetRemoveTime       bool
	RemoveTime          *time.Time
}

type AttributeFilters struct {
	MinAmbientNoiseAverage   *float64
	MinFootTrafficAverage    *float64
	MinResonanceAverage      *float64
	MinKeyTouchWeightAverage *float64
	MinTuningQualityAverage  *float64
}

type SearchInBBoxParams struct {
	MinLat           float64
	MinLng           float64
	MaxLat           float64
	MaxLng           float64
	Kind             *entity.PianoKind
	PianoType        *entity.PianoType
	PianoBrand       *string
	MinRatingAverage *float64
	Attributes       AttributeFilters
	Limit            int32
}

type SearchNearbyParams struct {
	CenterLat        float64
	CenterLng        float64
	RadiusM          float64
	Kind             *entity.PianoKind
	PianoType        *entity.PianoType
	PianoBrand       *string
	MinRatingAverage *float64
	Attributes       AttributeFilters
	Limit            int32
}

type SearchByTextParams struct {
	Query            string // ILIKE パターン (% 付きの形で渡す)
	Kind             *entity.PianoKind
	PianoType        *entity.PianoType
	PianoBrand       *string
	MinRatingAverage *float64
	Attributes       AttributeFilters
	Limit            int32
}

type CreatePianoEditParams struct {
	PianoID      ulid.ULID
	EditorUserID *ulid.ULID
	Operation    entity.PianoEditOperation
	Changes      []byte // JSONB の生バイト。nil 可
	Summary      *string
}

type ListPianoEditsParams struct {
	PianoID ulid.ULID
	AfterID *ulid.ULID
	Limit   int32
}

type Repository interface {
	GetPianoByID(ctx context.Context, id ulid.ULID) (*entity.Piano, error)
	SearchInBBox(ctx context.Context, params SearchInBBoxParams) ([]*entity.Piano, error)
	SearchNearby(ctx context.Context, params SearchNearbyParams) ([]*entity.Piano, error)
	SearchByText(ctx context.Context, params SearchByTextParams) ([]*entity.Piano, error)
	CreatePiano(ctx context.Context, params CreatePianoParams) error
	UpdatePiano(ctx context.Context, params UpdatePianoParams) error
	UpdatePianoLocation(ctx context.Context, id ulid.ULID, loc entity.LatLng) error
	CreatePianoEdit(ctx context.Context, params CreatePianoEditParams) error
	ListPianoEditsByPiano(ctx context.Context, params ListPianoEditsParams) ([]*entity.PianoEdit, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}

func (r *RepositoryImpl) GetPianoByID(ctx context.Context, id ulid.ULID) (*entity.Piano, error) {
	row, err := r.q.GetPianoByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToPianoFromGetRow(&row), nil
}

func toNullPianoKind(p *entity.PianoKind) sqlc.NullPianoKind {
	if p == nil {
		return sqlc.NullPianoKind{}
	}
	return sqlc.NullPianoKind{PianoKind: sqlc.PianoKind(*p), Valid: true}
}

func toNullPianoType(p *entity.PianoType) sqlc.NullPianoType {
	if p == nil {
		return sqlc.NullPianoType{}
	}
	return sqlc.NullPianoType{PianoType: sqlc.PianoType(*p), Valid: true}
}

func toNullPianoStatus(p *entity.PianoStatus) sqlc.NullPianoStatus {
	if p == nil {
		return sqlc.NullPianoStatus{}
	}
	return sqlc.NullPianoStatus{PianoStatus: sqlc.PianoStatus(*p), Valid: true}
}

func toNullPianoAvailability(p *entity.PianoAvailability) sqlc.NullPianoAvailability {
	if p == nil {
		return sqlc.NullPianoAvailability{}
	}
	return sqlc.NullPianoAvailability{PianoAvailability: sqlc.PianoAvailability(*p), Valid: true}
}

func (r *RepositoryImpl) SearchInBBox(ctx context.Context, params SearchInBBoxParams) ([]*entity.Piano, error) {
	rows, err := r.q.ListPianosInBBox(ctx, sqlc.ListPianosInBBoxParams{
		MinLng:           params.MinLng,
		MinLat:           params.MinLat,
		MaxLng:           params.MaxLng,
		MaxLat:           params.MaxLat,
		Kind:                     toNullPianoKind(params.Kind),
		PianoType:                toNullPianoType(params.PianoType),
		PianoBrand:               params.PianoBrand,
		MinRatingAverage:         params.MinRatingAverage,
		MinAmbientNoiseAverage:   params.Attributes.MinAmbientNoiseAverage,
		MinFootTrafficAverage:    params.Attributes.MinFootTrafficAverage,
		MinResonanceAverage:      params.Attributes.MinResonanceAverage,
		MinKeyTouchWeightAverage: params.Attributes.MinKeyTouchWeightAverage,
		MinTuningQualityAverage:  params.Attributes.MinTuningQualityAverage,
		LimitCount:               params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.Piano, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoFromBBoxRow(&rows[i])
	}
	return out, nil
}

func (r *RepositoryImpl) SearchNearby(ctx context.Context, params SearchNearbyParams) ([]*entity.Piano, error) {
	rows, err := r.q.ListPianosNearby(ctx, sqlc.ListPianosNearbyParams{
		CenterLng:        params.CenterLng,
		CenterLat:        params.CenterLat,
		RadiusM:          params.RadiusM,
		Kind:                     toNullPianoKind(params.Kind),
		PianoType:                toNullPianoType(params.PianoType),
		PianoBrand:               params.PianoBrand,
		MinRatingAverage:         params.MinRatingAverage,
		MinAmbientNoiseAverage:   params.Attributes.MinAmbientNoiseAverage,
		MinFootTrafficAverage:    params.Attributes.MinFootTrafficAverage,
		MinResonanceAverage:      params.Attributes.MinResonanceAverage,
		MinKeyTouchWeightAverage: params.Attributes.MinKeyTouchWeightAverage,
		MinTuningQualityAverage:  params.Attributes.MinTuningQualityAverage,
		LimitCount:               params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.Piano, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoFromNearbyRow(&rows[i])
	}
	return out, nil
}

func (r *RepositoryImpl) SearchByText(ctx context.Context, params SearchByTextParams) ([]*entity.Piano, error) {
	rows, err := r.q.SearchPianosByText(ctx, sqlc.SearchPianosByTextParams{
		QueryPattern:     params.Query,
		Kind:                     toNullPianoKind(params.Kind),
		PianoType:                toNullPianoType(params.PianoType),
		PianoBrand:               params.PianoBrand,
		MinRatingAverage:         params.MinRatingAverage,
		MinAmbientNoiseAverage:   params.Attributes.MinAmbientNoiseAverage,
		MinFootTrafficAverage:    params.Attributes.MinFootTrafficAverage,
		MinResonanceAverage:      params.Attributes.MinResonanceAverage,
		MinKeyTouchWeightAverage: params.Attributes.MinKeyTouchWeightAverage,
		MinTuningQualityAverage:  params.Attributes.MinTuningQualityAverage,
		LimitCount:               params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.Piano, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoFromTextSearchRow(&rows[i])
	}
	return out, nil
}

func (r *RepositoryImpl) CreatePiano(ctx context.Context, params CreatePianoParams) error {
	return r.q.CreatePiano(ctx, sqlc.CreatePianoParams{
		ID:               params.ID,
		Name:             params.Name,
		Description:      params.Description,
		Longitude:        params.Location.Longitude,
		Latitude:         params.Location.Latitude,
		Address:          params.Address,
		Prefecture:       params.Prefecture,
		City:             params.City,
		Kind:             sqlc.PianoKind(params.Kind),
		VenueType:        params.VenueType,
		PianoType:        sqlc.PianoType(params.PianoType),
		PianoBrand:       params.PianoBrand,
		PianoModel:       params.PianoModel,
		ManufactureYear:  params.ManufactureYear,
		Hours:            params.Hours,
		Availability:     sqlc.PianoAvailability(params.Availability),
		AvailabilityNote: params.AvailabilityNote,
		InstallTime:      params.InstallTime,
		CreatorUserID:    params.CreatorUserID,
	})
}

func (r *RepositoryImpl) UpdatePiano(ctx context.Context, params UpdatePianoParams) error {
	return r.q.UpdatePiano(ctx, sqlc.UpdatePianoParams{
		ID:                  params.ID,
		SetName:             params.SetName,
		Name:                params.Name,
		SetDescription:      params.SetDescription,
		Description:         params.Description,
		SetAddress:          params.SetAddress,
		Address:             params.Address,
		SetPrefecture:       params.SetPrefecture,
		Prefecture:          params.Prefecture,
		SetCity:             params.SetCity,
		City:                params.City,
		SetKind:             params.SetKind,
		Kind:                toNullPianoKind(params.Kind),
		SetVenueType:        params.SetVenueType,
		VenueType:           params.VenueType,
		SetPianoType:        params.SetPianoType,
		PianoType:           toNullPianoType(params.PianoType),
		SetPianoBrand:       params.SetPianoBrand,
		PianoBrand:          params.PianoBrand,
		SetPianoModel:       params.SetPianoModel,
		PianoModel:          params.PianoModel,
		SetManufactureYear:  params.SetManufactureYear,
		ManufactureYear:     params.ManufactureYear,
		SetHours:            params.SetHours,
		Hours:               params.Hours,
		SetStatus:           params.SetStatus,
		Status:              toNullPianoStatus(params.Status),
		SetAvailability:     params.SetAvailability,
		Availability:        toNullPianoAvailability(params.Availability),
		SetAvailabilityNote: params.SetAvailabilityNote,
		AvailabilityNote:    params.AvailabilityNote,
		SetInstallTime:      params.SetInstallTime,
		InstallTime:         params.InstallTime,
		SetRemoveTime:       params.SetRemoveTime,
		RemoveTime:          params.RemoveTime,
	})
}

func (r *RepositoryImpl) UpdatePianoLocation(ctx context.Context, id ulid.ULID, loc entity.LatLng) error {
	return r.q.UpdatePianoLocation(ctx, sqlc.UpdatePianoLocationParams{
		ID:        id,
		Longitude: loc.Longitude,
		Latitude:  loc.Latitude,
	})
}

func (r *RepositoryImpl) CreatePianoEdit(ctx context.Context, params CreatePianoEditParams) error {
	return r.q.CreatePianoEdit(ctx, sqlc.CreatePianoEditParams{
		ID:           ulid.New(),
		PianoID:      params.PianoID,
		EditorUserID: params.EditorUserID,
		Operation:    sqlc.PianoEditOperation(params.Operation),
		Changes:      params.Changes,
		Summary:      params.Summary,
	})
}

func (r *RepositoryImpl) ListPianoEditsByPiano(ctx context.Context, params ListPianoEditsParams) ([]*entity.PianoEdit, error) {
	rows, err := r.q.ListPianoEditsByPiano(ctx, sqlc.ListPianoEditsByPianoParams{
		PianoID:    params.PianoID,
		AfterID:    params.AfterID,
		LimitCount: params.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.PianoEdit, len(rows))
	for i := range rows {
		out[i] = mapper.ToPianoEdit(&rows[i])
	}
	return out, nil
}
