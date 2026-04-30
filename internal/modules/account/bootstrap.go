package account

import (
	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/gen/pb/account/v1/accountv1connect"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/account/handler"
	accountrepo "github.com/reverie-jp/piamap/internal/modules/account/repository"
	"github.com/reverie-jp/piamap/internal/modules/account/usecase"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/google"
	"github.com/reverie-jp/piamap/internal/platform/jwt"
)

func InitModule(
	q sqlc.Querier,
	userGateway usergw.Gateway,
	tx transaction.Runner,
	googleAuth *google.AuthClient,
	jwtManager *jwt.Manager,
) accountv1connect.AccountServiceHandler {
	repo := accountrepo.New(q)
	return handler.New(
		usecase.NewSocialLogin(repo, userGateway, tx, googleAuth, jwtManager),
		usecase.NewRefreshToken(repo, tx, jwtManager),
		usecase.NewLogout(repo),
		usecase.NewDeleteAccount(userGateway),
	)
}
