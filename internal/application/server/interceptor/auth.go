package interceptor

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/gen/pb/account/v1/accountv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano/v1/pianov1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post/v1/piano_postv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_comment/v1/piano_post_commentv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_post_like/v1/piano_post_likev1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/piano_user_list/v1/piano_user_listv1connect"
	"github.com/reverie-jp/piamap/internal/gen/pb/user/v1/userv1connect"
	"github.com/reverie-jp/piamap/internal/platform/jwt"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

type userIDKey struct{}

func ContextWithUserID(ctx context.Context, userID ulid.ULID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (ulid.ULID, bool) {
	v := ctx.Value(userIDKey{})
	if v == nil {
		return ulid.ULID{}, false
	}
	id, ok := v.(ulid.ULID)
	return id, ok
}

// publicProcedures: 認証スキップ。Authorization ヘッダがあっても無視する。
var publicProcedures = map[string]bool{
	accountv1connect.AccountServiceSocialLoginProcedure:  true,
	accountv1connect.AccountServiceRefreshTokenProcedure: true,
}

// optionalAuthProcedures: 認証はオプショナル。ヘッダ無しは guest 通過、有れば検証必須。
var optionalAuthProcedures = map[string]bool{
	userv1connect.UserServiceGetUserProcedure:                                    true,
	pianov1connect.PianoServiceGetPianoProcedure:                                 true,
	pianov1connect.PianoServiceSearchPianosProcedure:                             true,
	pianov1connect.PianoServiceListPianoEditsProcedure:                           true,
	piano_postv1connect.PianoPostServiceGetPianoPostProcedure:                    true,
	piano_postv1connect.PianoPostServiceListPianoPostsProcedure:                  true,
	piano_user_listv1connect.PianoUserListServiceListUserListPianosProcedure:     true,
	piano_post_likev1connect.PianoPostLikeServiceListLikedPianoPostsProcedure:    true,
	piano_post_commentv1connect.PianoPostCommentServiceListPianoPostCommentsProcedure: true,
}

type authInterceptor struct {
	jwtManager *jwt.Manager
}

func AuthInterceptor(jwtManager *jwt.Manager) connect.Interceptor {
	return &authInterceptor{jwtManager: jwtManager}
}

func (a *authInterceptor) authenticate(ctx context.Context, procedure string, header http.Header) (context.Context, error) {
	if publicProcedures[procedure] {
		return ctx, nil
	}

	rawHeader := header.Get("Authorization")
	optional := optionalAuthProcedures[procedure]
	if rawHeader == "" && optional {
		return ctx, nil
	}

	token, err := extractBearerToken(rawHeader)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("missing or invalid authorization header"))
	}

	claims, err := a.jwtManager.VerifyToken(token)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid or expired token"))
	}

	if claims.TokenType != jwt.TokenTypeAccess {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token type"))
	}

	userID, err := ulid.Parse(claims.Subject)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
	}

	return ContextWithUserID(ctx, userID), nil
}

func (a *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx, err := a.authenticate(ctx, req.Spec().Procedure, req.Header())
		if err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

func (a *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		ctx, err := a.authenticate(ctx, conn.Spec().Procedure, conn.RequestHeader())
		if err != nil {
			return err
		}
		return next(ctx, conn)
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("empty authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
