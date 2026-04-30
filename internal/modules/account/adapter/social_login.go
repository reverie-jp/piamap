package adapter

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	accountv1 "github.com/reverie-jp/piamap/internal/gen/pb/account/v1"
	"github.com/reverie-jp/piamap/internal/modules/account/usecase"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func FromSocialLoginRequest(_ context.Context, req *connect.Request[accountv1.SocialLoginRequest]) (usecase.SocialLoginInput, error) {
	provider, err := parseProvider(req.Msg.Provider)
	if err != nil {
		return usecase.SocialLoginInput{}, err
	}
	return usecase.SocialLoginInput{
		Provider: provider,
		Code:     req.Msg.Code,
	}, nil
}

func ToSocialLoginResponse(output *usecase.SocialLoginOutput) *connect.Response[accountv1.SocialLoginResponse] {
	return connect.NewResponse(&accountv1.SocialLoginResponse{
		TokenPair: &accountv1.TokenPair{
			AccessToken:  output.AccessToken,
			RefreshToken: output.RefreshToken,
		},
		IsNewAccount: output.IsNewAccount,
	})
}

func parseProvider(p accountv1.AuthProvider) (string, error) {
	switch p {
	case accountv1.AuthProvider_AUTH_PROVIDER_GOOGLE:
		return "google", nil
	default:
		return "", xerrors.ErrInvalidArgument.WithCause(errors.New("unsupported provider: " + strings.ToLower(p.String())))
	}
}
