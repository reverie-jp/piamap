package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/adapter"
)

func (h *Handler) ListPianoPosts(ctx context.Context, req *connect.Request[piano_postv1.ListPianoPostsRequest]) (*connect.Response[piano_postv1.ListPianoPostsResponse], error) {
	input, err := adapter.FromListPianoPostsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listPianoPosts.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListPianoPostsResponse(output), nil
}
