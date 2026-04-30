package entity

import (
	"time"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

// PianoEdit は piano_edits テーブルの 1 行 = ピアノに対する 1 編集。
// 公開ログとして表示する用途のため、editor_user_id は SET NULL の可能性あり。
type PianoEdit struct {
	ID           ulid.ULID
	PianoID      ulid.ULID
	EditorUserID *ulid.ULID
	Operation    PianoEditOperation
	Changes      []byte // JSONB。表示時はそのまま JSON 文字列として返す。
	Summary      *string
	CreateTime   time.Time
}
