package usecase

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type UpdatePiano struct {
	pianoGateway gateway.Gateway
	userGateway  usergw.Gateway
	tx           transaction.Runner
}

func NewUpdatePiano(pianoGateway gateway.Gateway, userGateway usergw.Gateway, tx transaction.Runner) *UpdatePiano {
	return &UpdatePiano{
		pianoGateway: pianoGateway,
		userGateway:  userGateway,
		tx:           tx,
	}
}

func (uc *UpdatePiano) Execute(ctx context.Context, input UpdatePianoInput) (*UpdatePianoOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if !input.HasAnyChange() {
		return nil, xerrors.ErrInvalidArgument.WithMessage("no fields to update")
	}

	existing, err := uc.pianoGateway.GetPiano(ctx, input.PianoID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, xerrors.ErrPianoNotFound
	}

	// 信頼ライン: status='removed' への変更 / 500m 以上の座標移動は trusted user のみ
	destructive := false
	if input.SetStatus && input.Status == entity.PianoStatusRemoved && existing.Status != entity.PianoStatusRemoved {
		destructive = true
	}
	if input.SetLocation && haversineDistanceM(existing.Location, input.Location) > TrustedRequiredMoveDistanceM {
		destructive = true
	}
	if destructive {
		editor, err := uc.userGateway.GetUserByID(ctx, input.RequesterID)
		if err != nil {
			return nil, err
		}
		if editor == nil || !editor.IsTrusted(time.Now()) {
			return nil, xerrors.ErrPermissionDenied.WithMessage("destructive edit requires trusted user")
		}
	}

	operation := chooseOperation(existing, input)
	changesJSON := buildChangesJSON(existing, input)

	err = uc.tx.WithTx(ctx, func(q sqlc.Querier) error {
		txPianoGw := gateway.New(q, uc.userGateway)
		if err := txPianoGw.UpdatePiano(ctx, gateway.UpdatePianoParams{
			ID:                  input.PianoID,
			SetName:             input.SetName,
			Name:                optStr(input.SetName, input.Name),
			SetDescription:      input.SetDescription,
			Description:         optPtrStr(input.SetDescription, input.Description),
			SetAddress:          input.SetAddress,
			Address:             optPtrStr(input.SetAddress, input.Address),
			SetPrefecture:       input.SetPrefecture,
			Prefecture:          optPtrStr(input.SetPrefecture, input.Prefecture),
			SetCity:             input.SetCity,
			City:                optPtrStr(input.SetCity, input.City),
			SetKind:             input.SetKind,
			Kind:                optKind(input.SetKind, input.Kind),
			SetVenueType:        input.SetVenueType,
			VenueType:           optPtrStr(input.SetVenueType, input.VenueType),
			SetPianoType:        input.SetPianoType,
			PianoType:           optPianoType(input.SetPianoType, input.PianoType),
			SetPianoBrand:       input.SetPianoBrand,
			PianoBrand:          optStr(input.SetPianoBrand, input.PianoBrand),
			SetPianoModel:       input.SetPianoModel,
			PianoModel:          optPtrStr(input.SetPianoModel, input.PianoModel),
			SetManufactureYear:  input.SetManufactureYear,
			ManufactureYear:     optPtrInt16(input.SetManufactureYear, input.ManufactureYear),
			SetHours:            input.SetHours,
			Hours:               optPtrStr(input.SetHours, input.Hours),
			SetStatus:           input.SetStatus,
			Status:              optStatus(input.SetStatus, input.Status),
			SetAvailability:     input.SetAvailability,
			Availability:        optAvailability(input.SetAvailability, input.Availability),
			SetAvailabilityNote: input.SetAvailabilityNote,
			AvailabilityNote:    optPtrStr(input.SetAvailabilityNote, input.AvailabilityNote),
			SetInstallTime:      input.SetInstallTime,
			InstallTime:         optPtrTime(input.SetInstallTime, input.InstallTime),
			SetRemoveTime:       input.SetRemoveTime,
			RemoveTime:          optPtrTime(input.SetRemoveTime, input.RemoveTime),
		}); err != nil {
			return err
		}
		if input.SetLocation {
			if err := txPianoGw.UpdatePianoLocation(ctx, input.PianoID, input.Location); err != nil {
				return err
			}
		}
		editorID := input.RequesterID
		return txPianoGw.CreatePianoEdit(ctx, gateway.CreatePianoEditParams{
			PianoID:      input.PianoID,
			EditorUserID: &editorID,
			Operation:    operation,
			Changes:      changesJSON,
			Summary:      input.EditSummary,
		})
	})
	if err != nil {
		return nil, err
	}

	piano, err := uc.pianoGateway.GetPiano(ctx, input.PianoID)
	if err != nil {
		return nil, err
	}
	view, err := uc.pianoGateway.BuildPianoView(ctx, input.RequesterID, piano)
	if err != nil {
		return nil, err
	}
	return &UpdatePianoOutput{View: view}, nil
}

func chooseOperation(existing *entity.Piano, input UpdatePianoInput) entity.PianoEditOperation {
	if input.SetStatus && input.Status != existing.Status {
		return entity.PianoEditOpStatusChange
	}
	if input.SetKind && input.Kind != existing.Kind {
		return entity.PianoEditOpKindChange
	}
	return entity.PianoEditOpUpdate
}

