package usecase

import (
	"context"
	"time"

	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type CreatePianoPostInput struct {
	RequesterID    ulid.ULID
	PianoID        ulid.ULID
	VisitTime      time.Time
	Rating         *int16
	Body           *string
	AmbientNoise   *int16
	FootTraffic    *int16
	Resonance      *int16
	KeyTouchWeight *int16
	TuningQuality  *int16
	Visibility     entity.PostVisibility
}

func (i CreatePianoPostInput) Validate() error {
	if i.RequesterID.IsZero() {
		return xerrors.ErrUnauthenticated
	}
	if i.PianoID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("piano id is required")
	}
	if i.Rating != nil && (*i.Rating < 1 || *i.Rating > 5) {
		return xerrors.ErrInvalidArgument.WithMessage("rating must be 1..5")
	}
	// rating または body のどちらかは必須。
	hasBody := i.Body != nil && len(*i.Body) > 0
	if i.Rating == nil && !hasBody {
		return xerrors.ErrInvalidArgument.WithMessage("rating or body is required")
	}
	if i.VisitTime.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("visit_time is required")
	}
	if i.VisitTime.After(time.Now().Add(5 * time.Minute)) {
		return xerrors.ErrInvalidArgument.WithMessage("visit_time must not be in the future")
	}
	if i.Body != nil && len(*i.Body) > 4000 {
		return xerrors.ErrInvalidArgument.WithMessage("body too long")
	}
	for _, v := range []*int16{i.AmbientNoise, i.FootTraffic, i.Resonance, i.KeyTouchWeight, i.TuningQuality} {
		if v != nil && (*v < 1 || *v > 5) {
			return xerrors.ErrInvalidArgument.WithMessage("attribute must be 1..5")
		}
	}
	switch i.Visibility {
	case entity.PostVisibilityPublic, entity.PostVisibilityPrivate:
	default:
		return xerrors.ErrInvalidArgument.WithMessage("invalid visibility")
	}
	return nil
}

type CreatePianoPostOutput struct {
	View *gateway.PianoPostView
}

type CreatePianoPost struct {
	gw           gateway.Gateway
	pianoGateway pianogw.Gateway
	userGateway  usergw.Gateway
	tx           transaction.Runner
}

func NewCreatePianoPost(
	gw gateway.Gateway,
	pianoGateway pianogw.Gateway,
	userGateway usergw.Gateway,
	tx transaction.Runner,
) *CreatePianoPost {
	return &CreatePianoPost{gw: gw, pianoGateway: pianoGateway, userGateway: userGateway, tx: tx}
}

func (uc *CreatePianoPost) Execute(ctx context.Context, input CreatePianoPostInput) (*CreatePianoPostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	piano, err := uc.pianoGateway.GetPiano(ctx, input.PianoID)
	if err != nil {
		return nil, err
	}
	if piano == nil || piano.Status == entity.PianoStatusRemoved {
		return nil, xerrors.ErrPianoNotFound
	}

	postID := ulid.New()
	err = uc.tx.WithTx(ctx, func(q sqlc.Querier) error {
		// Tx 内では BuildListPianoPostViews を呼ばないので likeLookup=nil で十分。
		txGw := gateway.New(q, uc.userGateway, uc.pianoGateway, nil)
		if err := txGw.CreatePianoPost(ctx, gateway.CreatePianoPostParams{
			ID:             postID,
			UserID:         input.RequesterID,
			PianoID:        input.PianoID,
			VisitTime:      input.VisitTime,
			Rating:         input.Rating,
			Body:           input.Body,
			AmbientNoise:   input.AmbientNoise,
			FootTraffic:    input.FootTraffic,
			Resonance:      input.Resonance,
			KeyTouchWeight: input.KeyTouchWeight,
			TuningQuality:  input.TuningQuality,
			Visibility:     input.Visibility,
		}); err != nil {
			return err
		}
		// 投稿成立 = 訪問成立。visited リストに UPSERT (冪等)。
		return txGw.UpsertPianoUserListVisited(ctx, input.RequesterID, input.PianoID)
	})
	if err != nil {
		return nil, err
	}

	post, err := uc.gw.GetPianoPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	view, err := uc.gw.BuildPianoPostView(ctx, input.RequesterID, post)
	if err != nil {
		return nil, err
	}
	return &CreatePianoPostOutput{View: view}, nil
}
