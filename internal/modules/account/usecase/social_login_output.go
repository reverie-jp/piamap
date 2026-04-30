package usecase

type SocialLoginOutput struct {
	AccessToken  string
	RefreshToken string
	IsNewAccount bool
}
