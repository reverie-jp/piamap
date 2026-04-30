package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type LogoutInput struct {
	UserID       ulid.ULID
	RefreshToken string
}

func (i LogoutInput) Validate() error {
	if i.UserID.IsZero() {
		return errors.New("user_id is required")
	}
	if i.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}
	return nil
}
