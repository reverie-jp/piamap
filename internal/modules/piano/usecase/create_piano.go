package usecase

import (
	"context"
	"encoding/json"

	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type CreatePiano struct {
	pianoGateway gateway.Gateway
	userGateway  usergw.Gateway
	tx           transaction.Runner
}

func NewCreatePiano(pianoGateway gateway.Gateway, userGateway usergw.Gateway, tx transaction.Runner) *CreatePiano {
	return &CreatePiano{
		pianoGateway: pianoGateway,
		userGateway:  userGateway,
		tx:           tx,
	}
}

func (uc *CreatePiano) Execute(ctx context.Context, input CreatePianoInput) (*CreatePianoOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	pianoID := ulid.New()
	creatorID := input.RequesterID

	changesJSON, _ := json.Marshal(map[string]any{
		"name":      input.Name,
		"latitude":  input.Location.Latitude,
		"longitude": input.Location.Longitude,
		"kind":      string(input.Kind),
	})

	err := uc.tx.WithTx(ctx, func(q sqlc.Querier) error {
		txPianoGw := gateway.New(q, uc.userGateway)
		if err := txPianoGw.CreatePiano(ctx, gateway.CreatePianoParams{
			ID:               pianoID,
			Name:             input.Name,
			Description:      input.Description,
			Location:         input.Location,
			Address:          input.Address,
			Prefecture:       input.Prefecture,
			City:             input.City,
			Kind:             input.Kind,
			VenueType:        input.VenueType,
			PianoType:        input.PianoType,
			PianoBrand:       input.PianoBrand,
			PianoModel:       input.PianoModel,
			ManufactureYear:  input.ManufactureYear,
			Hours:            input.Hours,
			Availability:     input.Availability,
			AvailabilityNote: input.AvailabilityNote,
			CreatorUserID:    &creatorID,
		}); err != nil {
			return err
		}
		return txPianoGw.CreatePianoEdit(ctx, gateway.CreatePianoEditParams{
			PianoID:      pianoID,
			EditorUserID: &creatorID,
			Operation:    entity.PianoEditOpCreate,
			Changes:      changesJSON,
		})
	})
	if err != nil {
		return nil, err
	}

	piano, err := uc.pianoGateway.GetPiano(ctx, pianoID)
	if err != nil {
		return nil, err
	}
	view, err := uc.pianoGateway.BuildPianoView(ctx, input.RequesterID, piano)
	if err != nil {
		return nil, err
	}
	return &CreatePianoOutput{View: view}, nil
}
