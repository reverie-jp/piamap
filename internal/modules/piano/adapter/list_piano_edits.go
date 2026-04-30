package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	pianov1 "github.com/reverie-jp/piamap/internal/gen/pb/piano/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano/usecase"
	"github.com/reverie-jp/piamap/internal/platform/resourcename"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromListPianoEditsRequest(_ context.Context, req *connect.Request[pianov1.ListPianoEditsRequest]) (usecase.ListPianoEditsInput, error) {
	parent := req.Msg.GetParent()
	if parent == "" {
		return usecase.ListPianoEditsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("parent is required"))
	}
	pianoID, err := resourcename.ParsePiano(parent)
	if err != nil {
		return usecase.ListPianoEditsInput{}, err
	}
	in := usecase.ListPianoEditsInput{PianoID: pianoID}
	if req.Msg.PageSize != nil {
		in.PageSize = *req.Msg.PageSize
	}
	if t := req.Msg.GetPageToken(); t != "" {
		id, err := ulid.Parse(t)
		if err != nil {
			return usecase.ListPianoEditsInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("invalid page_token"))
		}
		in.AfterID = &id
	}
	return in, nil
}

func ToListPianoEditsResponse(output *usecase.ListPianoEditsOutput) *connect.Response[pianov1.ListPianoEditsResponse] {
	edits := make([]*pianov1.PianoEdit, 0, len(output.Views))
	for _, v := range output.Views {
		if v == nil || v.Edit == nil {
			continue
		}
		edits = append(edits, toPbPianoEdit(v))
	}
	resp := &pianov1.ListPianoEditsResponse{Edits: edits}
	if output.NextID != nil {
		resp.NextPageToken = output.NextID.String()
	}
	return connect.NewResponse(resp)
}

func toPbPianoEdit(v *gateway.PianoEditView) *pianov1.PianoEdit {
	e := v.Edit
	pb := &pianov1.PianoEdit{
		Name:              resourcename.FormatPianoEdit(e.PianoID, e.ID),
		EditorDisplayName: v.EditorDisplayName,
		Operation:         toPbPianoEditOperation(e.Operation),
		Summary:           e.Summary,
		CreateTime:        timestamppb.New(e.CreateTime),
	}
	if v.EditorCustomID != "" {
		pb.Editor = resourcename.FormatUser(v.EditorCustomID)
	}
	if len(e.Changes) > 0 {
		s := string(e.Changes)
		pb.ChangesJson = &s
	}
	return pb
}

func toPbPianoEditOperation(op entity.PianoEditOperation) pianov1.PianoEditOperation {
	switch op {
	case entity.PianoEditOpCreate:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_CREATE
	case entity.PianoEditOpUpdate:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_UPDATE
	case entity.PianoEditOpPhotoAdd:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_PHOTO_ADD
	case entity.PianoEditOpPhotoRemove:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_PHOTO_REMOVE
	case entity.PianoEditOpStatusChange:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_STATUS_CHANGE
	case entity.PianoEditOpKindChange:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_KIND_CHANGE
	case entity.PianoEditOpRestore:
		return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_RESTORE
	}
	return pianov1.PianoEditOperation_PIANO_EDIT_OPERATION_UNSPECIFIED
}
