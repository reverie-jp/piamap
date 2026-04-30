package usecase

import (
	"errors"
	"time"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// 信頼ライン保護: 大幅な座標移動 (m) と「removed への変更」は trusted user のみ。
const TrustedRequiredMoveDistanceM = 500

type UpdatePianoInput struct {
	RequesterID ulid.ULID
	PianoID     ulid.ULID

	SetName        bool
	Name           string
	SetDescription bool
	Description    *string

	SetLocation bool
	Location    entity.LatLng

	SetAddress    bool
	Address       *string
	SetPrefecture bool
	Prefecture    *string
	SetCity       bool
	City          *string

	SetKind bool
	Kind    entity.PianoKind

	SetVenueType bool
	VenueType    *string

	SetPianoType  bool
	PianoType     entity.PianoType
	SetPianoBrand bool
	PianoBrand    string
	SetPianoModel bool
	PianoModel    *string

	SetManufactureYear bool
	ManufactureYear    *int16

	SetHours bool
	Hours    *string

	SetStatus bool
	Status    entity.PianoStatus

	SetAvailability     bool
	Availability        entity.PianoAvailability
	SetAvailabilityNote bool
	AvailabilityNote    *string

	SetInstallTime bool
	InstallTime    *time.Time
	SetRemoveTime  bool
	RemoveTime     *time.Time

	EditSummary *string
}

func (i UpdatePianoInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester is required")
	}
	if i.PianoID.IsZero() {
		return errors.New("piano id is required")
	}
	if i.SetName {
		if l := len(i.Name); l == 0 || l > 80 {
			return errors.New("name length must be 1..80")
		}
	}
	if i.SetLocation {
		if !validLat(i.Location.Latitude) || !validLng(i.Location.Longitude) {
			return errors.New("location out of range")
		}
	}
	if i.SetKind && !isValidKind(i.Kind) {
		return errors.New("invalid kind")
	}
	if i.SetPianoType && !isValidPianoType(i.PianoType) {
		return errors.New("invalid piano_type")
	}
	if i.SetAvailability && !isValidAvailability(i.Availability) {
		return errors.New("invalid availability")
	}
	if i.SetStatus && !isValidStatus(i.Status) {
		return errors.New("invalid status")
	}
	if i.SetPianoBrand && i.PianoBrand == "" {
		return errors.New("piano_brand cannot be empty")
	}
	return nil
}

// HasAnyChange は反映対象フィールドが何かしらあるかを返す。
func (i UpdatePianoInput) HasAnyChange() bool {
	return i.SetName || i.SetDescription || i.SetLocation ||
		i.SetAddress || i.SetPrefecture || i.SetCity ||
		i.SetKind || i.SetVenueType ||
		i.SetPianoType || i.SetPianoBrand || i.SetPianoModel ||
		i.SetManufactureYear || i.SetHours ||
		i.SetStatus || i.SetAvailability || i.SetAvailabilityNote ||
		i.SetInstallTime || i.SetRemoveTime
}
