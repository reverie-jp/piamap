package usecase

import "github.com/reverie-jp/piamap/internal/modules/piano/gateway"

type GetPianoOutput struct {
	View *gateway.PianoView
}
