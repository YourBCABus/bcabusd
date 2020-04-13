package auth

import (
	"github.com/go-pg/pg/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthProvider defines methods for authenticating a user
// using an external OAuth server.
type OAuthProvider interface {
	// Config returns an OAuth2 configuration suitable for requesting
	// and retrieving tokens.
	Config() (*oauth2.Config, error)

	Authenticate(token *oauth2.Token, db *pg.DB)
}

// GoogleProvider is an OAuthProvider that authenticates users
// with Google.
type GoogleProvider struct {
	ClientID string
	ClientSecret string
	RedirectURI string
}

// Config returns an OAuth2 configuration for Google.
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