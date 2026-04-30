package mapper

import (
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
)

func ToPianoEdit(row *sqlc.PianoEdit) *entity.PianoEdit {
	if row == nil {
		return nil
	}
	return &entity.PianoEdit{
		ID:           row.ID,
		PianoID:      row.PianoID,
		EditorUserID: row.EditorUserID,
		Operation:    entity.PianoEditOperation(row.Operation),
		Changes:      row.Changes,
		Summary:      row.Summary,
		CreateTime:   row.CreateTime,
	}
}
