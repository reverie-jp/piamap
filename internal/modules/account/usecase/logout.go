package usecase

import (
	"context"

	accountrepo "github.com/reverie-jp/piamap/internal/modules/account/repository"
)

type Logout struct {
	accountRepo accountrepo.Repository
}

func NewLogout(accountRepo accountrepo.Repository) *Logout {
	return &Logout{accountRepo: accountRepo}
}

func (uc *Logout) Execute(ctx context.Context, input LogoutInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	// Idempotent: missing rows succeed silently、user_id scoping で他人のトークンは消せない。
	return uc.accountRepo.DeleteRefreshTokenByRaw(ctx, input.RefreshToken, input.UserID)
}
