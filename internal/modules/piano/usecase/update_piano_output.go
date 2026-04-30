package usecase

import "github.com/reverie-jp/piamap/internal/modules/piano/gateway"

type UpdatePianoOutput struct {
	View *gateway.PianoView
}
