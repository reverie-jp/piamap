package mapper

import (
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
)

// PianoRow は ListPianosInBBoxRow / ListPianosNearbyRow / GetPianoByIDRow の共通フィールドを持つ
// インターフェイス。3つの sqlc 生成型はカラムが同じなので、同じマッパーで扱う。
// sqlc が共通型を生成しないので Go 側で薄く吸収する。
//
// 集約ヘルパとして、各 row 型から *entity.Piano に変換する関数を定義する。

func ToPianoFromGetRow(row *sqlc.GetPianoByIDRow) *entity.Piano {
	if row == nil {
		return nil
	}
	return &entity.Piano{
		ID:                  row.ID,
		Name:                row.Name,
		Description:         row.Description,
		Location:            entity.LatLng{Latitude: row.Latitude, Longitude: row.Longitude},
		Address:             row.Address,
		Prefecture:          row.Prefecture,
		City:                row.City,
		Kind:                entity.PianoKind(row.Kind),
		VenueType:           row.VenueType,
		PianoType:           entity.PianoType(row.PianoType),
		PianoBrand:          row.PianoBrand,
		PianoModel:          row.PianoModel,
		ManufactureYear:     row.ManufactureYear,
		Hours:               row.Hours,
		Status:              entity.PianoStatus(row.Status),
		Availability:        entity.PianoAvailability(row.Availability),
		AvailabilityNote:    row.AvailabilityNote,
		InstallTime:         row.InstallTime,
		RemoveTime:          row.RemoveTime,
		CreatorUserID:       row.CreatorUserID,
		PostCount:           row.PostCount,
		RatingCount:         row.RatingCount,
		RatingSum:           row.RatingSum,
		AmbientNoiseCount:   row.AmbientNoiseCount,
		AmbientNoiseSum:     row.AmbientNoiseSum,
		FootTrafficCount:    row.FootTrafficCount,
		FootTrafficSum:      row.FootTrafficSum,
		ResonanceCount:      row.ResonanceCount,
		ResonanceSum:        row.ResonanceSum,
		KeyTouchWeightCount: row.KeyTouchWeightCount,
		KeyTouchWeightSum:   row.KeyTouchWeightSum,
		TuningQualityCount:  row.TuningQualityCount,
		TuningQualitySum:    row.TuningQualitySum,
		WishlistCount:       row.WishlistCount,
		VisitedCount:        row.VisitedCount,
		FavoriteCount:       row.FavoriteCount,
		CreateTime:          row.CreateTime,
		UpdateTime:          row.UpdateTime,
		DistanceM:           row.DistanceM,
	}
}

func ToPianoFromBBoxRow(row *sqlc.ListPianosInBBoxRow) *entity.Piano {
	if row == nil {
		return nil
	}
	return &entity.Piano{
		ID:                  row.ID,
		Name:                row.Name,
		Description:         row.Description,
		Location:            entity.LatLng{Latitude: row.Latitude, Longitude: row.Longitude},
		Address:             row.Address,
		Prefecture:          row.Prefecture,
		City:                row.City,
		Kind:                entity.PianoKind(row.Kind),
		VenueType:           row.VenueType,
		PianoType:           entity.PianoType(row.PianoType),
		PianoBrand:          row.PianoBrand,
		PianoModel:          row.PianoModel,
		ManufactureYear:     row.ManufactureYear,
		Hours:               row.Hours,
		Status:              entity.PianoStatus(row.Status),
		Availability:        entity.PianoAvailability(row.Availability),
		AvailabilityNote:    row.AvailabilityNote,
		InstallTime:         row.InstallTime,
		RemoveTime:          row.RemoveTime,
		CreatorUserID:       row.CreatorUserID,
		PostCount:           row.PostCount,
		RatingCount:         row.RatingCount,
		RatingSum:           row.RatingSum,
		AmbientNoiseCount:   row.AmbientNoiseCount,
		AmbientNoiseSum:     row.AmbientNoiseSum,
		FootTrafficCount:    row.FootTrafficCount,
		FootTrafficSum:      row.FootTrafficSum,
		ResonanceCount:      row.ResonanceCount,
		ResonanceSum:        row.ResonanceSum,
		KeyTouchWeightCount: row.KeyTouchWeightCount,
		KeyTouchWeightSum:   row.KeyTouchWeightSum,
		TuningQualityCount:  row.TuningQualityCount,
		TuningQualitySum:    row.TuningQualitySum,
		WishlistCount:       row.WishlistCount,
		VisitedCount:        row.VisitedCount,
		FavoriteCount:       row.FavoriteCount,
		CreateTime:          row.CreateTime,
		UpdateTime:          row.UpdateTime,
		DistanceM:           row.DistanceM,
	}
}

