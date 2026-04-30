package handler

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/user/v1/userv1connect"
	"github.com/reverie-jp/piamap/internal/modules/user/usecase"
)

type Handler struct {
	userv1connect.UnimplementedUserServiceHandler
	getUser    *usecase.GetUser
	getMyUser  *usecase.GetMyUser
	updateUser *usecase.UpdateUser
}

func New(
	getUser *usecase.GetUser,
	getMyUser *usecase.GetMyUser,
	updateUser *usecase.UpdateUser,
) *Handler {
	return &Handler{
		getUser:    getUser,
		getMyUser:  getMyUser,
		updateUser: updateUser,
	}
}
