package usecase

import (
	"context"
	"time"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type UpdatePianoPostInput struct {
	RequesterID ulid.ULID
	PostID      ulid.ULID

	SetVisitTime      bool
	VisitTime         time.Time
	SetRating         bool
	Rating            int16
	SetBody           bool
	Body              *string
	SetAmbientNoise   bool
	AmbientNoise      *int16
	SetFootTraffic    bool
	FootTraffic       *int16
	SetResonance      bool
	Resonance         *int16
	SetKeyTouchWeight bool
	KeyTouchWeight    *int16
	SetTuningQuality  bool
	TuningQuality     *int16
	SetVisibility     bool
	Visibility        entity.PostVisibility
}

func (i UpdatePianoPostInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PostID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("post id is required")
	}
	// NOT NULL カラム (rating / visit_time / visibility) は SetX=true のとき値必須。
	if i.SetRating && (i.Rating < 1 || i.Rating > 5) {
		return xerrors.ErrInvalidArgument.WithMessage("rating must be 1..5")
	}
	if i.SetVisitTime {
		if i.VisitTime.IsZero() {
			return xerrors.ErrInvalidArgument.WithMessage("visit_time cannot be cleared")
		}
		if i.VisitTime.After(time.Now().Add(5 * time.Minute)) {
			return xerrors.ErrInvalidArgument.WithMessage("visit_time must not be in the future")
		}
	}
	if i.SetBody && i.Body != nil && len(*i.Body) > 4000 {
		return xerrors.ErrInvalidArgument.WithMessage("body too long")
	}
	type attr struct {
		set bool
		v   *int16
	}
	for _, a := range []attr{
		{i.SetAmbientNoise, i.AmbientNoise},
		{i.SetFootTraffic, i.FootTraffic},
		{i.SetResonance, i.Resonance},
		{i.SetKeyTouchWeight, i.KeyTouchWeight},
		{i.SetTuningQuality, i.TuningQuality},
	} {
		if a.set && a.v != nil && (*a.v < 1 || *a.v > 5) {
			return xerrors.ErrInvalidArgument.WithMessage("attribute must be 1..5")
		}
	}
	if i.SetVisibility {
		switch i.Visibility {
		case entity.PostVisibilityPublic, entity.PostVisibilityPrivate:
		default:
			return xerrors.ErrInvalidArgument.WithMessage("invalid visibility")
		}
	}
	return nil
}

func (i UpdatePianoPostInput) HasAnyChange() bool {
	return i.SetVisitTime || i.SetRating || i.SetBody ||
		i.SetAmbientNoise || i.SetFootTraffic || i.SetResonance ||
		i.SetKeyTouchWeight || i.SetTuningQuality || i.SetVisibility
}

type UpdatePianoPostOutput struct {
	View *gateway.PianoPostView
}

type UpdatePianoPost struct {
	gw gateway.Gateway
}

func NewUpdatePianoPost(gw gateway.Gateway) *UpdatePianoPost {
	return &UpdatePianoPost{gw: gw}
}

func (uc *UpdatePianoPost) Execute(ctx context.Context, input UpdatePianoPostInput) (*UpdatePianoPostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if !input.HasAnyChange() {
		return nil, xerrors.ErrInvalidArgument.WithMessage("no fields to update")
	}

	existing, err := uc.gw.GetPianoPost(ctx, input.PostID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, xerrors.ErrNotFound
	}
	if existing.UserID != input.RequesterID {
		return nil, xerrors.ErrPermissionDenied.WithMessage("not the author")
	}

	params := gateway.UpdatePianoPostParams{
		ID:                input.PostID,
		SetVisitTime:      input.SetVisitTime,
		SetRating:         input.SetRating,
		SetBody:           input.SetBody,
		SetAmbientNoise:   input.SetAmbientNoise,
		SetFootTraffic:    input.SetFootTraffic,
		SetResonance:      input.SetResonance,
		SetKeyTouchWeight: input.SetKeyTouchWeight,
		SetTuningQuality:  input.SetTuningQuality,
		SetVisibility:     input.SetVisibility,
	}
	if input.SetVisitTime {
		v := input.VisitTime
		params.VisitTime = &v
	}
	if input.SetRating {
		v := input.Rating
		params.Rating = &v
	}
	if input.SetBody {
		params.Body = input.Body
	}
	if input.SetAmbientNoise {
		params.AmbientNoise = input.AmbientNoise
	}
	if input.SetFootTraffic {
		params.FootTraffic = input.FootTraffic
	}
	if input.SetResonance {
		params.Resonance = input.Resonance
	}
	if input.SetKeyTouchWeight {
		params.KeyTouchWeight = input.KeyTouchWeight
	}
	if input.SetTuningQuality {
		params.TuningQuality = input.TuningQuality
	}
	if input.SetVisibility {
		v := input.Visibility
		params.Visibility = &v
	}

	if err := uc.gw.UpdatePianoPost(ctx, params); err != nil {
		return nil, err
	}

	post, err := uc.gw.GetPianoPost(ctx, input.PostID)
	if err != nil {
		return nil, err
	}
	view, err := uc.gw.BuildPianoPostView(ctx, input.RequesterID, post)
	if err != nil {
		return nil, err
	}
	return &UpdatePianoPostOutput{View: view}, nil
}
