package auth

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type redirectHandler struct {
	providers   map[string]OAuthProvider
	stateMaxAge int
	stateLength int
	jwtSecret   []byte
	jwtAudience string
}

func (h redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	provider, providerName := providerFor(w, r, h.providers)
	if provider == nil {
		return
	}

	config, err := provider.Config(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := ""

	stateLength := h.stateLength
	if stateLength == 0 {
		stateLength = 8
	}

	stateBytes := make([]byte, stateLength)
	_, err = rand.Read(stateBytes)
	if err != nil {
		fmt.Println("Error generating random state:", err)
		http.Error(w, "Error generating state", http.StatusInternalServerError)
	}

	now := time.Now().Unix()

	stateMaxAge := h.stateMaxAge
	if stateMaxAge == 0 {
		stateMaxAge = 3600
	}

	state = Base64Encoding.EncodeToString(stateBytes)
	stateToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: now + int64(stateMaxAge),
		Audience:  h.jwtAudience + "-" + providerName,
		NotBefore: now,
		Subject:   r.URL.Query().Get("login_challenge"),
		Id:        state,
	}).SignedString(h.jwtSecret)
	http.SetCookie(w, &http.Cookie{Name: stateCookieName(providerName), Value: stateToken, Path: "/auth", MaxAge: stateMaxAge})

	url := config.AuthCodeURL(stateToken)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
