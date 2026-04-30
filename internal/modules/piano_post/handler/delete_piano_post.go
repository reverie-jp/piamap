package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/adapter"
)

func (h *Handler) DeletePianoPost(ctx context.Context, req *connect.Request[piano_postv1.DeletePianoPostRequest]) (*connect.Response[piano_postv1.DeletePianoPostResponse], error) {
	input, err := adapter.FromDeletePianoPostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if _, err := h.deletePianoPost.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToDeletePianoPostResponse(), nil
}
