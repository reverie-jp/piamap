package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

const (
	SearchPianosDefaultLimit = 200
	SearchPianosMaxLimit     = 500
	SearchPianosMaxRadiusM   = 50_000 // 50 km
)

type SearchPianosInput struct {
	RequesterID ulid.ULID // ゲスト可

	// bbox 検索 (UI のマップビューポート優先)。指定があれば center/radius は無視。
	HasBounds bool
	MinLat    float64
	MinLng    float64
	MaxLat    float64
	MaxLng    float64

	// center+radius 検索。
	HasCenter bool
	CenterLat float64
	CenterLng float64
	RadiusM   float64

	Limit            int32
	Kind             *entity.PianoKind
	PianoType        *entity.PianoType
	MinRatingAverage *float64
}

func (i SearchPianosInput) Validate() error {
	if !i.HasBounds && !i.HasCenter {
		return errors.New("bounds or center+radius is required")
	}
	if i.HasBounds {
		if i.MinLat >= i.MaxLat || i.MinLng >= i.MaxLng {
			return errors.New("invalid bounds (min must be < max)")
		}
		if !validLat(i.MinLat) || !validLat(i.MaxLat) || !validLng(i.MinLng) || !validLng(i.MaxLng) {
			return errors.New("bounds out of range")
		}
	}
	if i.HasCenter {
		if !validLat(i.CenterLat) || !validLng(i.CenterLng) {
			return errors.New("center out of range")
		}
		if i.RadiusM <= 0 || i.RadiusM > SearchPianosMaxRadiusM {
			return errors.New("radius_m must be in (0, 50000]")
		}
	}
	if i.Limit < 0 || i.Limit > SearchPianosMaxLimit {
		return errors.New("limit out of range")
	}
	if i.MinRatingAverage != nil {
		if *i.MinRatingAverage < 0 || *i.MinRatingAverage > 5 {
			return errors.New("min_rating_average must be in [0, 5]")
		}
	}
	return nil
}

func validLat(lat float64) bool { return lat >= -90 && lat <= 90 }
func validLng(lng float64) bool { return lng >= -180 && lng <= 180 }
