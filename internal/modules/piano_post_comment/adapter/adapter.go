package adapter

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	piano_post_commentv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post_comment/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func ToPianoPostComment(view *gateway.PianoPostCommentView) *piano_post_commentv1.PianoPostComment {
	if view == nil || view.Comment == nil {
		return nil
	}
	c := view.Comment
	pb := &piano_post_commentv1.PianoPostComment{
		Name:              resourcename.FormatPianoPostComment(view.PianoID, c.PianoPostID, c.ID),
		Body:              c.Body,
		CreateTime:        timestamppb.New(c.CreateTime),
		UpdateTime:        timestamppb.New(c.UpdateTime),
		AuthorDisplayName: view.AuthorDisplayName,
		AuthorAvatarUrl:   view.AuthorAvatarURL,
	}
	if view.AuthorCustomID != "" {
		pb.Author = resourcename.FormatUser(view.AuthorCustomID)
	}
	if c.ParentCommentID != nil {
		pn := resourcename.FormatPianoPostComment(view.PianoID, c.PianoPostID, *c.ParentCommentID)
		pb.ParentComment = &pn
	}
	return pb
}

func FromCreateRequest(ctx context.Context, req *connect.Request[piano_post_commentv1.CreatePianoPostCommentRequest]) (usecase.CreatePianoPostCommentInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.CreatePianoPostCommentInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	_, postID, err := resourcename.ParsePianoPost(req.Msg.GetParent())
	if err != nil {
		return usecase.CreatePianoPostCommentInput{}, err
	}
	in := usecase.CreatePianoPostCommentInput{
		RequesterID: requesterID,
		PianoPostID: postID,
	}
	if c := req.Msg.GetPianoPostComment(); c != nil {
		in.Body = c.GetBody()
		if pc := c.GetParentComment(); pc != "" {
			_, _, parentID, err := resourcename.ParsePianoPostComment(pc)
			if err != nil {
				return usecase.CreatePianoPostCommentInput{}, err
			}
			in.ParentCommentID = &parentID
		}
	}
	return in, nil
}

func ToCreateResponse(output *usecase.CreatePianoPostCommentOutput) *connect.Response[piano_post_commentv1.CreatePianoPostCommentResponse] {
	return connect.NewResponse(&piano_post_commentv1.CreatePianoPostCommentResponse{
		PianoPostComment: ToPianoPostComment(output.View),
	})
}

func FromListRequest(ctx context.Context, req *connect.Request[piano_post_commentv1.ListPianoPostCommentsRequest]) (usecase.ListPianoPostCommentsInput, error) {
	requesterID, _ := interceptor.UserIDFromContext(ctx)
	parent := req.Msg.GetParent()
	if parent == "" {
		return usecase.ListPianoPostCommentsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("parent is required"))
	}
	in := usecase.ListPianoPostCommentsInput{RequesterID: requesterID}
	switch {
	case strings.HasPrefix(parent, "pianos/"):
		_, postID, err := resourcename.ParsePianoPost(parent)
		if err != nil {
			return usecase.ListPianoPostCommentsInput{}, err
		}
		in.ParentKind = usecase.ListParentPost
		in.PianoPostID = postID
	case strings.HasPrefix(parent, "users/"):
		customID, err := resourcename.ParseUser(parent)
		if err != nil {
			return usecase.ListPianoPostCommentsInput{}, err
		}
		in.ParentKind = usecase.ListParentUser
		in.UserCustomID = customID
	default:
		return usecase.ListPianoPostCommentsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid parent"))
	}
	if req.Msg.PageSize != nil {
		in.PageSize = *req.Msg.PageSize
	}
	if t := req.Msg.GetPageToken(); t != "" {
		id, err := ulid.Parse(t)
		if err != nil {
			return usecase.ListPianoPostCommentsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid page_token"))
		}
		in.AfterID = &id
	}
	return in, nil
}

func ToListResponse(output *usecase.ListPianoPostCommentsOutput) *connect.Response[piano_post_commentv1.ListPianoPostCommentsResponse] {
	resp := &piano_post_commentv1.ListPianoPostCommentsResponse{
		PianoPostComments: make([]*piano_post_commentv1.PianoPostComment, 0, len(output.Views)),
	}
	for _, v := range output.Views {
		if pb := ToPianoPostComment(v); pb != nil {
			resp.PianoPostComments = append(resp.PianoPostComments, pb)
		}
	}
	if output.NextID != nil {
		resp.NextPageToken = output.NextID.String()
	}
	return connect.NewResponse(resp)
}

func FromDeleteRequest(ctx context.Context, req *connect.Request[piano_post_commentv1.DeletePianoPostCommentRequest]) (usecase.DeletePianoPostCommentInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.DeletePianoPostCommentInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	_, _, commentID, err := resourcename.ParsePianoPostComment(req.Msg.GetName())
	if err != nil {
		return usecase.DeletePianoPostCommentInput{}, err
	}
	return usecase.DeletePianoPostCommentInput{RequesterID: requesterID, CommentID: commentID}, nil
}

func ToDeleteResponse(_ *usecase.DeletePianoPostCommentOutput) *connect.Response[piano_post_commentv1.DeletePianoPostCommentResponse] {
	return connect.NewResponse(&piano_post_commentv1.DeletePianoPostCommentResponse{})
}
