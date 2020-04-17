package auth

import (
	"context"
	"errors"

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

func (p *GoogleProvider) setupProvider() error {
	if p.oidcProvider == nil {
		var err error
		p.oidcProvider, err = oidc.NewProvider(context.Background(), "https://accounts.google.com")
		if err != nil {
			return err
		}
	}

	return nil
}

// Config returns an OAuth2 configuration for Google.
func (p *GoogleProvider) Config(ctx context.Context) (*oauth2.Config, error) {
	if err := p.setupProvider(); err != nil {
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
func (p *GoogleProvider) Authenticate(ctx context.Context, token *oauth2.Token, conn *pg.DB, createNewUser func(db.Meta) (string, error)) (string, error) {
	if err := p.setupProvider(); err != nil {
		return "", err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("Missing ID token")
	}

	verifier := p.oidcProvider.Verifier(&oidc.Config{ClientID: p.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	var claims struct {
		Subject  string `json:"sub"`
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", err
	}

	var providers []db.AuthProvider
	if err := conn.Model(&providers).Column("user_id").Where("provider = ?", "google").Where("subject = ?", claims.Subject).Select(); err != nil {
		return "", err
	}

	if len(providers) == 0 {
		id, err := createNewUser(map[string]interface{}{"email": claims.Email})
		if err != nil {
			return id, err
		}

		if err := conn.Insert(&db.AuthProvider{
			UserID:   id,
			Provider: "google",
			Email:    claims.Email,
			Subject:  claims.Subject,
			Meta:     map[string]interface{}{"email_verified": claims.Verified},
		}); err != nil {
			return "", err
		}

		return id, nil
	}

	return providers[0].UserID, nil
}
