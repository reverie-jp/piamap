package adapter

import (
	"github.com/reverie-jp/piamap/internal/domain/entity"
	piano_user_listv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1"
)

func toPbListKind(k entity.PianoListKind) piano_user_listv1.PianoListKind {
	switch k {
	case entity.PianoListKindWishlist:
		return piano_user_listv1.PianoListKind_PIANO_LIST_KIND_WISHLIST
	case entity.PianoListKindVisited:
		return piano_user_listv1.PianoListKind_PIANO_LIST_KIND_VISITED
	case entity.PianoListKindFavorite:
		return piano_user_listv1.PianoListKind_PIANO_LIST_KIND_FAVORITE
	}
	return piano_user_listv1.PianoListKind_PIANO_LIST_KIND_UNSPECIFIED
}

func fromPbListKind(k piano_user_listv1.PianoListKind) (entity.PianoListKind, bool) {
	switch k {
	case piano_user_listv1.PianoListKind_PIANO_LIST_KIND_WISHLIST:
		return entity.PianoListKindWishlist, true
	case piano_user_listv1.PianoListKind_PIANO_LIST_KIND_VISITED:
		return entity.PianoListKindVisited, true
	case piano_user_listv1.PianoListKind_PIANO_LIST_KIND_FAVORITE:
		return entity.PianoListKindFavorite, true
	}
	return "", false
}
