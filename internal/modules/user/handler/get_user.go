package handler

import (
	"context"

	"connectrpc.com/connect"

	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/adapter"
)

func (h *Handler) GetUser(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (*connect.Response[userv1.GetUserResponse], error) {
	input, err := adapter.FromGetUserRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getUser.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetUserResponse(output), nil
}
