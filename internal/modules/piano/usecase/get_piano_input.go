package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type GetPianoInput struct {
	RequesterID ulid.ULID // ゲスト可
	PianoID     ulid.ULID
}

func (i GetPianoInput) Validate() error {
	if i.PianoID.IsZero() {
		return errors.New("piano id is required")
	}
	return nil
}
