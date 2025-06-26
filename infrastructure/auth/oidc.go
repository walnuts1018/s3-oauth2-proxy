package auth

import (
	"context"
	"errors"

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

func (r *authRepository) GetAuthorizationURL(state string) string {
	return r.oauth2Config.AuthCodeURL(state)
}

func (r *authRepository) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return r.oauth2Config.Exchange(ctx, code)
}

func (r *authRepository) VerifyIDToken(ctx context.Context, token *oauth2.Token) (string, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("id_token not found")
	}

	verifier := r.provider.Verifier(&oidc.Config{ClientID: r.cfg.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	var claims struct {
		Groups []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", err
	}

	for _, group := range claims.Groups {
		if group == r.cfg.AllowedGroup {
			return idToken.Subject, nil
		}
	}

	return "", errors.New("not in allowed group")
}
