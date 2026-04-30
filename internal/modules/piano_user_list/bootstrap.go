package piano_user_list

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1/piano_user_listv1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/handler"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/usecase"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
)

func InitModule(
	listGateway gateway.Gateway,
	pianoGateway pianogw.Gateway,
	userGateway usergw.Gateway,
) piano_user_listv1connect.PianoUserListServiceHandler {
	return handler.New(
		usecase.NewAddPianoToUserList(listGateway, pianoGateway),
		usecase.NewRemovePianoFromUserList(listGateway),
		usecase.NewListUserListPianos(listGateway, userGateway, pianoGateway),
		usecase.NewGetMyPianoUserLists(listGateway),
	)
}
