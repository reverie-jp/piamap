package piano

import (
	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano/v1/pianov1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano/handler"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
)

func InitModule(
	pianoGateway gateway.Gateway,
	userGateway usergw.Gateway,
	tx transaction.Runner,
) pianov1connect.PianoServiceHandler {
	return handler.New(
		usecase.NewGetPiano(pianoGateway),
		usecase.NewSearchPianos(pianoGateway),
		usecase.NewCreatePiano(pianoGateway, userGateway, tx),
		usecase.NewUpdatePiano(pianoGateway, userGateway, tx),
	)
}
