package handler

import (
	"context"

	"connectrpc.com/connect"

	piano_user_listv1 "github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/adapter"
)

func (h *Handler) AddPianoToUserList(ctx context.Context, req *connect.Request[piano_user_listv1.AddPianoToUserListRequest]) (*connect.Response[piano_user_listv1.AddPianoToUserListResponse], error) {
	input, err := adapter.FromAddRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.add.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToAddResponse(output), nil
}

func (h *Handler) RemovePianoFromUserList(ctx context.Context, req *connect.Request[piano_user_listv1.RemovePianoFromUserListRequest]) (*connect.Response[piano_user_listv1.RemovePianoFromUserListResponse], error) {
	input, err := adapter.FromRemoveRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.remove.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToRemoveResponse(output), nil
}

func (h *Handler) ListUserListPianos(ctx context.Context, req *connect.Request[piano_user_listv1.ListUserListPianosRequest]) (*connect.Response[piano_user_listv1.ListUserListPianosResponse], error) {
	input, err := adapter.FromListRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listPianos.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListResponse(output), nil
}

func (h *Handler) GetMyPianoUserLists(ctx context.Context, req *connect.Request[piano_user_listv1.GetMyPianoUserListsRequest]) (*connect.Response[piano_user_listv1.GetMyPianoUserListsResponse], error) {
	input, err := adapter.FromGetMyRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getMine.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetMyResponse(output), nil
}
