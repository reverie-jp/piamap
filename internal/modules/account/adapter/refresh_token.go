package adapter

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "github.com/reverie-jp/piamap/internal/gen/pb/account/v1"
	"github.com/reverie-jp/piamap/internal/modules/account/usecase"
)

func FromRefreshTokenRequest(_ context.Context, req *connect.Request[accountv1.RefreshTokenRequest]) (usecase.RefreshTokenInput, error) {
	return usecase.RefreshTokenInput{RefreshToken: req.Msg.RefreshToken}, nil
}

func ToRefreshTokenResponse(output *usecase.RefreshTokenOutput) *connect.Response[accountv1.RefreshTokenResponse] {
	return connect.NewResponse(&accountv1.RefreshTokenResponse{
		TokenPair: &accountv1.TokenPair{
			AccessToken:  output.AccessToken,
			RefreshToken: output.RefreshToken,
		},
	})
}
