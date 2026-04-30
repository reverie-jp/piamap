package handler

import (
	"context"

	"connectrpc.com/connect"

	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/adapter"
)

func (h *Handler) GetMyUser(ctx context.Context, req *connect.Request[userv1.GetMyUserRequest]) (*connect.Response[userv1.GetMyUserResponse], error) {
	input, err := adapter.FromGetMyUserRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getMyUser.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetMyUserResponse(output), nil
}
