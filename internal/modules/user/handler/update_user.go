package handler

import (
	"context"

	"connectrpc.com/connect"

	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/adapter"
)

func (h *Handler) UpdateUser(ctx context.Context, req *connect.Request[userv1.UpdateUserRequest]) (*connect.Response[userv1.UpdateUserResponse], error) {
	input, err := adapter.FromUpdateUserRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.updateUser.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToUpdateUserResponse(output), nil
}
