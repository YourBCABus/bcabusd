package auth

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

type redirectHandler struct {
	providers map[string]OAuthProvider
	stateMaxAge int
	stateLength int
}

func (h redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	provider, providerName := providerFor(w, r, h.providers)
	if provider == nil {
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
		http.SetCookie(w, &http.Cookie{Name: stateCookieName(providerName), Value: state, Path: "/auth", MaxAge: stateMaxAge})
	}

	url := config.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusSeeOther)
}