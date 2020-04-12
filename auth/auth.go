package auth

import (
	"crypto/rand"
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

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth things")
}

type Config struct {
	Providers map[string]OAuthProvider
	StateMaxAge int
	StateLength int
}

type RedirectHandler struct {
	providers map[string]OAuthProvider
	stateMaxAge int
	stateLength int
}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	providerNames := r.URL.Query()["provider"]
	if len(providerNames) != 1 {
		http.Error(w, "Must provide exactly 1 auth provider", http.StatusBadRequest)
		return
	}

	providerName := providerNames[0]
	provider := h.providers[providerName]
	if provider == nil {
		http.Error(w, fmt.Sprintf("%v is not a provider", providerName), http.StatusBadRequest)
		return
	}

	config, err := provider.Config()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := ""

	stateLength := h.stateLength
	if stateLength == 0 {
		stateLength = 12
	}

	if stateLength > 0 {
		stateBytes := make([]byte, stateLength)
		_, err := rand.Read(stateBytes)
		if err != nil {
			fmt.Println("Error generating random state:", err)
			http.Error(w, "Error generating state", http.StatusInternalServerError)
		}

		stateMaxAge := h.stateMaxAge
		if stateMaxAge == 0 {
			stateMaxAge = 3600
		}

		state = Base64Encoding.EncodeToString(stateBytes)
		http.SetCookie(w, &http.Cookie{Name: fmt.Sprintf("%s-external-auth-state", providerName), Value: state, Path: "/auth", MaxAge: stateMaxAge})
	}

	url := config.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// ApplyRoutes applies authentication-related routes
// to a given router using a given database.
func ApplyRoutes(router *mux.Router, db *pg.DB, config Config) {
	router.Handle("/redirect", RedirectHandler{providers: config.Providers, stateMaxAge: config.StateMaxAge, stateLength: config.StateLength})
	router.HandleFunc("", index)
}
