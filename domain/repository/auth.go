package repository

import (
	"context"

	"golang.org/x/oauth2"
)

//go:generate go tool mockgen -source=auth.go -destination=mock_repository/mock_auth_repository.go -package=mock_repository
type AuthRepository interface {
	GetAuthorizationURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	VerifyIDToken(ctx context.Context, token *oauth2.Token, expectedNonce string) (string, error)
}