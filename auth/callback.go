package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v9"

	"github.com/YourBCABus/bcabusd/db"
)

type callbackHandler struct {
	checkState      bool
	providers       map[string]OAuthProvider
	db              *pg.DB
	userTokenLength int
	jwtSecret       []byte
	jwtExpiresIn    int64
	jwtAudience     string
	jwtCookieName   string
	jwtCookieMaxAge int
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

	config, err := provider.Config(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	id, err := provider.Authenticate(r.Context(), token, h.db, func(meta db.Meta) (string, error) {
		user := db.User{IsBot: false, Meta: meta}
		tokenLength := h.userTokenLength
		if tokenLength == 0 {
			tokenLength = 24
		}
		_, err := h.db.Model(&user).Value("user_token", fmt.Sprintf("gen_random_bytes(%d)", tokenLength)).Returning("id").Insert()
		if err != nil {
			return "", err
		}

		return user.ID, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if id == "" {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	user := db.User{ID: id}
	if err := h.db.Model(&user).Column("user_token").Where("id = ?id").Select(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiresIn := h.jwtExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}

	now := time.Now().Unix()

	userToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: now + expiresIn,
		Audience:  h.jwtAudience,
		NotBefore: now,
		Subject:   Base64Encoding.EncodeToString(user.UserToken),
	}).SignedString(h.jwtSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: h.jwtCookieName, Value: userToken, Path: "/auth", MaxAge: h.jwtCookieMaxAge})
	fmt.Fprintf(w, "Done!")
}
