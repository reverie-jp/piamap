package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type DeleteAccountInput struct {
	UserID          ulid.ULID
	ConfirmCustomID string
}

func (i DeleteAccountInput) Validate() error {
	if i.UserID.IsZero() {
		return errors.New("user_id is required")
	}
	if i.ConfirmCustomID == "" {
		return errors.New("confirm_custom_id is required")
	}
	return nil
}
