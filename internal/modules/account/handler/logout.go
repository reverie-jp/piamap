package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "github.com/reverie-jp/piamap/internal/gen/pb/account/v1"
	"github.com/reverie-jp/piamap/internal/modules/account/adapter"
)

func (h *Handler) Logout(ctx context.Context, req *connect.Request[accountv1.LogoutRequest]) (*connect.Response[accountv1.LogoutResponse], error) {
	input, err := adapter.FromLogoutRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.logout.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToLogoutResponse(), nil
}
