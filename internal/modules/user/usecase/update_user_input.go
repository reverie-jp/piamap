package usecase

import (
	"errors"
	"regexp"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

var customIDPattern = regexp.MustCompile(`^[a-z0-9_]{3,20}$`)

// UpdateUserInput は FieldMask で部分更新する。Set* フラグが true のフィールドだけ反映する。
type UpdateUserInput struct {
	RequesterID ulid.ULID

	SetCustomID bool
	CustomID    string

	SetDisplayName bool
	DisplayName    string

	SetBiography bool
	Biography    *string

	SetAvatarURL bool
	AvatarURL    *string

	SetHometown bool
	Hometown    *string

	SetPianoStartedYear bool
	PianoStartedYear    *int16

	SetYearsOfExperience bool
	YearsOfExperience    *int16
}

func (i UpdateUserInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester is required")
	}
	if i.SetCustomID {
		if !customIDPattern.MatchString(i.CustomID) {
			return errors.New("custom_id must match ^[a-z0-9_]{3,20}$")
		}
	}
	if i.SetDisplayName {
		if l := len(i.DisplayName); l == 0 || l > 30 {
			return errors.New("display_name length must be 1..30")
		}
	}
	if i.SetPianoStartedYear && i.PianoStartedYear != nil {
		if *i.PianoStartedYear < 1900 || *i.PianoStartedYear > 2100 {
			return errors.New("piano_started_year out of range")
		}
	}
	if i.SetYearsOfExperience && i.YearsOfExperience != nil {
		if *i.YearsOfExperience < 0 || *i.YearsOfExperience > 100 {
			return errors.New("years_of_experience out of range")
		}
	}
	return nil
}
