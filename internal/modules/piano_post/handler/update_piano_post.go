package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/adapter"
)

func (h *Handler) UpdatePianoPost(ctx context.Context, req *connect.Request[piano_postv1.UpdatePianoPostRequest]) (*connect.Response[piano_postv1.UpdatePianoPostResponse], error) {
	input, err := adapter.FromUpdatePianoPostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.updatePianoPost.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToUpdatePianoPostResponse(output), nil
}
