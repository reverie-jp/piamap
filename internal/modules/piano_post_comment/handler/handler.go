package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_post_commentv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post_comment/v1"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_comment/v1/piano_post_commentv1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/adapter"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment/usecase"
)

type Handler struct {
	piano_post_commentv1connect.UnimplementedPianoPostCommentServiceHandler
	create *usecase.CreatePianoPostComment
	list   *usecase.ListPianoPostComments
	delete *usecase.DeletePianoPostComment
}

func New(
	create *usecase.CreatePianoPostComment,
	list *usecase.ListPianoPostComments,
	del *usecase.DeletePianoPostComment,
) *Handler {
	return &Handler{create: create, list: list, delete: del}
}

func (h *Handler) CreatePianoPostComment(ctx context.Context, req *connect.Request[piano_post_commentv1.CreatePianoPostCommentRequest]) (*connect.Response[piano_post_commentv1.CreatePianoPostCommentResponse], error) {
	input, err := adapter.FromCreateRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.create.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToCreateResponse(output), nil
}

func (h *Handler) ListPianoPostComments(ctx context.Context, req *connect.Request[piano_post_commentv1.ListPianoPostCommentsRequest]) (*connect.Response[piano_post_commentv1.ListPianoPostCommentsResponse], error) {
	input, err := adapter.FromListRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.list.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListResponse(output), nil
}

func (h *Handler) DeletePianoPostComment(ctx context.Context, req *connect.Request[piano_post_commentv1.DeletePianoPostCommentRequest]) (*connect.Response[piano_post_commentv1.DeletePianoPostCommentResponse], error) {
	input, err := adapter.FromDeleteRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.delete.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToDeleteResponse(output), nil
}
