package usecase

import "github.com/reverie-jp/piamap/internal/modules/piano/gateway"

type SearchPianosOutput struct {
	Views []*gateway.PianoView
}
