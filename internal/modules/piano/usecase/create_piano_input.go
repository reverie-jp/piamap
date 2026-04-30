package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type CreatePianoInput struct {
	RequesterID ulid.ULID

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
}

func (i CreatePianoInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester is required")
	}
	if l := len(i.Name); l == 0 || l > 80 {
		return errors.New("name length must be 1..80")
	}
	if !validLat(i.Location.Latitude) || !validLng(i.Location.Longitude) {
		return errors.New("location out of range")
	}
	if !isValidKind(i.Kind) {
		return errors.New("invalid kind")
	}
	if !isValidPianoType(i.PianoType) {
		return errors.New("invalid piano_type")
	}
	if !isValidAvailability(i.Availability) {
		return errors.New("invalid availability")
	}
	if i.PianoBrand == "" {
		return errors.New("piano_brand is required (use 'unknown' if unsure)")
	}
	if i.ManufactureYear != nil {
		if *i.ManufactureYear < 1700 || *i.ManufactureYear > 2100 {
			return errors.New("manufacture_year out of range")
		}
	}
	return nil
}

func isValidKind(k entity.PianoKind) bool {
	switch k {
	case entity.PianoKindStreet, entity.PianoKindPracticeRoom, entity.PianoKindOther:
		return true
	}
	return false
}

func isValidPianoType(t entity.PianoType) bool {
	switch t {
	case entity.PianoTypeGrand, entity.PianoTypeUpright, entity.PianoTypeElectronic, entity.PianoTypeUnknown:
		return true
	}
	return false
}

func isValidAvailability(a entity.PianoAvailability) bool {
	switch a {
	case entity.PianoAvailabilityRegular,
		entity.PianoAvailabilityIrregular,
		entity.PianoAvailabilityEventOnly,
		entity.PianoAvailabilityWeatherDependent:
		return true
	}
	return false
}

func isValidStatus(s entity.PianoStatus) bool {
	switch s {
	case entity.PianoStatusPending, entity.PianoStatusActive, entity.PianoStatusTemporary, entity.PianoStatusRemoved:
		return true
	}
	return false
}
