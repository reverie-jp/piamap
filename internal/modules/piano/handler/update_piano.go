package handler

import (
	"context"

	"connectrpc.com/connect"

	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/adapter"
)

func (h *Handler) UpdatePiano(ctx context.Context, req *connect.Request[pianov1.UpdatePianoRequest]) (*connect.Response[pianov1.UpdatePianoResponse], error) {
	input, err := adapter.FromUpdatePianoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.updatePiano.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToUpdatePianoResponse(output), nil
}
