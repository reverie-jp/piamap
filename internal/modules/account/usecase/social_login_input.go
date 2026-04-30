package usecase

import "errors"

type SocialLoginInput struct {
	Provider string
	Code     string
}

func (i SocialLoginInput) Validate() error {
	if i.Provider != "google" {
		return errors.New("provider must be 'google'")
	}
	if i.Code == "" {
		return errors.New("code is required")
	}
	return nil
}
