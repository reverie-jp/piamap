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

func FromGetMyRequest(ctx context.Context, req *connect.Request[piano_user_listv1.GetMyPianoUserListsRequest]) (usecase.GetMyPianoUserListsInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.GetMyPianoUserListsInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	pianoID, err := resourcename.ParsePiano(req.Msg.GetParent())
	if err != nil {
		return usecase.GetMyPianoUserListsInput{}, err
	}
	return usecase.GetMyPianoUserListsInput{
		RequesterID: requesterID,
		PianoID:     pianoID,
	}, nil
}

func ToGetMyResponse(output *usecase.GetMyPianoUserListsOutput) *connect.Response[piano_user_listv1.GetMyPianoUserListsResponse] {
	pb := make([]piano_user_listv1.PianoListKind, 0, len(output.ListKinds))
	for _, k := range output.ListKinds {
		pb = append(pb, toPbListKind(k))
	}
	return connect.NewResponse(&piano_user_listv1.GetMyPianoUserListsResponse{ListKinds: pb})
}
