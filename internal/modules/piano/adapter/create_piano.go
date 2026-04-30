package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	"github.com/reverie-jp/piamap/internal/domain/entity"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromCreatePianoRequest(ctx context.Context, req *connect.Request[pianov1.CreatePianoRequest]) (usecase.CreatePianoInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.CreatePianoInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	piano := req.Msg.GetPiano()
	if piano == nil {
		return usecase.CreatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("piano is required"))
	}
	loc := piano.GetLocation()
	if loc == nil {
		return usecase.CreatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("location is required"))
	}
	kind, ok := fromPbKind(piano.GetKind())
	if !ok {
		kind = entity.PianoKindStreet // MVP デフォルト
	}
	pianoType, ok := fromPbPianoType(piano.GetPianoType())
	if !ok {
		pianoType = entity.PianoTypeUnknown
	}
	availability, ok := fromPbAvailability(piano.GetAvailability())
	if !ok {
		availability = entity.PianoAvailabilityRegular
	}

	in := usecase.CreatePianoInput{
		RequesterID:      requesterID,
		Name:             piano.GetDisplayName(),
		Description:      piano.Description,
		Location:         entity.LatLng{Latitude: loc.GetLatitude(), Longitude: loc.GetLongitude()},
		Address:          piano.Address,
		Prefecture:       piano.Prefecture,
		City:             piano.City,
		Kind:             kind,
		VenueType:        piano.VenueType,
		PianoType:        pianoType,
		PianoBrand:       piano.GetPianoBrand(),
		PianoModel:       piano.PianoModel,
		Hours:            piano.Hours,
		Availability:     availability,
		AvailabilityNote: piano.AvailabilityNote,
	}
	if piano.PianoBrand == "" {
		in.PianoBrand = "unknown"
	}
	if piano.ManufactureYear != nil {
		v := int16(*piano.ManufactureYear)
		in.ManufactureYear = &v
	}
	return in, nil
}

func ToCreatePianoResponse(output *usecase.CreatePianoOutput) *connect.Response[pianov1.CreatePianoResponse] {
	return connect.NewResponse(&pianov1.CreatePianoResponse{Piano: ToPiano(output.View)})
}
