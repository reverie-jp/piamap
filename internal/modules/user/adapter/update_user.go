package adapter

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	userv1 "github.com/reverie-jp/piamap/internal/gen/pb/user/v1"
	"github.com/reverie-jp/piamap/internal/modules/user/usecase"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromUpdateUserRequest(ctx context.Context, req *connect.Request[userv1.UpdateUserRequest]) (usecase.UpdateUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UpdateUserInput{}, xerrors.ErrUnauthenticated.WithCause(errors.New("missing user id"))
	}
	input := usecase.UpdateUserInput{RequesterID: requesterID}

	user := req.Msg.GetUser()
	if user == nil {
		return usecase.UpdateUserInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("user is required"))
	}
	mask := req.Msg.GetUpdateMask()
	if mask == nil {
		return usecase.UpdateUserInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("update_mask is required"))
	}

	// クライアントは FieldMask paths を camelCase で送るが、protojson は受信時に
	// proto field name (snake_case) へ変換するため、switch は snake_case で行う。
	for _, p := range mask.Paths {
		switch p {
		case "custom_id":
			input.SetCustomID = true
			input.CustomID = user.GetCustomId()
		case "display_name":
			input.SetDisplayName = true
			input.DisplayName = user.GetDisplayName()
		case "biography":
			input.SetBiography = true
			input.Biography = user.Biography
		case "avatar_url":
			input.SetAvatarURL = true
			input.AvatarURL = user.AvatarUrl
		case "hometown":
			input.SetHometown = true
			input.Hometown = user.Hometown
		case "piano_started_year":
			input.SetPianoStartedYear = true
			if user.PianoStartedYear != nil {
				v := int16(*user.PianoStartedYear)
				input.PianoStartedYear = &v
			}
		case "years_of_experience":
			input.SetYearsOfExperience = true
			if user.YearsOfExperience != nil {
				v := int16(*user.YearsOfExperience)
				input.YearsOfExperience = &v
			}
		default:
			return usecase.UpdateUserInput{}, xerrors.ErrInvalidArgument.WithCause(errors.New("unknown field in update_mask: " + p))
		}
	}
	return input, nil
}

func ToUpdateUserResponse(output *usecase.UpdateUserOutput) *connect.Response[userv1.UpdateUserResponse] {
	return connect.NewResponse(&userv1.UpdateUserResponse{User: ToUser(output.View)})
}
