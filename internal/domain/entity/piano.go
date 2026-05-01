package entity

import (
	"time"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type PianoKind string
type PianoStatus string
type PianoType string
type PianoAvailability string
type PianoEditOperation string

const (
	PianoKindStreet       PianoKind = "street"
	PianoKindPracticeRoom PianoKind = "practice_room"
	PianoKindOther        PianoKind = "other"

	PianoStatusPending   PianoStatus = "pending"
	PianoStatusActive    PianoStatus = "active"
	PianoStatusTemporary PianoStatus = "temporary"
	PianoStatusRemoved   PianoStatus = "removed"

	PianoTypeGrand      PianoType = "grand"
	PianoTypeUpright    PianoType = "upright"
	PianoTypeElectronic PianoType = "electronic"
	PianoTypeUnknown    PianoType = "unknown"

	PianoAvailabilityRegular          PianoAvailability = "regular"
	PianoAvailabilityIrregular        PianoAvailability = "irregular"
	PianoAvailabilityEventOnly        PianoAvailability = "event_only"
	PianoAvailabilityWeatherDependent PianoAvailability = "weather_dependent"

	PianoEditOpCreate       PianoEditOperation = "create"
	PianoEditOpUpdate       PianoEditOperation = "update"
	PianoEditOpPhotoAdd     PianoEditOperation = "photo_add"
	PianoEditOpPhotoRemove  PianoEditOperation = "photo_remove"
	PianoEditOpStatusChange PianoEditOperation = "status_change"
	PianoEditOpKindChange   PianoEditOperation = "kind_change"
	PianoEditOpRestore      PianoEditOperation = "restore"
)

// LatLng は WGS84 (EPSG:4326) 表現の緯度経度。
type LatLng struct {
	Latitude  float64
	Longitude float64
}

// Piano はピアノの主要属性 + 集計列。
// `Aggregates` を embed して View 層でも同じ計算 (avg = sum/count) を使えるようにする。
type Piano struct {
	ID               ulid.ULID
	Name             string
	Description      *string
	Location         LatLng
	Address          *string
	Prefecture       *string
	City             *string
	Kind             PianoKind
	VenueType        *string
	PianoType        PianoType
	PianoBrand       string
	PianoModel       *string
	ManufactureYear  *int16
	Hours            *string
	Status           PianoStatus
	Availability     PianoAvailability
	AvailabilityNote *string
	InstallTime      *time.Time
	RemoveTime       *time.Time
	CreatorUserID    *ulid.ULID

	PostCount           int32
	RatingCount         int32
	RatingSum           int32
	AmbientNoiseCount   int32
	AmbientNoiseSum     int32
	FootTrafficCount    int32
	FootTrafficSum      int32
	ResonanceCount      int32
	ResonanceSum        int32
	KeyTouchWeightCount int32
	KeyTouchWeightSum   int32
	TuningQualityCount  int32
	TuningQualitySum    int32

	WishlistCount int32
	VisitedCount  int32
	FavoriteCount int32

	CreateTime time.Time
	UpdateTime time.Time

	// SearchPianos の結果でセットされる (それ以外は 0)。
	DistanceM float64
}

// RatingAverage: rating_count = 0 の場合は 0 を返す (UI 側で「未評価」表示)。
func (p *Piano) RatingAverage() float64 {
	if p == nil || p.RatingCount == 0 {
		return 0
	}
	return float64(p.RatingSum) / float64(p.RatingCount)
}

func avg(sum, count int32) float64 {
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
}

func (p *Piano) AmbientNoiseAverage() float64 {
	if p == nil {
		return 0
	}
	return avg(p.AmbientNoiseSum, p.AmbientNoiseCount)
}

func (p *Piano) FootTrafficAverage() float64 {
	if p == nil {
		return 0
	}
	return avg(p.FootTrafficSum, p.FootTrafficCount)
}

func (p *Piano) ResonanceAverage() float64 {
	if p == nil {
		return 0
	}
	return avg(p.ResonanceSum, p.ResonanceCount)
}

func (p *Piano) KeyTouchWeightAverage() float64 {
	if p == nil {
		return 0
	}
	return avg(p.KeyTouchWeightSum, p.KeyTouchWeightCount)
}

func (p *Piano) TuningQualityAverage() float64 {
	if p == nil {
		return 0
	}
	return avg(p.TuningQualitySum, p.TuningQualityCount)
}
