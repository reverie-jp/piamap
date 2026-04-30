package xerrors

import "connectrpc.com/connect"

var (
	ErrSocialLoginFailed   = New("social_login_failed", "social login failed", connect.CodeUnauthenticated)
	ErrInvalidRefreshToken = New("invalid_refresh_token", "invalid refresh token", connect.CodeUnauthenticated)
	ErrAccountNotFound     = New("account_not_found", "account not found", connect.CodeNotFound)
	ErrCustomIDMismatch    = New("custom_id_mismatch", "custom_id confirmation mismatch", connect.CodePermissionDenied)
)
