package handler

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/piano/v1/pianov1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
)

type Handler struct {
	pianov1connect.UnimplementedPianoServiceHandler
	getPiano     *usecase.GetPiano
	searchPianos *usecase.SearchPianos
	createPiano  *usecase.CreatePiano
	updatePiano  *usecase.UpdatePiano
}

func New(
	getPiano *usecase.GetPiano,
	searchPianos *usecase.SearchPianos,
	createPiano *usecase.CreatePiano,
	updatePiano *usecase.UpdatePiano,
) *Handler {
	return &Handler{
		getPiano:     getPiano,
		searchPianos: searchPianos,
		createPiano:  createPiano,
		updatePiano:  updatePiano,
	}
}
