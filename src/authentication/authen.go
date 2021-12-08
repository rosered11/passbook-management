package authentication

import (
	"context"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Authentication interface {
	GetProvider() (*oidc.Provider, *oauth2.Config, error)
}

type DefaultAuthentication struct {
	context context.Context
}

func (authen DefaultAuthentication) GetProvider() (*oidc.Provider, *oauth2.Config, error) {
	provider, err := oidc.NewProvider(authen.context, "https://identity-server-4213.herokuapp.com")
	if err != nil {
		return nil, nil, err
	}
	config := oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  os.Getenv("redirect_url"),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}
	return provider, &config, nil
}

func NewAuthentication(context context.Context) DefaultAuthentication {
	return DefaultAuthentication{context: context}
}
