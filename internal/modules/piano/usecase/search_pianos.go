package usecase

import (
	"context"
	"strings"

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

	attrs := gateway.AttributeFilters{
		MinAmbientNoiseAverage:   input.MinAmbientNoiseAverage,
		MinFootTrafficAverage:    input.MinFootTrafficAverage,
		MinResonanceAverage:      input.MinResonanceAverage,
		MinKeyTouchWeightAverage: input.MinKeyTouchWeightAverage,
		MinTuningQualityAverage:  input.MinTuningQualityAverage,
	}
	var pianos []*entity.Piano
	if q := strings.TrimSpace(input.Query); q != "" {
		ps, err := uc.pianoGateway.SearchByText(ctx, gateway.SearchByTextParams{
			Query:            "%" + escapeLike(q) + "%",
			Kind:             input.Kind,
			PianoType:        input.PianoType,
			PianoBrand:       input.PianoBrand,
			MinRatingAverage: input.MinRatingAverage,
			Attributes:       attrs,
			Limit:            limit,
		})
		if err != nil {
			return nil, err
		}
		pianos = ps
	} else if input.HasBounds {
		ps, err := uc.pianoGateway.SearchInBBox(ctx, gateway.SearchInBBoxParams{
			MinLat:           input.MinLat,
			MinLng:           input.MinLng,
			MaxLat:           input.MaxLat,
			MaxLng:           input.MaxLng,
			Kind:             input.Kind,
			PianoType:        input.PianoType,
			PianoBrand:       input.PianoBrand,
			MinRatingAverage: input.MinRatingAverage,
			Attributes:       attrs,
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
			PianoBrand:       input.PianoBrand,
			MinRatingAverage: input.MinRatingAverage,
			Attributes:       attrs,
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

// escapeLike は ILIKE のメタ文字 (% _ \) を無効化する。
func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}
