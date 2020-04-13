package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/gorilla/mux"
)

var Base64Encoding = base64.StdEncoding

type Error struct {
	status int
	message string
}

func (e Error) Error() string {
	return e.message
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth things")
}

type Config struct {
	Providers map[string]OAuthProvider
	StateMaxAge int
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