func ToPianoFromNearbyRow(row *sqlc.ListPianosNearbyRow) *entity.Piano {
	if row == nil {
		return nil
	}
	return &entity.Piano{
		ID:                  row.ID,
		Name:                row.Name,
		Description:         row.Description,
		Location:            entity.LatLng{Latitude: row.Latitude, Longitude: row.Longitude},
		Address:             row.Address,
		Prefecture:          row.Prefecture,
		City:                row.City,
		Kind:                entity.PianoKind(row.Kind),
		VenueType:           row.VenueType,
		PianoType:           entity.PianoType(row.PianoType),
		PianoBrand:          row.PianoBrand,
		PianoModel:          row.PianoModel,
		ManufactureYear:     row.ManufactureYear,
		Hours:               row.Hours,
		Status:              entity.PianoStatus(row.Status),
		Availability:        entity.PianoAvailability(row.Availability),
		AvailabilityNote:    row.AvailabilityNote,
		InstallTime:         row.InstallTime,
		RemoveTime:          row.RemoveTime,
		CreatorUserID:       row.CreatorUserID,
		PostCount:           row.PostCount,
		RatingCount:         row.RatingCount,
		RatingSum:           row.RatingSum,
		AmbientNoiseCount:   row.AmbientNoiseCount,
		AmbientNoiseSum:     row.AmbientNoiseSum,
		FootTrafficCount:    row.FootTrafficCount,
		FootTrafficSum:      row.FootTrafficSum,
		ResonanceCount:      row.ResonanceCount,
		ResonanceSum:        row.ResonanceSum,
		KeyTouchWeightCount: row.KeyTouchWeightCount,
		KeyTouchWeightSum:   row.KeyTouchWeightSum,
		TuningQualityCount:  row.TuningQualityCount,
		TuningQualitySum:    row.TuningQualitySum,
		WishlistCount:       row.WishlistCount,
		VisitedCount:        row.VisitedCount,
		FavoriteCount:       row.FavoriteCount,
		CreateTime:          row.CreateTime,
		UpdateTime:          row.UpdateTime,
		DistanceM:           row.DistanceM,
	}
}

func ToPianoFromTextSearchRow(row *sqlc.SearchPianosByTextRow) *entity.Piano {
	if row == nil {
		return nil
	}
	return &entity.Piano{
		ID:                  row.ID,
		Name:                row.Name,
		Description:         row.Description,
		Location:            entity.LatLng{Latitude: row.Latitude, Longitude: row.Longitude},
		Address:             row.Address,
		Prefecture:          row.Prefecture,
		City:                row.City,
		Kind:                entity.PianoKind(row.Kind),
		VenueType:           row.VenueType,
		PianoType:           entity.PianoType(row.PianoType),
		PianoBrand:          row.PianoBrand,
		PianoModel:          row.PianoModel,
		ManufactureYear:     row.ManufactureYear,
		Hours:               row.Hours,
		Status:              entity.PianoStatus(row.Status),
		Availability:        entity.PianoAvailability(row.Availability),
		AvailabilityNote:    row.AvailabilityNote,
		InstallTime:         row.InstallTime,
		RemoveTime:          row.RemoveTime,
		CreatorUserID:       row.CreatorUserID,
		PostCount:           row.PostCount,
		RatingCount:         row.RatingCount,
		RatingSum:           row.RatingSum,
		AmbientNoiseCount:   row.AmbientNoiseCount,
		AmbientNoiseSum:     row.AmbientNoiseSum,
		FootTrafficCount:    row.FootTrafficCount,
		FootTrafficSum:      row.FootTrafficSum,
		ResonanceCount:      row.ResonanceCount,
		ResonanceSum:        row.ResonanceSum,
		KeyTouchWeightCount: row.KeyTouchWeightCount,
		KeyTouchWeightSum:   row.KeyTouchWeightSum,
		TuningQualityCount:  row.TuningQualityCount,
		TuningQualitySum:    row.TuningQualitySum,
		WishlistCount:       row.WishlistCount,
		VisitedCount:        row.VisitedCount,
		FavoriteCount:       row.FavoriteCount,
		CreateTime:          row.CreateTime,
		UpdateTime:          row.UpdateTime,
		DistanceM:           row.DistanceM,
	}
}
