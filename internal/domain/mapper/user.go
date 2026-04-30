package mapper

import (
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
)

func ToUser(row *sqlc.User) *entity.User {
	if row == nil {
		return nil
	}
	return &entity.User{
		ID:                  row.ID,
		CustomID:            row.CustomID,
		CustomIDChangeTime:  row.CustomIDChangeTime,
		DisplayName:         row.DisplayName,
		Biography:           row.Biography,
		AvatarURL:           row.AvatarUrl,
		Hometown:            row.Hometown,
		PianoStartedYear:    row.PianoStartedYear,
		YearsOfExperience:   row.YearsOfExperience,
		PostCount:           row.PostCount,
		EditCount:           row.EditCount,
		ReportReceivedCount: row.ReportReceivedCount,
		CreateTime:          row.CreateTime,
		UpdateTime:          row.UpdateTime,
	}
}
