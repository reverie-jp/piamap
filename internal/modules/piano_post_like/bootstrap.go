package piano_post_like

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_like/v1/piano_post_likev1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/handler"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/usecase"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
)

func InitModule(
	likeGateway gateway.Gateway,
	postGateway postgw.Gateway,
	userGateway usergw.Gateway,
) piano_post_likev1connect.PianoPostLikeServiceHandler {
	return handler.New(
		usecase.NewLikePianoPost(likeGateway, postGateway),
		usecase.NewUnlikePianoPost(likeGateway),
		usecase.NewListLikedPianoPosts(likeGateway, userGateway, postGateway),
	)
}
