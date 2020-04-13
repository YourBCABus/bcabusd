package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/gorilla/mux"
)

// Base64Encoding represents the base 64 encoding used for auth tokens.
var Base64Encoding = base64.StdEncoding

// Error is a type of error that encapsulates an HTTP status code and an error message.
type Error struct {
	Status int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth things")
}

// Config contains variables for configuring authentication routes.
type Config struct {
	// Providers is a map containing OAuthProviders to use.
	Providers map[string]OAuthProvider

	// StateMaxAge is the MaxAge of an OAuth state cookie.
	// The default is 3600 (1 hour).
	StateMaxAge int

	// StateLength is the length of the OAuth state cookie, in bytes.
	// The default is 12 bytes; set this to a negative value to disable
	// state checking.
	StateLength int
}

func providerFor(w http.ResponseWriter, r *http.Request, providers map[string]OAuthProvider) (OAuthProvider, string) {
	providerNames := r.URL.Query()["provider"]
	if len(providerNames) != 1 {
		http.Error(w, "Must provide exactly 1 auth provider", http.StatusBadRequest)
		return nil, ""
	}

	providerName := providerNames[0]
	provider := providers[providerName]
	if provider == nil {
		http.Error(w, fmt.Sprintf("%v is not a provider", providerName), http.StatusBadRequest)
		return nil, providerName
	}

	return provider, providerName
}

func stateCookieName(providerName string) string {
	return fmt.Sprintf("%s-external-auth-state", providerName)
}

// ApplyRoutes applies authentication-related routes
// to a given router using a given database.
func ApplyRoutes(router *mux.Router, db *pg.DB, config Config) {
	router.Handle("/redirect", redirectHandler{providers: config.Providers, stateMaxAge: config.StateMaxAge, stateLength: config.StateLength})
	router.Handle("/callback", callbackHandler{checkState: config.StateLength >= 0, providers: config.Providers})
	router.HandleFunc("", index)
}
