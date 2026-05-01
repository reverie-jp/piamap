package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromUpdatePianoPostRequest(ctx context.Context, req *connect.Request[piano_postv1.UpdatePianoPostRequest]) (usecase.UpdatePianoPostInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UpdatePianoPostInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}

	post := req.Msg.GetPianoPost()
	if post == nil {
		return usecase.UpdatePianoPostInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("piano_post is required"))
	}
	mask := req.Msg.GetUpdateMask()
	if mask == nil {
		return usecase.UpdatePianoPostInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("update_mask is required"))
	}

	_, postID, err := resourcename.ParsePianoPost(post.GetName())
	if err != nil {
		return usecase.UpdatePianoPostInput{}, err
	}

	input := usecase.UpdatePianoPostInput{
		RequesterID: requesterID,
		PostID:      postID,
	}

	for _, p := range mask.Paths {
		switch p {
		case "visit_time":
			if post.GetVisitTime() == nil {
				return usecase.UpdatePianoPostInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("visit_time is required when in mask"))
			}
			input.SetVisitTime = true
			input.VisitTime = post.GetVisitTime().AsTime()
		case "rating":
			input.SetRating = true
			input.Rating = optInt16(post.Rating)
		case "body":
			input.SetBody = true
			input.Body = post.Body
		case "ambient_noise":
			input.SetAmbientNoise = true
			input.AmbientNoise = optInt16(post.AmbientNoise)
		case "foot_traffic":
			input.SetFootTraffic = true
			input.FootTraffic = optInt16(post.FootTraffic)
		case "resonance":
			input.SetResonance = true
			input.Resonance = optInt16(post.Resonance)
		case "key_touch_weight":
			input.SetKeyTouchWeight = true
			input.KeyTouchWeight = optInt16(post.KeyTouchWeight)
		case "tuning_quality":
			input.SetTuningQuality = true
			input.TuningQuality = optInt16(post.TuningQuality)
		case "visibility":
			input.SetVisibility = true
			input.Visibility = fromPbVisibility(post.GetVisibility())
		default:
			return usecase.UpdatePianoPostInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("unknown field in update_mask: " + p))
		}
	}
	return input, nil
}

func ToUpdatePianoPostResponse(output *usecase.UpdatePianoPostOutput) *connect.Response[piano_postv1.UpdatePianoPostResponse] {
	return connect.NewResponse(&piano_postv1.UpdatePianoPostResponse{
		PianoPost: ToPianoPost(output.View),
	})
}
