package xerrors

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
)

type Error struct {
	Code        string
	Message     string
	Cause       error
	ConnectCode connect.Code
}

func New(code, message string, connectCode connect.Code) *Error {
	return &Error{
		Code:        code,
		Message:     message,
		ConnectCode: connectCode,
	}
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) WithMessage(message string) *Error {
	return &Error{
		Code:        e.Code,
		Message:     message,
		Cause:       e.Cause,
		ConnectCode: e.ConnectCode,
	}
}

func (e *Error) WithCause(cause error) *Error {
	return &Error{
		Code:        e.Code,
		Message:     e.Message,
		Cause:       cause,
		ConnectCode: e.ConnectCode,
	}
}

// AsError extracts an *Error from err if one exists in the chain.
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// ToConnectError converts a domain error to a connect.Error suitable for
// returning from a handler. Non-domain errors are wrapped as Internal.
func ToConnectError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := AsError(err); ok {
		return connect.NewError(e.ConnectCode, e)
	}
	return connect.NewError(connect.CodeInternal, err)
}
