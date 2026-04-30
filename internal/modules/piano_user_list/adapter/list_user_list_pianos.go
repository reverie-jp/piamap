package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	piano_user_listv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1"
	pianoadapter "github.com/reverie-jp/piamap/internal/modules/piano/adapter"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromListRequest(ctx context.Context, req *connect.Request[piano_user_listv1.ListUserListPianosRequest]) (usecase.ListUserListPianosInput, error) {
	requesterID, _ := interceptor.UserIDFromContext(ctx)
	customID, err := resourcename.ParseUser(req.Msg.GetParent())
	if err != nil {
		return usecase.ListUserListPianosInput{}, err
	}
	kind, ok := fromPbListKind(req.Msg.GetListKind())
	if !ok {
		return usecase.ListUserListPianosInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid list_kind"))
	}
	in := usecase.ListUserListPianosInput{
		RequesterID:  requesterID,
		UserCustomID: customID,
		ListKind:     kind,
	}
	if req.Msg.PageSize != nil {
		in.PageSize = *req.Msg.PageSize
	}
	if t := req.Msg.GetPageToken(); t != "" {
		id, err := ulid.Parse(t)
		if err != nil {
			return usecase.ListUserListPianosInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid page_token"))
		}
		in.AfterPianoID = &id
	}
	return in, nil
}

func ToListResponse(output *usecase.ListUserListPianosOutput) *connect.Response[piano_user_listv1.ListUserListPianosResponse] {
	pianos := make([]*pianov1.Piano, 0, len(output.Views))
	for _, v := range output.Views {
		if pb := pianoadapter.ToPiano(v); pb != nil {
			pianos = append(pianos, pb)
		}
	}
	resp := &piano_user_listv1.ListUserListPianosResponse{Pianos: pianos}
	if output.NextPianoID != nil {
		resp.NextPageToken = output.NextPianoID.String()
	}
	return connect.NewResponse(resp)
}