func buildChangesJSON(existing *entity.Piano, input UpdatePianoInput) []byte {
	changes := map[string]any{}
	diffStr := func(key, oldV, newV string) {
		if oldV != newV {
			changes[key] = map[string]string{"old": oldV, "new": newV}
		}
	}
	diffOptStr := func(key string, oldV *string, newV *string) {
		o := derefStr(oldV)
		n := derefStr(newV)
		if o != n {
			changes[key] = map[string]string{"old": o, "new": n}
		}
	}
	diffOptInt16 := func(key string, oldV *int16, newV *int16) {
		o := derefInt16(oldV)
		n := derefInt16(newV)
		if o != n {
			changes[key] = map[string]any{"old": o, "new": n}
		}
	}
	diffOptTime := func(key string, oldV *time.Time, newV *time.Time) {
		o := formatOptTime(oldV)
		n := formatOptTime(newV)
		if o != n {
			changes[key] = map[string]string{"old": o, "new": n}
		}
	}

	if input.SetName {
		diffStr("name", existing.Name, input.Name)
	}
	if input.SetDescription {
		diffOptStr("description", existing.Description, input.Description)
	}
	if input.SetAddress {
		diffOptStr("address", existing.Address, input.Address)
	}
	if input.SetPrefecture {
		diffOptStr("prefecture", existing.Prefecture, input.Prefecture)
	}
	if input.SetCity {
		diffOptStr("city", existing.City, input.City)
	}
	if input.SetVenueType {
		diffOptStr("venue_type", existing.VenueType, input.VenueType)
	}
	if input.SetPianoBrand {
		diffStr("piano_brand", existing.PianoBrand, input.PianoBrand)
	}
	if input.SetPianoModel {
		diffOptStr("piano_model", existing.PianoModel, input.PianoModel)
	}
	if input.SetManufactureYear {
		diffOptInt16("manufacture_year", existing.ManufactureYear, input.ManufactureYear)
	}
	if input.SetHours {
		diffOptStr("hours", existing.Hours, input.Hours)
	}
	if input.SetPianoType && input.PianoType != existing.PianoType {
		changes["piano_type"] = map[string]string{"old": string(existing.PianoType), "new": string(input.PianoType)}
	}
	if input.SetAvailability && input.Availability != existing.Availability {
		changes["availability"] = map[string]string{"old": string(existing.Availability), "new": string(input.Availability)}
	}
	if input.SetAvailabilityNote {
		diffOptStr("availability_note", existing.AvailabilityNote, input.AvailabilityNote)
	}
	if input.SetInstallTime {
		diffOptTime("install_time", existing.InstallTime, input.InstallTime)
	}
	if input.SetRemoveTime {
		diffOptTime("remove_time", existing.RemoveTime, input.RemoveTime)
	}
	if input.SetLocation && (input.Location.Latitude != existing.Location.Latitude || input.Location.Longitude != existing.Location.Longitude) {
		changes["location"] = map[string]any{
			"old": map[string]float64{"lat": existing.Location.Latitude, "lng": existing.Location.Longitude},
			"new": map[string]float64{"lat": input.Location.Latitude, "lng": input.Location.Longitude},
		}
	}
	if input.SetKind && input.Kind != existing.Kind {
		changes["kind"] = map[string]string{"old": string(existing.Kind), "new": string(input.Kind)}
	}
	if input.SetStatus && input.Status != existing.Status {
		changes["status"] = map[string]string{"old": string(existing.Status), "new": string(input.Status)}
	}
	if len(changes) == 0 {
		return nil
	}
	b, _ := json.Marshal(changes)
	return b
}

func derefStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func derefInt16(v *int16) any {
	if v == nil {
		return nil
	}
	return *v
}

func formatOptTime(v *time.Time) string {
	if v == nil {
		return ""
	}
	return v.UTC().Format(time.RFC3339)
}

// haversineDistanceM は WGS84 ふたつの点の概算距離 (m)。
// PostGIS の ST_Distance(geography) と同じ意図だが Go 側で済ませる (DB 往復を避けるため)。
func haversineDistanceM(a, b entity.LatLng) float64 {
	const earthRadiusM = 6371_000.0
	rad := math.Pi / 180
	lat1 := a.Latitude * rad
	lat2 := b.Latitude * rad
	dLat := (b.Latitude - a.Latitude) * rad
	dLng := (b.Longitude - a.Longitude) * rad
	h := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * earthRadiusM * math.Asin(math.Min(1, math.Sqrt(h)))
}

func optStr(set bool, v string) *string {
	if !set {
		return nil
	}
	s := v
	return &s
}

func optPtrStr(set bool, v *string) *string {
	if !set {
		return nil
	}
	return v
}

func optPtrInt16(set bool, v *int16) *int16 {
	if !set {
		return nil
	}
	return v
}

func optPtrTime(set bool, v *time.Time) *time.Time {
	if !set {
		return nil
	}
	return v
}

func optKind(set bool, v entity.PianoKind) *entity.PianoKind {
	if !set {
		return nil
	}
	k := v
	return &k
}

func optPianoType(set bool, v entity.PianoType) *entity.PianoType {
	if !set {
		return nil
	}
	t := v
	return &t
}

func optStatus(set bool, v entity.PianoStatus) *entity.PianoStatus {
	if !set {
		return nil
	}
	s := v
	return &s
}

func optAvailability(set bool, v entity.PianoAvailability) *entity.PianoAvailability {
	if !set {
		return nil
	}
	a := v
	return &a
}
