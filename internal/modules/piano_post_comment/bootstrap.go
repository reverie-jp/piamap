package piano_post_comment

import (
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_comment/v1/piano_post_commentv1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/handler"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/usecase"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
)

func InitModule(
	commentGateway gateway.Gateway,
	postGateway postgw.Gateway,
	userGateway usergw.Gateway,
) piano_post_commentv1connect.PianoPostCommentServiceHandler {
	return handler.New(
		usecase.NewCreatePianoPostComment(commentGateway, postGateway),
		usecase.NewListPianoPostComments(commentGateway, postGateway, userGateway),
		usecase.NewDeletePianoPostComment(commentGateway),
	)
}
