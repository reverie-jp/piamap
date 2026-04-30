package xerrors

import "connectrpc.com/connect"

var (
	ErrPianoNotFound = New("piano_not_found", "piano not found", connect.CodeNotFound)
	ErrPianoHidden   = New("piano_hidden", "piano is hidden", connect.CodeNotFound)
)
