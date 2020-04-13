package auth

import (
	"fmt"
	"net/http"
)

type callbackHandler struct {
	checkState bool
	providers map[string]OAuthProvider
}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	provider, providerName := providerFor(w, r, h.providers)
	if provider == nil {
		return
	}

	if h.checkState {
		cookie, _ := r.Cookie(stateCookieName(providerName))
		if cookie == nil || r.URL.Query().Get("state") != cookie.Value {
			http.Error(w, "Bad state", http.StatusBadRequest)
			return
		}
	}

	config, err := provider.Config()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	fmt.Fprintf(w, "Callback! Token: %v", token.AccessToken)
}
