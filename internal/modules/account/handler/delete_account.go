package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "github.com/reverie-jp/piamap/internal/gen/pb/account/v1"
	"github.com/reverie-jp/piamap/internal/modules/account/adapter"
)

func (h *Handler) DeleteAccount(ctx context.Context, req *connect.Request[accountv1.DeleteAccountRequest]) (*connect.Response[accountv1.DeleteAccountResponse], error) {
	input, err := adapter.FromDeleteAccountRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.deleteAccount.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToDeleteAccountResponse(), nil
}
