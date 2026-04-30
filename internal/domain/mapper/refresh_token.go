package mapper

import (
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
)

func ToRefreshToken(row *sqlc.RefreshToken) *entity.RefreshToken {
	if row == nil {
		return nil
	}
	return &entity.RefreshToken{
		ID:         row.ID,
		UserID:     row.UserID,
		TokenHash:  row.TokenHash,
		ExpireTime: row.ExpireTime,
		CreateTime: row.CreateTime,
	}
}
