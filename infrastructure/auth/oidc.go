package auth

import (
	"context"
	"errors"
	"slices"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/walnuts1018/s3-oauth2-proxy/config"
	"github.com/walnuts1018/s3-oauth2-proxy/domain/repository"
	"golang.org/x/oauth2"
)

type authRepository struct {
	cfg          config.OIDCConfig
	provider     *oidc.Provider
	oauth2Config oauth2.Config
}

func NewAuthRepository(cfg config.OIDCConfig) (repository.AuthRepository, error) {
	oauth2Config := cfg.ToOAuth2Config()
	provider, err := oidc.NewProvider(context.Background(), cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauth2Config.Endpoint = provider.Endpoint()

	return &authRepository{
		cfg:          cfg,
		provider:     provider,
		oauth2Config: oauth2Config,
	}, nil
}

func (r *authRepository) GetAuthorizationURL(state, nonce string) string {
	return r.oauth2Config.AuthCodeURL(state, oauth2.SetAuthURLParam("nonce", nonce))
}

func (r *authRepository) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return r.oauth2Config.Exchange(ctx, code)
}

func (r *authRepository) VerifyIDToken(ctx context.Context, token *oauth2.Token, expectedNonce string) (string, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("id_token not found")
	}

	verifier := r.provider.Verifier(&oidc.Config{ClientID: r.cfg.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	if idToken.Nonce != expectedNonce {
		return "", errors.New("invalid nonce")
	}

	userinfo, err := r.provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return "", err
	}

	if len(r.cfg.AllowedGroups) == 0 {
		return idToken.Subject, nil
	}

	var claims map[string]any
	if err := userinfo.Claims(&claims); err != nil {
		return "", err
	}

	groupsClaim, ok := claims[r.cfg.GroupClaim].([]any)
	if !ok {
		return "", errors.New("group claim not found or not a list")
	}

	for _, group := range groupsClaim {
		if groupStr, ok := group.(string); ok && slices.Contains(r.cfg.AllowedGroups, groupStr) {
			return idToken.Subject, nil
		}
	}

	return "", errors.New("not in allowed group")
}
