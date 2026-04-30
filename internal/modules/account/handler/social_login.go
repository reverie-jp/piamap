package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "github.com/reverie-jp/piamap/internal/gen/pb/account/v1"
	"github.com/reverie-jp/piamap/internal/modules/account/adapter"
)

func (h *Handler) SocialLogin(ctx context.Context, req *connect.Request[accountv1.SocialLoginRequest]) (*connect.Response[accountv1.SocialLoginResponse], error) {
	input, err := adapter.FromSocialLoginRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.socialLogin.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToSocialLoginResponse(output), nil
}
