package usecase

import "github.com/reverie-jp/piamap/internal/modules/piano/gateway"

type CreatePianoOutput struct {
	View *gateway.PianoView
}
