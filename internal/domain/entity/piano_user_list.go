package entity

type PianoListKind string

const (
	PianoListKindWishlist PianoListKind = "wishlist"
	PianoListKindVisited  PianoListKind = "visited"
	PianoListKindFavorite PianoListKind = "favorite"
)

func (k PianoListKind) Valid() bool {
	switch k {
	case PianoListKindWishlist, PianoListKindVisited, PianoListKindFavorite:
		return true
	}
	return false
}
