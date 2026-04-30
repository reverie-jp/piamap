package usecase

import (
	"context"

	"github.com/reverie-jp/piamap/internal/domain/entity"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list/gateway"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

const (
	ListUserListPianosDefaultLimit = 20
	ListUserListPianosMaxLimit     = 100
)

type ListUserListPianosInput struct {
	RequesterID  ulid.ULID
	UserCustomID string
	ListKind     entity.PianoListKind
	PageSize     int32
	AfterPianoID *ulid.ULID
}

func (i ListUserListPianosInput) Validate() error {
	if i.UserCustomID == "" {
		return xerrors.ErrInvalidArgument.WithMessage("user custom_id is required")
	}
	if !i.ListKind.Valid() {
		return xerrors.ErrInvalidArgument.WithMessage("invalid list_kind")
	}
	return nil
}

type ListUserListPianosOutput struct {
	Views       []*pianogw.PianoView
	NextPianoID *ulid.ULID
}

type ListUserListPianos struct {
	gw           gateway.Gateway
	userGateway  usergw.Gateway
	pianoGateway pianogw.Gateway
}

func NewListUserListPianos(
	gw gateway.Gateway,
	userGateway usergw.Gateway,
	pianoGateway pianogw.Gateway,
) *ListUserListPianos {
	return &ListUserListPianos{gw: gw, userGateway: userGateway, pianoGateway: pianoGateway}
}

func (uc *ListUserListPianos) Execute(ctx context.Context, input ListUserListPianosInput) (*ListUserListPianosOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	user, err := uc.userGateway.GetUserByCustomID(ctx, input.UserCustomID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}

	limit := input.PageSize
	if limit <= 0 {
		limit = ListUserListPianosDefaultLimit
	} else if limit > ListUserListPianosMaxLimit {
		limit = ListUserListPianosMaxLimit
	}
	queryLimit := limit + 1

	pianoIDs, err := uc.gw.ListByUser(ctx, gateway.ListByUserParams{
		UserID:       user.ID,
		ListKind:     input.ListKind,
		AfterPianoID: input.AfterPianoID,
		Limit:        queryLimit,
	})
	if err != nil {
		return nil, err
	}

	var nextID *ulid.ULID
	if len(pianoIDs) > int(limit) {
		pianoIDs = pianoIDs[:limit]
		last := pianoIDs[len(pianoIDs)-1]
		nextID = &last
	}

	pianos := make([]*entity.Piano, 0, len(pianoIDs))
	for _, id := range pianoIDs {
		p, err := uc.pianoGateway.GetPiano(ctx, id)
		if err != nil {
			return nil, err
		}
		if p != nil {
			pianos = append(pianos, p)
		}
	}
	views, err := uc.pianoGateway.BuildListPianoViews(ctx, input.RequesterID, pianos)
	if err != nil {
		return nil, err
	}
	return &ListUserListPianosOutput{Views: views, NextPianoID: nextID}, nil
}
