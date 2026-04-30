package adapter

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromListPianoPostsRequest(ctx context.Context, req *connect.Request[piano_postv1.ListPianoPostsRequest]) (usecase.ListPianoPostsInput, error) {
	requesterID, _ := interceptor.UserIDFromContext(ctx)
	parent := req.Msg.GetParent()
	if parent == "" {
		return usecase.ListPianoPostsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("parent is required"))
	}

	input := usecase.ListPianoPostsInput{RequesterID: requesterID}
	switch {
	case parent == "-":
		input.ParentKind = usecase.ListParentGlobal
	case strings.HasPrefix(parent, "pianos/"):
		id, err := resourcename.ParsePiano(parent)
		if err != nil {
			return usecase.ListPianoPostsInput{}, err
		}
		input.ParentKind = usecase.ListParentPiano
		input.PianoID = id
	case strings.HasPrefix(parent, "users/"):
		customID, err := resourcename.ParseUser(parent)
		if err != nil {
			return usecase.ListPianoPostsInput{}, err
		}
		input.ParentKind = usecase.ListParentUser
		input.UserCustom = customID
	default:
		return usecase.ListPianoPostsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid parent"))
	}

	if req.Msg.PageSize != nil {
		input.PageSize = *req.Msg.PageSize
	}
	if t := req.Msg.GetPageToken(); t != "" {
		id, err := ulid.Parse(t)
		if err != nil {
			return usecase.ListPianoPostsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid page_token"))
		}
		input.AfterID = &id
	}
	return input, nil
}

func ToListPianoPostsResponse(output *usecase.ListPianoPostsOutput) *connect.Response[piano_postv1.ListPianoPostsResponse] {
	posts := make([]*piano_postv1.PianoPost, 0, len(output.Views))
	for _, v := range output.Views {
		if pb := ToPianoPost(v); pb != nil {
			posts = append(posts, pb)
		}
	}
	resp := &piano_postv1.ListPianoPostsResponse{PianoPosts: posts}
	if output.NextID != nil {
		resp.NextPageToken = output.NextID.String()
	}
	return connect.NewResponse(resp)
}
