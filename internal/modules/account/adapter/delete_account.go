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

func FromDeleteAccountRequest(ctx context.Context, req *connect.Request[accountv1.DeleteAccountRequest]) (usecase.DeleteAccountInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.DeleteAccountInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	return usecase.DeleteAccountInput{
		UserID:          userID,
		ConfirmCustomID: req.Msg.ConfirmCustomId,
	}, nil
}

func ToDeleteAccountResponse() *connect.Response[accountv1.DeleteAccountResponse] {
	return connect.NewResponse(&accountv1.DeleteAccountResponse{})
}
