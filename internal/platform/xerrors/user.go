package xerrors

import "connectrpc.com/connect"

var (
	ErrUserNotFound        = New("user_not_found", "user not found", connect.CodeNotFound)
	ErrUserCustomIDInUse   = New("user_custom_id_in_use", "custom_id already in use", connect.CodeAlreadyExists)
	ErrUserSuspended       = New("user_suspended", "user is suspended", connect.CodePermissionDenied)
)
