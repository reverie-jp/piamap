package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/reverie-jp/piamap/internal/config"
	"github.com/reverie-jp/piamap/internal/platform/xerrors"
)

func ErrorInterceptor(env config.Env) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				return nil, connectErr
			}

			var appErr *xerrors.Error
			if errors.As(err, &appErr) {
				return nil, connect.NewError(appErr.ConnectCode, appErr)
			}

			slog.ErrorContext(ctx, "unhandled error", slog.String("error", err.Error()))
			if env == config.EnvDevelopment {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
		}
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
