package handler

import (
	"context"

	"connectrpc.com/connect"

	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/adapter"
)

func (h *Handler) GetPiano(ctx context.Context, req *connect.Request[pianov1.GetPianoRequest]) (*connect.Response[pianov1.GetPianoResponse], error) {
	input, err := adapter.FromGetPianoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getPiano.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetPianoResponse(output), nil
}
