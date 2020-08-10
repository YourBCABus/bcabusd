package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"

	"github.com/YourBCABus/bcabusd/db"
	"github.com/go-pg/pg/v9"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
)

type fositeStore struct {
	db *pg.DB
}

func (s *fositeStore) GetClient(_ context.Context, id string) (fosite.Client, error) {
	client := &db.AuthClient{ClientID: id}
	err := s.db.Select(client)
	if err != nil {
		return nil, err
	}

	return &fosite.DefaultClient{
		ID:           client.ClientID,
		Secret:       client.ClientSecret,
		RedirectURIs: client.RedirectURIs,
		GrantTypes:   client.GrantTypes,
	}, nil
}

type FositeConfig struct {
	Secret     []byte
	RSAKeySize int
}

func FositeProvider(config FositeConfig, db *pg.DB) (fosite.OAuth2Provider, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, config.RSAKeySize)
	if err != nil {
		return nil, err
	}

	return compose.ComposeAllEnabled(&compose.Config{}, &fositeStore{db}, config.Secret, privateKey), nil
}
