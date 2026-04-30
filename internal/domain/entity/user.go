package entity

import (
	"time"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// 信頼ライン: 登録から N 日経過 + post + edit が M 件以上 + 通報 resolved 件数 = 0 で trusted。
const (
	TrustedAccountAgeDays    = 30
	TrustedContributionCount = 10
)

type User struct {
	ID                  ulid.ULID
	CustomID            string
	CustomIDChangeTime  *time.Time
	DisplayName         string
	Biography           *string
	AvatarURL           *string
	Hometown            *string
	PianoStartedYear    *int16
	YearsOfExperience   *int16
	PostCount           int32
	EditCount           int32
	ReportReceivedCount int32
	CreateTime          time.Time
	UpdateTime          time.Time
}

// IsTrusted は破壊的編集 (削除 / 大幅な座標移動 / 名前全文置換) のゲートに使う。
func (u *User) IsTrusted(now time.Time) bool {
	if u == nil {
		return false
	}
	if u.ReportReceivedCount > 0 {
		return false
	}
	if now.Sub(u.CreateTime) < TrustedAccountAgeDays*24*time.Hour {
		return false
	}
	return u.PostCount+u.EditCount >= TrustedContributionCount
}
