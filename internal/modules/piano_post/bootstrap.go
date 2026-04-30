package piano_post

import (
	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1/piano_postv1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/handler"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
)

func InitModule(
	postGateway gateway.Gateway,
	pianoGateway pianogw.Gateway,
	userGateway usergw.Gateway,
	tx transaction.Runner,
) piano_postv1connect.PianoPostServiceHandler {
	return handler.New(
		usecase.NewGetPianoPost(postGateway),
		usecase.NewListPianoPosts(postGateway, userGateway),
		usecase.NewCreatePianoPost(postGateway, pianoGateway, userGateway, tx),
		usecase.NewUpdatePianoPost(postGateway),
		usecase.NewDeletePianoPost(postGateway),
	)
}
