package handler

import (
	"context"

	"connectrpc.com/connect"

	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/adapter"
)

func (h *Handler) SearchPianos(ctx context.Context, req *connect.Request[pianov1.SearchPianosRequest]) (*connect.Response[pianov1.SearchPianosResponse], error) {
	input, err := adapter.FromSearchPianosRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.searchPianos.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToSearchPianosResponse(output), nil
}
