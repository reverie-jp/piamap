package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/adapter"
)

func (h *Handler) CreatePianoPost(ctx context.Context, req *connect.Request[piano_postv1.CreatePianoPostRequest]) (*connect.Response[piano_postv1.CreatePianoPostResponse], error) {
	input, err := adapter.FromCreatePianoPostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.createPianoPost.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToCreatePianoPostResponse(output), nil
}
