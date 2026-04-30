package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/usecase"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromGetMyUserRequest(ctx context.Context, _ *connect.Request[userv1.GetMyUserRequest]) (usecase.GetMyUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.GetMyUserInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	return usecase.GetMyUserInput{RequesterID: requesterID}, nil
}

func ToGetMyUserResponse(output *usecase.GetMyUserOutput) *connect.Response[userv1.GetMyUserResponse] {
	return connect.NewResponse(&userv1.GetMyUserResponse{User: ToUser(output.View)})
}
