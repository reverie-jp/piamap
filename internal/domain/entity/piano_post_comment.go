package entity

import (
	"time"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type PianoPostComment struct {
	ID              ulid.ULID
	PianoPostID     ulid.ULID
	UserID          ulid.ULID
	ParentCommentID *ulid.ULID
	Body            string
	CreateTime      time.Time
	UpdateTime      time.Time
}
