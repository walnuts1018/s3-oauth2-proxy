package usecase

import (
	"context"

	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository"
)

type AuthUsecase interface {
	GetAuthorizationURL(state, nonce string) string
	Login(ctx context.Context, code, expectedNonce string) (string, error)
}

type authUsecase struct {
	authRepo repository.AuthRepository
}

func NewAuthUsecase(authRepo repository.AuthRepository) AuthUsecase {
	return &authUsecase{authRepo: authRepo}
}

func (u *authUsecase) GetAuthorizationURL(state, nonce string) string {
	return u.authRepo.GetAuthorizationURL(state, nonce)
}

func (u *authUsecase) Login(ctx context.Context, code, expectedNonce string) (string, error) {
	token, err := u.authRepo.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	return u.authRepo.VerifyIDToken(ctx, token, expectedNonce)
}
