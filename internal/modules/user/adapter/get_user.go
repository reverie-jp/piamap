package adapter

import (
	"context"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
)

func FromGetUserRequest(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (usecase.GetUserInput, error) {
	customID, err := resourcename.ParseUser(req.Msg.Name)
	if err != nil {
		return usecase.GetUserInput{}, err
	}
	input := usecase.GetUserInput{TargetCustomID: customID}
	if requesterID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = requesterID
	}
	return input, nil
}

func ToGetUserResponse(output *usecase.GetUserOutput) *connect.Response[userv1.GetUserResponse] {
	return connect.NewResponse(&userv1.GetUserResponse{User: ToUser(output.View)})
}
