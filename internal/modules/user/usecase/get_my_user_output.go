package usecase

import "github.com/reverie-jp/piamap/internal/modules/user/gateway"

type GetMyUserOutput struct {
	View *gateway.UserView
}
