package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	accountv1 "github.com/reverie-jp/piamap/internal/gen/pb/account/v1"
	"github.com/reverie-jp/piamap/internal/modules/account/usecase"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromLogoutRequest(ctx context.Context, req *connect.Request[accountv1.LogoutRequest]) (usecase.LogoutInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.LogoutInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	return usecase.LogoutInput{
		UserID:       userID,
		RefreshToken: req.Msg.RefreshToken,
	}, nil
}

func ToLogoutResponse() *connect.Response[accountv1.LogoutResponse] {
	return connect.NewResponse(&accountv1.LogoutResponse{})
}
