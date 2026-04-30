package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	piano_user_listv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromRemoveRequest(ctx context.Context, req *connect.Request[piano_user_listv1.RemovePianoFromUserListRequest]) (usecase.RemovePianoFromUserListInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.RemovePianoFromUserListInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	pianoID, err := resourcename.ParsePiano(req.Msg.GetParent())
	if err != nil {
		return usecase.RemovePianoFromUserListInput{}, err
	}
	kind, ok := fromPbListKind(req.Msg.GetListKind())
	if !ok {
		return usecase.RemovePianoFromUserListInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid list_kind"))
	}
	return usecase.RemovePianoFromUserListInput{
		RequesterID: requesterID,
		PianoID:     pianoID,
		ListKind:    kind,
	}, nil
}

func ToRemoveResponse(_ *usecase.RemovePianoFromUserListOutput) *connect.Response[piano_user_listv1.RemovePianoFromUserListResponse] {
	return connect.NewResponse(&piano_user_listv1.RemovePianoFromUserListResponse{})
}
