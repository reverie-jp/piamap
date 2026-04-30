package handler

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/account/v1/accountv1connect"
	"github.com/reverie-jp/piamap/internal/modules/account/usecase"
)

type Handler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	socialLogin   *usecase.SocialLogin
	refreshToken  *usecase.RefreshToken
	logout        *usecase.Logout
	deleteAccount *usecase.DeleteAccount
}

func New(
	socialLogin *usecase.SocialLogin,
	refreshToken *usecase.RefreshToken,
	logout *usecase.Logout,
	deleteAccount *usecase.DeleteAccount,
) *Handler {
	return &Handler{
		socialLogin:   socialLogin,
		refreshToken:  refreshToken,
		logout:        logout,
		deleteAccount: deleteAccount,
	}
}
