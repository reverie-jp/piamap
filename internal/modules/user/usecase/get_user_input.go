package usecase

import (
	"errors"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// 公開プロフィール取得 (custom_id 経由)。RequesterID は guest 可。
type GetUserInput struct {
	RequesterID    ulid.ULID
	TargetCustomID string
}

func (i GetUserInput) Validate() error {
	if i.TargetCustomID == "" {
		return errors.New("target custom_id is required")
	}
	return nil
}
