package user

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/user/v1/userv1connect"
	"github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/modules/user/handler"
	"github.com/reverie-jp/piamap/internal/modules/user/usecase"
)

func InitModule(userGateway gateway.Gateway) userv1connect.UserServiceHandler {
	return handler.New(
		usecase.NewGetUser(userGateway),
		usecase.NewGetMyUser(userGateway),
		usecase.NewUpdateUser(userGateway),
	)
}
