package usecase

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	accountrepo "github.com/reverie-jp/piamap/internal/modules/account/repository"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/google"
	"github.com/reverie-jp/piamap/internal/platform/jwt"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

type SocialLogin struct {
	accountRepo accountrepo.Repository
	userGateway usergw.Gateway
	tx          transaction.Runner
	googleAuth  *google.AuthClient
	jwtManager  *jwt.Manager
}

func NewSocialLogin(
	accountRepo accountrepo.Repository,
	userGateway usergw.Gateway,
	tx transaction.Runner,
	googleAuth *google.AuthClient,
	jwtManager *jwt.Manager,
) *SocialLogin {
	return &SocialLogin{
		accountRepo: accountRepo,
		userGateway: userGateway,
		tx:          tx,
		googleAuth:  googleAuth,
		jwtManager:  jwtManager,
	}
}

func (uc *SocialLogin) Execute(ctx context.Context, input SocialLoginInput) (*SocialLoginOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	userInfo, err := uc.googleAuth.Exchange(ctx, input.Code)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	authProvider, err := uc.accountRepo.GetAuthProviderByProvider(ctx, input.Provider, userInfo.Sub)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	var userID ulid.ULID
	isNewAccount := authProvider == nil

	if isNewAccount {
		userID, err = uc.createNewUser(ctx, userInfo)
		if err != nil {
			return nil, err
		}
	} else {
		userID = authProvider.UserID
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	refreshToken, refreshExpireTime, err := uc.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	if err := uc.accountRepo.CreateRefreshToken(ctx, accountrepo.CreateRefreshTokenParams{
		UserID:     userID,
		RawToken:   refreshToken,
		ExpireTime: refreshExpireTime,
	}); err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	return &SocialLoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsNewAccount: isNewAccount,
	}, nil
}

func (uc *SocialLogin) createNewUser(ctx context.Context, userInfo *google.UserInfo) (ulid.ULID, error) {
	userID := ulid.New()

	customID, err := generateCustomID()
	if err != nil {
		return ulid.ULID{}, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	displayName := userInfo.Name
	if displayName == "" {
		displayName = "unknown"
	}

	var avatarURL *string
	if userInfo.Picture != "" {
		// TODO: R2 にコピーして安定 URL 化する。MVP は Google の URL を直接使う。
		avatarURL = &userInfo.Picture
	}

	err = uc.tx.WithTx(ctx, func(q sqlc.Querier) error {
		txUserGw := usergw.New(q)
		txAccountRepo := accountrepo.New(q)

		if err := txUserGw.CreateUser(ctx, usergw.CreateUserParams{
			ID:          userID,
			CustomID:    customID,
			DisplayName: displayName,
			AvatarURL:   avatarURL,
		}); err != nil {
			return err
		}
		return txAccountRepo.CreateAuthProvider(ctx, accountrepo.CreateAuthProviderParams{
			UserID:         userID,
			Provider:       "google",
			ProviderUserID: userInfo.Sub,
		})
	})
	if err != nil {
		return ulid.ULID{}, err
	}
	return userID, nil
}

// generateCustomID: a-z0-9 の 10 文字。重複時のリトライは未実装(衝突確率は十分低い)。
func generateCustomID() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 10
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}
