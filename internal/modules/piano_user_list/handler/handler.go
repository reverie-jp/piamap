package handler

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1/piano_user_listv1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/usecase"
)

type Handler struct {
	piano_user_listv1connect.UnimplementedPianoUserListServiceHandler
	add        *usecase.AddPianoToUserList
	remove     *usecase.RemovePianoFromUserList
	listPianos *usecase.ListUserListPianos
	getMine    *usecase.GetMyPianoUserLists
}

func New(
	add *usecase.AddPianoToUserList,
	remove *usecase.RemovePianoFromUserList,
	listPianos *usecase.ListUserListPianos,
	getMine *usecase.GetMyPianoUserLists,
) *Handler {
	return &Handler{add: add, remove: remove, listPianos: listPianos, getMine: getMine}
}
