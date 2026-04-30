package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/adapter"
)

func (h *Handler) GetPianoPost(ctx context.Context, req *connect.Request[piano_postv1.GetPianoPostRequest]) (*connect.Response[piano_postv1.GetPianoPostResponse], error) {
	input, err := adapter.FromGetPianoPostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getPianoPost.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetPianoPostResponse(output), nil
}
