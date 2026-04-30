package usecase

import "github.com/reverie-jp/piamap/internal/modules/user/gateway"

type UpdateUserOutput struct {
	View *gateway.UserView
}
