package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type GetMyUserInput struct {
	RequesterID ulid.ULID
}

func (i GetMyUserInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester is required")
	}
	return nil
}
