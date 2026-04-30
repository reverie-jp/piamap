package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

const (
	ListPianoEditsDefaultLimit = 20
	ListPianoEditsMaxLimit     = 100
)

type ListPianoEditsInput struct {
	PianoID  ulid.ULID
	PageSize int32
	AfterID  *ulid.ULID
}

func (i ListPianoEditsInput) Validate() error {
	if i.PianoID.IsZero() {
		return xerrors.ErrInvalidArgument.WithMessage("piano id is required")
	}
	return nil
}

type ListPianoEditsOutput struct {
	Views  []*gateway.PianoEditView
	NextID *ulid.ULID
}

type ListPianoEdits struct {
	pianoGateway gateway.Gateway
}

func NewListPianoEdits(pianoGateway gateway.Gateway) *ListPianoEdits {
	return &ListPianoEdits{pianoGateway: pianoGateway}
}

func (uc *ListPianoEdits) Execute(ctx context.Context, input ListPianoEditsInput) (*ListPianoEditsOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	limit := input.PageSize
	if limit <= 0 {
		limit = ListPianoEditsDefaultLimit
	} else if limit > ListPianoEditsMaxLimit {
		limit = ListPianoEditsMaxLimit
	}
	queryLimit := limit + 1

	edits, err := uc.pianoGateway.ListPianoEditsByPiano(ctx, gateway.ListPianoEditsParams{
		PianoID: input.PianoID,
		AfterID: input.AfterID,
		Limit:   queryLimit,
	})
	if err != nil {
		return nil, err
	}
	next, edits := splitNextEdit(edits, int(limit))
	views, err := uc.pianoGateway.BuildListPianoEditViews(ctx, edits)
	if err != nil {
		return nil, err
	}
	return &ListPianoEditsOutput{Views: views, NextID: next}, nil
}

func splitNextEdit(items []*entity.PianoEdit, limit int) (*ulid.ULID, []*entity.PianoEdit) {
	if len(items) <= limit {
		return nil, items
	}
	cut := items[:limit]
	id := cut[len(cut)-1].ID
	return &id, cut
}
