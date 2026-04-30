package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	"github.com/reverie-jp/piamap/internal/domain/entity"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromUpdatePianoRequest(ctx context.Context, req *connect.Request[pianov1.UpdatePianoRequest]) (usecase.UpdatePianoInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UpdatePianoInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	piano := req.Msg.GetPiano()
	if piano == nil {
		return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("piano is required"))
	}
	pianoID, err := resourcename.ParsePiano(piano.GetName())
	if err != nil {
		return usecase.UpdatePianoInput{}, err
	}
	mask := req.Msg.GetUpdateMask()
	if mask == nil {
		return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("update_mask is required"))
	}

	in := usecase.UpdatePianoInput{
		RequesterID: requesterID,
		PianoID:     pianoID,
		EditSummary: req.Msg.EditSummary,
	}

	// クライアントは FieldMask paths を camelCase で送るが、protojson は受信時に
	// proto field name (snake_case) へ変換するため、switch は snake_case で行う。
	for _, p := range mask.Paths {
		switch p {
		case "display_name":
			in.SetName = true
			in.Name = piano.GetDisplayName()
		case "description":
			in.SetDescription = true
			in.Description = piano.Description
		case "location":
			loc := piano.GetLocation()
			if loc == nil {
				return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("location is required when in update_mask"))
			}
			in.SetLocation = true
			in.Location = entity.LatLng{Latitude: loc.GetLatitude(), Longitude: loc.GetLongitude()}
		case "address":
			in.SetAddress = true
			in.Address = piano.Address
		case "prefecture":
			in.SetPrefecture = true
			in.Prefecture = piano.Prefecture
		case "city":
			in.SetCity = true
			in.City = piano.City
		case "kind":
			k, ok := fromPbKind(piano.GetKind())
			if !ok {
				return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid kind"))
			}
			in.SetKind = true
			in.Kind = k
		case "venue_type":
			in.SetVenueType = true
			in.VenueType = piano.VenueType
		case "piano_type":
			t, ok := fromPbPianoType(piano.GetPianoType())
			if !ok {
				return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid piano_type"))
			}
			in.SetPianoType = true
			in.PianoType = t
		case "piano_brand":
			in.SetPianoBrand = true
			in.PianoBrand = piano.GetPianoBrand()
		case "piano_model":
			in.SetPianoModel = true
			in.PianoModel = piano.PianoModel
		case "manufacture_year":
			in.SetManufactureYear = true
			if piano.ManufactureYear != nil {
				v := int16(*piano.ManufactureYear)
				in.ManufactureYear = &v
			}
		case "hours":
			in.SetHours = true
			in.Hours = piano.Hours
		case "status":
			s, ok := fromPbStatus(piano.GetStatus())
			if !ok {
				return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid status"))
			}
			in.SetStatus = true
			in.Status = s
		case "availability":
			a, ok := fromPbAvailability(piano.GetAvailability())
			if !ok {
				return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid availability"))
			}
			in.SetAvailability = true
			in.Availability = a
		case "availability_note":
			in.SetAvailabilityNote = true
			in.AvailabilityNote = piano.AvailabilityNote
		case "install_time":
			in.SetInstallTime = true
			if piano.InstallTime != nil {
				t := piano.InstallTime.AsTime()
				in.InstallTime = &t
			}
		case "remove_time":
			in.SetRemoveTime = true
			if piano.RemoveTime != nil {
				t := piano.RemoveTime.AsTime()
				in.RemoveTime = &t
			}
		default:
			return usecase.UpdatePianoInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("unknown field in update_mask: " + p))
		}
	}
	return in, nil
}

func ToUpdatePianoResponse(output *usecase.UpdatePianoOutput) *connect.Response[pianov1.UpdatePianoResponse] {
	return connect.NewResponse(&pianov1.UpdatePianoResponse{Piano: ToPiano(output.View)})
}
