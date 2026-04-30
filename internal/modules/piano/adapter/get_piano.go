package adapter

import (
	"context"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
)

func FromGetPianoRequest(ctx context.Context, req *connect.Request[pianov1.GetPianoRequest]) (usecase.GetPianoInput, error) {
	id, err := resourcename.ParsePiano(req.Msg.Name)
	if err != nil {
		return usecase.GetPianoInput{}, err
	}
	input := usecase.GetPianoInput{PianoID: id}
	if requesterID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = requesterID
	}
	return input, nil
}

func ToGetPianoResponse(output *usecase.GetPianoOutput) *connect.Response[pianov1.GetPianoResponse] {
	return connect.NewResponse(&pianov1.GetPianoResponse{Piano: ToPiano(output.View)})
}
