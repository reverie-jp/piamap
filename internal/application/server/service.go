package server

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/config"
	"github.com/reverie-jp/piamap/internal/gen/pb/account/v1/accountv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/user/v1/userv1connect"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/account"
	"github.com/reverie-jp/piamap/internal/modules/user"
	usergw "github.com/reverie-jp/piamap/internal/modules/user/gateway"
	"github.com/reverie-jp/piamap/internal/platform/google"
	"github.com/reverie-jp/piamap/internal/platform/jwt"
)

type Service struct {
	Name                   string
	RegisterConnectHandler func(mux *http.ServeMux)
}

func initServices(cfg *config.Config, db *pgxpool.Pool, jwtManager *jwt.Manager) []Service {
	q := sqlc.New(db)
	tx := transaction.NewRunner(db)
	googleAuth := google.NewAuthClient(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL)

	errorInterceptor := interceptor.ErrorInterceptor(cfg.Env)
	authInterceptor := interceptor.AuthInterceptor(jwtManager)

	userGateway := usergw.New(q)
	accountService := account.InitModule(q, userGateway, tx, googleAuth, jwtManager)
	userService := user.InitModule(userGateway)

	return []Service{
		{
			Name: accountv1connect.AccountServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(accountv1connect.NewAccountServiceHandler(
					accountService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
		{
			Name: userv1connect.UserServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(userv1connect.NewUserServiceHandler(
					userService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
	}
}
