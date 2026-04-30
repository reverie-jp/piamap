package entity

import (
	"time"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type RefreshToken struct {
	ID         ulid.ULID
	UserID     ulid.ULID
	TokenHash  string
	ExpireTime time.Time
	CreateTime time.Time
}
