package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	"github.com/reverie-jp/piamap/internal/domain/entity"
	piano_postv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromCreatePianoPostRequest(ctx context.Context, req *connect.Request[piano_postv1.CreatePianoPostRequest]) (usecase.CreatePianoPostInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.CreatePianoPostInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}

	pianoID, err := resourcename.ParsePiano(req.Msg.GetParent())
	if err != nil {
		return usecase.CreatePianoPostInput{}, err
	}

	post := req.Msg.GetPianoPost()
	if post == nil {
		return usecase.CreatePianoPostInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("piano_post is required"))
	}
	if post.GetVisitTime() == nil {
		return usecase.CreatePianoPostInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("visit_time is required"))
	}

	input := usecase.CreatePianoPostInput{
		RequesterID: requesterID,
		PianoID:     pianoID,
		VisitTime:   post.GetVisitTime().AsTime(),
		Rating:      int16(post.GetRating()),
		Body:        post.Body,
		Visibility:  fromPbVisibility(post.GetVisibility()),
	}
	input.AmbientNoise = optInt16(post.AmbientNoise)
	input.FootTraffic = optInt16(post.FootTraffic)
	input.Resonance = optInt16(post.Resonance)
	input.KeyTouchWeight = optInt16(post.KeyTouchWeight)
	input.TuningQuality = optInt16(post.TuningQuality)
	if input.Visibility == "" {
		input.Visibility = entity.PostVisibilityPublic
	}
	return input, nil
}

func ToCreatePianoPostResponse(output *usecase.CreatePianoPostOutput) *connect.Response[piano_postv1.CreatePianoPostResponse] {
	return connect.NewResponse(&piano_postv1.CreatePianoPostResponse{
		PianoPost: ToPianoPost(output.View),
	})
}

func optInt16(v *int32) *int16 {
	if v == nil {
		return nil
	}
	x := int16(*v)
	return &x
}
