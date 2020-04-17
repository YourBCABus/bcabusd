package auth

import (
	"context"

	"github.com/coreos/go-oidc"
	"github.com/go-pg/pg/v9"
	"golang.org/x/oauth2"

	"github.com/YourBCABus/bcabusd/db"
)

// OAuthProvider defines methods for authenticating a user
// using an external OAuth server.
type OAuthProvider interface {
	// Config returns an OAuth2 configuration suitable for requesting
	// and retrieving tokens.
	Config(context context.Context) (*oauth2.Config, error)

	// Authenticate returns a matching user ID for a given OAuth2 token,
	// creating a new user as needed.
	Authenticate(context context.Context, token *oauth2.Token, db *pg.DB, createNewUser func(db.Meta) (string, error)) (string, error)
}

// GoogleProvider is an OAuthProvider that authenticates users
// with Google.
type GoogleProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	oidcProvider *oidc.Provider
}

func (p *GoogleProvider) setupProvider(ctx context.Context) error {
	if p.oidcProvider == nil {
		var err error
		p.oidcProvider, err = oidc.NewProvider(ctx, "https://accounts.google.com")
		if err != nil {
			return err
		}
	}

	return nil
}

// Config returns an OAuth2 configuration for Google.
func (p *GoogleProvider) Config(ctx context.Context) (*oauth2.Config, error) {
	if err := p.setupProvider(ctx); err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.RedirectURI,
		Scopes:       []string{"profile", "email"},
		Endpoint:     p.oidcProvider.Endpoint(),
	}, nil
}

// Authenticate returns a matching user ID for a given Google OAuth2 token,
// creating a new user as needed.
func (p *GoogleProvider) Authenticate(ctx context.Context, token *oauth2.Token, db *pg.DB, createNewUser func(db.Meta) (string, error)) (string, error) {
	if err := p.setupProvider(ctx); err != nil {
		return "", err
	}

	id, err := createNewUser(map[string]interface{}{"Hello": "World"})
	return id, err
}
