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

func FromAddRequest(ctx context.Context, req *connect.Request[piano_user_listv1.AddPianoToUserListRequest]) (usecase.AddPianoToUserListInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.AddPianoToUserListInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	pianoID, err := resourcename.ParsePiano(req.Msg.GetParent())
	if err != nil {
		return usecase.AddPianoToUserListInput{}, err
	}
	kind, ok := fromPbListKind(req.Msg.GetListKind())
	if !ok {
		return usecase.AddPianoToUserListInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid list_kind"))
	}
	return usecase.AddPianoToUserListInput{
		RequesterID: requesterID,
		PianoID:     pianoID,
		ListKind:    kind,
	}, nil
}

func ToAddResponse(_ *usecase.AddPianoToUserListOutput) *connect.Response[piano_user_listv1.AddPianoToUserListResponse] {
	return connect.NewResponse(&piano_user_listv1.AddPianoToUserListResponse{})
}
