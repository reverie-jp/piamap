package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	piano_post_likev1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post_like/v1"
	postadapter "github.com/reverie-jp/piamap/internal/modules/piano_post/adapter"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromLikeRequest(ctx context.Context, req *connect.Request[piano_post_likev1.LikePianoPostRequest]) (usecase.LikePianoPostInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.LikePianoPostInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	_, postID, err := resourcename.ParsePianoPost(req.Msg.GetParent())
	if err != nil {
		return usecase.LikePianoPostInput{}, err
	}
	return usecase.LikePianoPostInput{RequesterID: requesterID, PostID: postID}, nil
}

func ToLikeResponse(_ *usecase.LikePianoPostOutput) *connect.Response[piano_post_likev1.LikePianoPostResponse] {
	return connect.NewResponse(&piano_post_likev1.LikePianoPostResponse{})
}

func FromUnlikeRequest(ctx context.Context, req *connect.Request[piano_post_likev1.UnlikePianoPostRequest]) (usecase.UnlikePianoPostInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UnlikePianoPostInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	_, postID, err := resourcename.ParsePianoPost(req.Msg.GetParent())
	if err != nil {
		return usecase.UnlikePianoPostInput{}, err
	}
	return usecase.UnlikePianoPostInput{RequesterID: requesterID, PostID: postID}, nil
}

func ToUnlikeResponse(_ *usecase.UnlikePianoPostOutput) *connect.Response[piano_post_likev1.UnlikePianoPostResponse] {
	return connect.NewResponse(&piano_post_likev1.UnlikePianoPostResponse{})
}

func FromListLikedRequest(ctx context.Context, req *connect.Request[piano_post_likev1.ListLikedPianoPostsRequest]) (usecase.ListLikedPianoPostsInput, error) {
	requesterID, _ := interceptor.UserIDFromContext(ctx)
	customID, err := resourcename.ParseUser(req.Msg.GetParent())
	if err != nil {
		return usecase.ListLikedPianoPostsInput{}, err
	}
	in := usecase.ListLikedPianoPostsInput{RequesterID: requesterID, UserCustomID: customID}
	if req.Msg.PageSize != nil {
		in.PageSize = *req.Msg.PageSize
	}
	if t := req.Msg.GetPageToken(); t != "" {
		id, err := ulid.Parse(t)
		if err != nil {
			return usecase.ListLikedPianoPostsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid page_token"))
		}
		in.AfterPostID = &id
	}
	return in, nil
}

func ToListLikedResponse(output *usecase.ListLikedPianoPostsOutput) *connect.Response[piano_post_likev1.ListLikedPianoPostsResponse] {
	posts := make([]*piano_postv1.PianoPost, 0, len(output.Views))
	for _, v := range output.Views {
		if pb := postadapter.ToPianoPost(v); pb != nil {
			posts = append(posts, pb)
		}
	}
	resp := &piano_post_likev1.ListLikedPianoPostsResponse{PianoPosts: posts}
	if output.NextPostID != nil {
		resp.NextPageToken = output.NextPostID.String()
	}
	return connect.NewResponse(resp)
}
