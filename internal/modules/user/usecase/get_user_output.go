package usecase

import "github.com/reverie-jp/piamap/internal/modules/user/gateway"

type GetUserOutput struct {
	View *gateway.UserView
}
