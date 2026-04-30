package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
)

type SearchPianos struct {
	pianoGateway gateway.Gateway
}

func NewSearchPianos(pianoGateway gateway.Gateway) *SearchPianos {
	return &SearchPianos{pianoGateway: pianoGateway}
}

func (uc *SearchPianos) Execute(ctx context.Context, input SearchPianosInput) (*SearchPianosOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	limit := input.Limit
	if limit == 0 {
		limit = SearchPianosDefaultLimit
	}

	var pianos []*entity.Piano
	if input.HasBounds {
		ps, err := uc.pianoGateway.SearchInBBox(ctx, gateway.SearchInBBoxParams{
			MinLat:           input.MinLat,
			MinLng:           input.MinLng,
			MaxLat:           input.MaxLat,
			MaxLng:           input.MaxLng,
			Kind:             input.Kind,
			PianoType:        input.PianoType,
			MinRatingAverage: input.MinRatingAverage,
			Limit:            limit,
		})
		if err != nil {
			return nil, err
		}
		pianos = ps
	} else {
		ps, err := uc.pianoGateway.SearchNearby(ctx, gateway.SearchNearbyParams{
			CenterLat:        input.CenterLat,
			CenterLng:        input.CenterLng,
			RadiusM:          input.RadiusM,
			Kind:             input.Kind,
			PianoType:        input.PianoType,
			MinRatingAverage: input.MinRatingAverage,
			Limit:            limit,
		})
		if err != nil {
			return nil, err
		}
		pianos = ps
	}

	views, err := uc.pianoGateway.BuildListPianoViews(ctx, input.RequesterID, pianos)
	if err != nil {
		return nil, err
	}
	return &SearchPianosOutput{Views: views}, nil
}
