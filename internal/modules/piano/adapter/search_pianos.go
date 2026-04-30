package adapter

import (
	"context"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
)

func FromSearchPianosRequest(ctx context.Context, req *connect.Request[pianov1.SearchPianosRequest]) (usecase.SearchPianosInput, error) {
	in := usecase.SearchPianosInput{}
	if requesterID, ok := interceptor.UserIDFromContext(ctx); ok {
		in.RequesterID = requesterID
	}
	if b := req.Msg.GetBounds(); b != nil && b.GetSouthwest() != nil && b.GetNortheast() != nil {
		in.HasBounds = true
		in.MinLat = b.GetSouthwest().GetLatitude()
		in.MinLng = b.GetSouthwest().GetLongitude()
		in.MaxLat = b.GetNortheast().GetLatitude()
		in.MaxLng = b.GetNortheast().GetLongitude()
	} else if c := req.Msg.GetCenter(); c != nil && req.Msg.RadiusM != nil {
		in.HasCenter = true
		in.CenterLat = c.GetLatitude()
		in.CenterLng = c.GetLongitude()
		in.RadiusM = *req.Msg.RadiusM
	}
	if req.Msg.Limit != nil {
		in.Limit = *req.Msg.Limit
	}
	if req.Msg.Kind != nil {
		if k, ok := fromPbKind(*req.Msg.Kind); ok {
			in.Kind = &k
		}
	}
	if req.Msg.PianoType != nil {
		if t, ok := fromPbPianoType(*req.Msg.PianoType); ok {
			in.PianoType = &t
		}
	}
	if req.Msg.MinRatingAverage != nil {
		v := *req.Msg.MinRatingAverage
		in.MinRatingAverage = &v
	}
	if req.Msg.Query != nil {
		in.Query = *req.Msg.Query
	}
	if req.Msg.PianoBrand != nil {
		v := *req.Msg.PianoBrand
		in.PianoBrand = &v
	}
	if req.Msg.MinAmbientNoiseAverage != nil {
		v := *req.Msg.MinAmbientNoiseAverage
		in.MinAmbientNoiseAverage = &v
	}
	if req.Msg.MinFootTrafficAverage != nil {
		v := *req.Msg.MinFootTrafficAverage
		in.MinFootTrafficAverage = &v
	}
	if req.Msg.MinResonanceAverage != nil {
		v := *req.Msg.MinResonanceAverage
		in.MinResonanceAverage = &v
	}
	if req.Msg.MinKeyTouchWeightAverage != nil {
		v := *req.Msg.MinKeyTouchWeightAverage
		in.MinKeyTouchWeightAverage = &v
	}
	if req.Msg.MinTuningQualityAverage != nil {
		v := *req.Msg.MinTuningQualityAverage
		in.MinTuningQualityAverage = &v
	}
	return in, nil
}

func ToSearchPianosResponse(output *usecase.SearchPianosOutput) *connect.Response[pianov1.SearchPianosResponse] {
	pianos := make([]*pianov1.Piano, 0, len(output.Views))
	for _, v := range output.Views {
		if pb := ToPiano(v); pb != nil {
			pianos = append(pianos, pb)
		}
	}
	return connect.NewResponse(&pianov1.SearchPianosResponse{Pianos: pianos})
}
