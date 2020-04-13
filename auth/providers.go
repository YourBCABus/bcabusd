package auth

import (
	"github.com/go-pg/pg/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthProvider interface {
	Config() (*oauth2.Config, error)
	Authenticate(token *oauth2.Token, db *pg.DB)
}

type GoogleProvider struct {
	ClientID string
	ClientSecret string
	RedirectURI string
}

func (p GoogleProvider) Config() (*oauth2.Config, error) {
	return &oauth2.Config{
		ClientID: p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL: p.RedirectURI,
		Scopes: []string{"profile", "email"},
		Endpoint: google.Endpoint,
	}, nil
}

func (p GoogleProvider) Authenticate(token *oauth2.Token, db *pg.DB) {

}