package handler

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1/piano_postv1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
)

type Handler struct {
	piano_postv1connect.UnimplementedPianoPostServiceHandler
	getPianoPost    *usecase.GetPianoPost
	listPianoPosts  *usecase.ListPianoPosts
	createPianoPost *usecase.CreatePianoPost
	updatePianoPost *usecase.UpdatePianoPost
	deletePianoPost *usecase.DeletePianoPost
}

func New(
	getPianoPost *usecase.GetPianoPost,
	listPianoPosts *usecase.ListPianoPosts,
	createPianoPost *usecase.CreatePianoPost,
	updatePianoPost *usecase.UpdatePianoPost,
	deletePianoPost *usecase.DeletePianoPost,
) *Handler {
	return &Handler{
		getPianoPost:    getPianoPost,
		listPianoPosts:  listPianoPosts,
		createPianoPost: createPianoPost,
		updatePianoPost: updatePianoPost,
		deletePianoPost: deletePianoPost,
	}
}
