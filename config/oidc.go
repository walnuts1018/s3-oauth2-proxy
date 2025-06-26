package config

import "golang.org/x/oauth2"

type OIDCConfig struct {
	IssuerURL    string `env:"OIDC_ISSUER_URL,required"`
	ClientID     string `env:"OIDC_CLIENT_ID,required"`
	ClientSecret string `env:"OIDC_CLIENT_SECRET,required"`
	RedirectURL  string `env:"OIDC_REDIRECT_URL,required"`
	AllowedGroup string `env:"OIDC_ALLOWED_GROUPS,required"`
	GroupClaim   string `env:"OIDC_GROUP_CLAIM" envDefault:"groups"`
}

func (c *OIDCConfig) ToOAuth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
	}
}
