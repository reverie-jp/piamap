package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_post_likev1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post_like/v1"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_like/v1/piano_post_likev1connect"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/adapter"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like/usecase"
)

type Handler struct {
	piano_post_likev1connect.UnimplementedPianoPostLikeServiceHandler
	like     *usecase.LikePianoPost
	unlike   *usecase.UnlikePianoPost
	listLiked *usecase.ListLikedPianoPosts
}

func New(
	like *usecase.LikePianoPost,
	unlike *usecase.UnlikePianoPost,
	listLiked *usecase.ListLikedPianoPosts,
) *Handler {
	return &Handler{like: like, unlike: unlike, listLiked: listLiked}
}

func (h *Handler) LikePianoPost(ctx context.Context, req *connect.Request[piano_post_likev1.LikePianoPostRequest]) (*connect.Response[piano_post_likev1.LikePianoPostResponse], error) {
	input, err := adapter.FromLikeRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.like.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToLikeResponse(output), nil
}

func (h *Handler) UnlikePianoPost(ctx context.Context, req *connect.Request[piano_post_likev1.UnlikePianoPostRequest]) (*connect.Response[piano_post_likev1.UnlikePianoPostResponse], error) {
	input, err := adapter.FromUnlikeRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.unlike.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToUnlikeResponse(output), nil
}

func (h *Handler) ListLikedPianoPosts(ctx context.Context, req *connect.Request[piano_post_likev1.ListLikedPianoPostsRequest]) (*connect.Response[piano_post_likev1.ListLikedPianoPostsResponse], error) {
	input, err := adapter.FromListLikedRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listLiked.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListLikedResponse(output), nil
}
