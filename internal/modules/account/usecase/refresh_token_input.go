package usecase

import "errors"

type RefreshTokenInput struct {
	RefreshToken string
}

func (i RefreshTokenInput) Validate() error {
	if i.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}
	return nil
}
