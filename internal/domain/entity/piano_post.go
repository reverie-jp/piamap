package entity

import (
	"time"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type PostVisibility string

const (
	PostVisibilityPublic  PostVisibility = "public"
	PostVisibilityPrivate PostVisibility = "private"
)

// PianoPost は 1 訪問 = 1 投稿。rating または body のどちらかは必須、5 環境属性は任意。
type PianoPost struct {
	ID             ulid.ULID
	UserID         ulid.ULID
	PianoID        ulid.ULID
	VisitTime      time.Time
	Rating         *int16
	Body           *string
	AmbientNoise   *int16
	FootTraffic    *int16
	Resonance      *int16
	KeyTouchWeight *int16
	TuningQuality  *int16
	Visibility     PostVisibility
	CommentCount   int32
	LikeCount      int32
	CreateTime     time.Time
	UpdateTime     time.Time
}
