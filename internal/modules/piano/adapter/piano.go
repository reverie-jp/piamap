package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
)

func ToPiano(view *gateway.PianoView) *pianov1.Piano {
	if view == nil || view.Piano == nil {
		return nil
	}
	p := view.Piano
	pb := &pianov1.Piano{
		Name:        resourcename.FormatPiano(p.ID),
		DisplayName: p.Name,
		Description: p.Description,
		Location: &pianov1.LatLng{
			Latitude:  p.Location.Latitude,
			Longitude: p.Location.Longitude,
		},
		Address:                p.Address,
		Prefecture:             p.Prefecture,
		City:                   p.City,
		Kind:                   toPbKind(p.Kind),
		VenueType:              p.VenueType,
		PianoType:              toPbPianoType(p.PianoType),
		PianoBrand:             p.PianoBrand,
		PianoModel:             p.PianoModel,
		Hours:                  p.Hours,
		Status:                 toPbStatus(p.Status),
		Availability:           toPbAvailability(p.Availability),
		AvailabilityNote:       p.AvailabilityNote,
		PostCount:              p.PostCount,
		RatingCount:            p.RatingCount,
		RatingAverage:          p.RatingAverage(),
		AmbientNoiseAverage:    p.AmbientNoiseAverage(),
		FootTrafficAverage:     p.FootTrafficAverage(),
		ResonanceAverage:       p.ResonanceAverage(),
		KeyTouchWeightAverage:  p.KeyTouchWeightAverage(),
		TuningQualityAverage:   p.TuningQualityAverage(),
		WishlistCount:          p.WishlistCount,
		VisitedCount:           p.VisitedCount,
		FavoriteCount:          p.FavoriteCount,
		CreateTime:             timestamppb.New(p.CreateTime),
		UpdateTime:             timestamppb.New(p.UpdateTime),
	}
	if p.ManufactureYear != nil {
		v := int32(*p.ManufactureYear)
		pb.ManufactureYear = &v
	}
	if p.InstallTime != nil {
		pb.InstallTime = timestamppb.New(*p.InstallTime)
	}
	if p.RemoveTime != nil {
		pb.RemoveTime = timestamppb.New(*p.RemoveTime)
	}
	if view.CreatorCustomID != "" {
		pb.Creator = resourcename.FormatUser(view.CreatorCustomID)
	}
	if p.DistanceM > 0 {
		d := p.DistanceM
		pb.DistanceM = &d
	}
	return pb
}

func toPbKind(k entity.PianoKind) pianov1.PianoKind {
	switch k {
	case entity.PianoKindStreet:
		return pianov1.PianoKind_PIANO_KIND_STREET
	case entity.PianoKindPracticeRoom:
		return pianov1.PianoKind_PIANO_KIND_PRACTICE_ROOM
	case entity.PianoKindOther:
		return pianov1.PianoKind_PIANO_KIND_OTHER
	}
	return pianov1.PianoKind_PIANO_KIND_UNSPECIFIED
}

func fromPbKind(k pianov1.PianoKind) (entity.PianoKind, bool) {
	switch k {
	case pianov1.PianoKind_PIANO_KIND_STREET:
		return entity.PianoKindStreet, true
	case pianov1.PianoKind_PIANO_KIND_PRACTICE_ROOM:
		return entity.PianoKindPracticeRoom, true
	case pianov1.PianoKind_PIANO_KIND_OTHER:
		return entity.PianoKindOther, true
	}
	return "", false
}

func toPbPianoType(t entity.PianoType) pianov1.PianoType {
	switch t {
	case entity.PianoTypeGrand:
		return pianov1.PianoType_PIANO_TYPE_GRAND
	case entity.PianoTypeUpright:
		return pianov1.PianoType_PIANO_TYPE_UPRIGHT
	case entity.PianoTypeElectronic:
		return pianov1.PianoType_PIANO_TYPE_ELECTRONIC
	case entity.PianoTypeUnknown:
		return pianov1.PianoType_PIANO_TYPE_UNKNOWN
	}
	return pianov1.PianoType_PIANO_TYPE_UNSPECIFIED
}

func fromPbPianoType(t pianov1.PianoType) (entity.PianoType, bool) {
	switch t {
	case pianov1.PianoType_PIANO_TYPE_GRAND:
		return entity.PianoTypeGrand, true
	case pianov1.PianoType_PIANO_TYPE_UPRIGHT:
		return entity.PianoTypeUpright, true
	case pianov1.PianoType_PIANO_TYPE_ELECTRONIC:
		return entity.PianoTypeElectronic, true
	case pianov1.PianoType_PIANO_TYPE_UNKNOWN:
		return entity.PianoTypeUnknown, true
	}
	return "", false
}

func toPbStatus(s entity.PianoStatus) pianov1.PianoStatus {
	switch s {
	case entity.PianoStatusPending:
		return pianov1.PianoStatus_PIANO_STATUS_PENDING
	case entity.PianoStatusActive:
		return pianov1.PianoStatus_PIANO_STATUS_ACTIVE
	case entity.PianoStatusTemporary:
		return pianov1.PianoStatus_PIANO_STATUS_TEMPORARY
	case entity.PianoStatusRemoved:
		return pianov1.PianoStatus_PIANO_STATUS_REMOVED
	}
	return pianov1.PianoStatus_PIANO_STATUS_UNSPECIFIED
}

func fromPbStatus(s pianov1.PianoStatus) (entity.PianoStatus, bool) {
	switch s {
	case pianov1.PianoStatus_PIANO_STATUS_PENDING:
		return entity.PianoStatusPending, true
	case pianov1.PianoStatus_PIANO_STATUS_ACTIVE:
		return entity.PianoStatusActive, true
	case pianov1.PianoStatus_PIANO_STATUS_TEMPORARY:
		return entity.PianoStatusTemporary, true
	case pianov1.PianoStatus_PIANO_STATUS_REMOVED:
		return entity.PianoStatusRemoved, true
	}
	return "", false
}

func toPbAvailability(a entity.PianoAvailability) pianov1.PianoAvailability {
	switch a {
	case entity.PianoAvailabilityRegular:
		return pianov1.PianoAvailability_PIANO_AVAILABILITY_REGULAR
	case entity.PianoAvailabilityIrregular:
		return pianov1.PianoAvailability_PIANO_AVAILABILITY_IRREGULAR
	case entity.PianoAvailabilityEventOnly:
		return pianov1.PianoAvailability_PIANO_AVAILABILITY_EVENT_ONLY
	case entity.PianoAvailabilityWeatherDependent:
		return pianov1.PianoAvailability_PIANO_AVAILABILITY_WEATHER_DEPENDENT
	}
	return pianov1.PianoAvailability_PIANO_AVAILABILITY_UNSPECIFIED
}

func fromPbAvailability(a pianov1.PianoAvailability) (entity.PianoAvailability, bool) {
	switch a {
	case pianov1.PianoAvailability_PIANO_AVAILABILITY_REGULAR:
		return entity.PianoAvailabilityRegular, true
	case pianov1.PianoAvailability_PIANO_AVAILABILITY_IRREGULAR:
		return entity.PianoAvailabilityIrregular, true
	case pianov1.PianoAvailability_PIANO_AVAILABILITY_EVENT_ONLY:
		return entity.PianoAvailabilityEventOnly, true
	case pianov1.PianoAvailability_PIANO_AVAILABILITY_WEATHER_DEPENDENT:
		return entity.PianoAvailabilityWeatherDependent, true
	}
	return "", false
}
