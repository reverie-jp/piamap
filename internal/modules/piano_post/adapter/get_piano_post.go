package adapter

import (
	"context"

	"connectrpc.com/connect"

	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
)

func FromGetPianoPostRequest(ctx context.Context, req *connect.Request[piano_postv1.GetPianoPostRequest]) (usecase.GetPianoPostInput, error) {
	requesterID, _ := interceptor.UserIDFromContext(ctx)
	_, postID, err := resourcename.ParsePianoPost(req.Msg.GetName())
	if err != nil {
		return usecase.GetPianoPostInput{}, err
	}
	return usecase.GetPianoPostInput{
		RequesterID: requesterID,
		PostID:      postID,
	}, nil
}

func ToGetPianoPostResponse(output *usecase.GetPianoPostOutput) *connect.Response[piano_postv1.GetPianoPostResponse] {
	return connect.NewResponse(&piano_postv1.GetPianoPostResponse{
		PianoPost: ToPianoPost(output.View),
	})
}
