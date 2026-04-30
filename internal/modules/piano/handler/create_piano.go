package handler

import (
	"context"

	"connectrpc.com/connect"

	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/adapter"
)

func (h *Handler) CreatePiano(ctx context.Context, req *connect.Request[pianov1.CreatePianoRequest]) (*connect.Response[pianov1.CreatePianoResponse], error) {
	input, err := adapter.FromCreatePianoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.createPiano.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToCreatePianoResponse(output), nil
}
