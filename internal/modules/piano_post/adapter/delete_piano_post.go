package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromDeletePianoPostRequest(ctx context.Context, req *connect.Request[piano_postv1.DeletePianoPostRequest]) (usecase.DeletePianoPostInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.DeletePianoPostInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	_, postID, err := resourcename.ParsePianoPost(req.Msg.GetName())
	if err != nil {
		return usecase.DeletePianoPostInput{}, err
	}
	return usecase.DeletePianoPostInput{
		RequesterID: requesterID,
		PostID:      postID,
	}, nil
}

func ToDeletePianoPostResponse() *connect.Response[piano_postv1.DeletePianoPostResponse] {
	return connect.NewResponse(&piano_postv1.DeletePianoPostResponse{})
}
