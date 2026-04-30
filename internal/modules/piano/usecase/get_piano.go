package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type GetPiano struct {
	pianoGateway gateway.Gateway
}

func NewGetPiano(pianoGateway gateway.Gateway) *GetPiano {
	return &GetPiano{pianoGateway: pianoGateway}
}

func (uc *GetPiano) Execute(ctx context.Context, input GetPianoInput) (*GetPianoOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	piano, err := uc.pianoGateway.GetPiano(ctx, input.PianoID)
	if err != nil {
		return nil, err
	}
	if piano == nil {
		return nil, xerrors.ErrPianoNotFound
	}
	// removed / pending は creator 本人以外には見せない (= 404 扱い)。
	if piano.Status == entity.PianoStatusRemoved || piano.Status == entity.PianoStatusPending {
		if piano.CreatorUserID == nil || input.RequesterID.IsZero() || *piano.CreatorUserID != input.RequesterID {
			return nil, xerrors.ErrPianoHidden
		}
	}
	view, err := uc.pianoGateway.BuildPianoView(ctx, input.RequesterID, piano)
	if err != nil {
		return nil, err
	}
	return &GetPianoOutput{View: view}, nil
}
