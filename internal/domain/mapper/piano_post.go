package mapper

import (
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
)

func ToPianoPostFromGetRow(row *sqlc.GetPianoPostByIDRow) *entity.PianoPost {
	if row == nil {
		return nil
	}
	return &entity.PianoPost{
		ID:             row.ID,
		UserID:         row.UserID,
		PianoID:        row.PianoID,
		VisitTime:      row.VisitTime,
		Rating:         row.Rating,
		Body:           row.Body,
		AmbientNoise:   row.AmbientNoise,
		FootTraffic:    row.FootTraffic,
		Resonance:      row.Resonance,
		KeyTouchWeight: row.KeyTouchWeight,
		TuningQuality:  row.TuningQuality,
		Visibility:     entity.PostVisibility(row.Visibility),
		CommentCount:   row.CommentCount,
		LikeCount:      row.LikeCount,
		CreateTime:     row.CreateTime,
		UpdateTime:     row.UpdateTime,
	}
}

func ToPianoPostFromListByPianoRow(row *sqlc.ListPianoPostsByPianoRow) *entity.PianoPost {
	if row == nil {
		return nil
	}
	return &entity.PianoPost{
		ID:             row.ID,
		UserID:         row.UserID,
		PianoID:        row.PianoID,
		VisitTime:      row.VisitTime,
		Rating:         row.Rating,
		Body:           row.Body,
		AmbientNoise:   row.AmbientNoise,
		FootTraffic:    row.FootTraffic,
		Resonance:      row.Resonance,
		KeyTouchWeight: row.KeyTouchWeight,
		TuningQuality:  row.TuningQuality,
		Visibility:     entity.PostVisibility(row.Visibility),
		CommentCount:   row.CommentCount,
		LikeCount:      row.LikeCount,
		CreateTime:     row.CreateTime,
		UpdateTime:     row.UpdateTime,
	}
}

func ToPianoPostFromListByUserRow(row *sqlc.ListPianoPostsByUserRow) *entity.PianoPost {
	if row == nil {
		return nil
	}
	return &entity.PianoPost{
		ID:             row.ID,
		UserID:         row.UserID,
		PianoID:        row.PianoID,
		VisitTime:      row.VisitTime,
		Rating:         row.Rating,
		Body:           row.Body,
		AmbientNoise:   row.AmbientNoise,
		FootTraffic:    row.FootTraffic,
		Resonance:      row.Resonance,
		KeyTouchWeight: row.KeyTouchWeight,
		TuningQuality:  row.TuningQuality,
		Visibility:     entity.PostVisibility(row.Visibility),
		CommentCount:   row.CommentCount,
		LikeCount:      row.LikeCount,
		CreateTime:     row.CreateTime,
		UpdateTime:     row.UpdateTime,
	}
}

func ToPianoPostFromListPublicRow(row *sqlc.ListPublicPianoPostsRow) *entity.PianoPost {
	if row == nil {
		return nil
	}
	return &entity.PianoPost{
		ID:             row.ID,
		UserID:         row.UserID,
		PianoID:        row.PianoID,
		VisitTime:      row.VisitTime,
		Rating:         row.Rating,
		Body:           row.Body,
		AmbientNoise:   row.AmbientNoise,
		FootTraffic:    row.FootTraffic,
		Resonance:      row.Resonance,
		KeyTouchWeight: row.KeyTouchWeight,
		TuningQuality:  row.TuningQuality,
		Visibility:     entity.PostVisibility(row.Visibility),
		CommentCount:   row.CommentCount,
		LikeCount:      row.LikeCount,
		CreateTime:     row.CreateTime,
		UpdateTime:     row.UpdateTime,
	}
}
