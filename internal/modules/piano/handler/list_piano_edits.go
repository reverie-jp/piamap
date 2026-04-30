package handler

import (
	"context"

	"connectrpc.com/connect"

	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/adapter"
)

func (h *Handler) ListPianoEdits(ctx context.Context, req *connect.Request[pianov1.ListPianoEditsRequest]) (*connect.Response[pianov1.ListPianoEditsResponse], error) {
	input, err := adapter.FromListPianoEditsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listPianoEdits.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListPianoEditsResponse(output), nil
}
