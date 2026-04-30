package server

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reverie-jp/piamap/internal/application/server/interceptor"
	"github.com/reverie-jp/piamap/internal/application/transaction"
	"github.com/reverie-jp/piamap/internal/config"
	"github.com/reverie-jp/piamap/internal/gen/pb/account/v1/accountv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano/v1/pianov1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1/piano_postv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_comment/v1/piano_post_commentv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_like/v1/piano_post_likev1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1/piano_user_listv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/user/v1/userv1connect"
	"github.com/reverie-jp/piamap/internal/gen/sqlc"
	"github.com/reverie-jp/piamap/internal/modules/account"
	"github.com/reverie-jp/piamap/internal/modules/piano"
	pianogw "github.com/reverie-jp/piamap/internal/modules/piano/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post"
	postgw "github.com/reverie-jp/piamap/internal/modules/piano_post/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_comment"
	commentgw "github.com/reverie-jp/piamap/internal/modules/piano_post_comment/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_post_like"
	likegw "github.com/reverie-jp/piamap/internal/modules/piano_post_like/gateway"
	"github.com/reverie-jp/piamap/internal/modules/piano_user_list"
	listgw "github.com/reverie-jp/piamap/internal/modules/piano_user_list/gateway"
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
	pianoGateway := pianogw.New(q, userGateway)
	pianoPostLikeGateway := likegw.New(q)
	pianoPostGateway := postgw.New(q, userGateway, pianoGateway, pianoPostLikeGateway)
	pianoPostCommentGateway := commentgw.New(q, userGateway, pianoPostGateway)
	pianoUserListGateway := listgw.New(q)

	accountService := account.InitModule(q, userGateway, tx, googleAuth, jwtManager)
	userService := user.InitModule(userGateway)
	pianoService := piano.InitModule(pianoGateway, userGateway, tx)
	pianoPostService := piano_post.InitModule(pianoPostGateway, pianoGateway, userGateway, tx)
	pianoPostLikeService := piano_post_like.InitModule(pianoPostLikeGateway, pianoPostGateway, userGateway)
	pianoPostCommentService := piano_post_comment.InitModule(pianoPostCommentGateway, pianoPostGateway, userGateway)
	pianoUserListService := piano_user_list.InitModule(pianoUserListGateway, pianoGateway, userGateway)

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
		{
			Name: pianov1connect.PianoServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(pianov1connect.NewPianoServiceHandler(
					pianoService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
		{
			Name: piano_postv1connect.PianoPostServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(piano_postv1connect.NewPianoPostServiceHandler(
					pianoPostService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
		{
			Name: piano_user_listv1connect.PianoUserListServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(piano_user_listv1connect.NewPianoUserListServiceHandler(
					pianoUserListService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
		{
			Name: piano_post_likev1connect.PianoPostLikeServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(piano_post_likev1connect.NewPianoPostLikeServiceHandler(
					pianoPostLikeService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
		{
			Name: piano_post_commentv1connect.PianoPostCommentServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(piano_post_commentv1connect.NewPianoPostCommentServiceHandler(
					pianoPostCommentService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
	}
}
