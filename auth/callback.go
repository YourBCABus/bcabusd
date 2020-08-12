package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v9"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"

	"github.com/YourBCABus/bcabusd/db"
)

type callbackHandler struct {
	providers        map[string]OAuthProvider
	db               *pg.DB
	jwtSecret        []byte
	jwtAudience      string
	userJWTSecret    []byte
	userJWTAudience  string
	userJWTExpiresIn int64
	userJWTCookie    string
	userTokenLength  int
	hydraClient      admin.ClientService
	remember         bool
	rememberFor      int64
}

func (h callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	provider, providerName := providerFor(w, r, h.providers)
	if provider == nil {
		return
	}

	cookie, _ := r.Cookie(stateCookieName(providerName))
	state := r.URL.Query().Get("state")
	if cookie == nil || state != cookie.Value {
		http.Error(w, "Bad state", http.StatusBadRequest)
		return
	}

	stateToken, err := jwt.ParseWithClaims(state, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return h.jwtSecret, nil
	})

	if err != nil || !stateToken.Claims.(*jwt.StandardClaims).VerifyAudience(h.jwtAudience+"-"+providerName, true) {
		http.Error(w, "Bad state token", http.StatusBadRequest)
		return
	}

	challenge := stateToken.Claims.(*jwt.StandardClaims).Subject

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

	expiresIn := h.userJWTExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}

	now := time.Now().Unix()

	userToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: now + expiresIn,
		Audience:  h.userJWTAudience,
		NotBefore: now,
		Subject:   Base64Encoding.EncodeToString(user.UserToken),
	}).SignedString(h.userJWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: h.userJWTCookie, Value: userToken, Path: "/auth", MaxAge: int(expiresIn)})

	if challenge != "" {
		accept, err := h.hydraClient.AcceptLoginRequest(&admin.AcceptLoginRequestParams{LoginChallenge: challenge, Context: r.Context(), Body: &models.AcceptLoginRequest{
			Subject:     &id,
			Remember:    h.remember,
			RememberFor: h.rememberFor,
		}})
		if err != nil {
			http.Error(w, "Internal auth server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, *accept.Payload.RedirectTo, http.StatusSeeOther)
	} else {
		fmt.Fprintf(w, "Done! Logged in as user %v", id)
	}
}
