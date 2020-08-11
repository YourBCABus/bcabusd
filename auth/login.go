package auth

import (
	"html/template"
	"net/http"

	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

type loginHandler struct {
	client   admin.ClientService
	template *template.Template
}

func (h loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("login_challenge")
	if challenge == "" {
		http.Error(w, "Login challenge required", http.StatusBadRequest)
		return
	}

	loginRequest, err := h.client.GetLoginRequest(&admin.GetLoginRequestParams{LoginChallenge: challenge, Context: r.Context()})
	if err != nil {
		http.Error(w, "Bad login challenge", http.StatusBadRequest)
		return
	}

	if *loginRequest.Payload.Skip {
		accept, err := h.client.AcceptLoginRequest(&admin.AcceptLoginRequestParams{LoginChallenge: challenge, Context: r.Context(), Body: &models.AcceptLoginRequest{
			Subject: loginRequest.Payload.Subject,
		}})
		if err != nil {
			http.Error(w, "Internal auth server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, *accept.Payload.RedirectTo, http.StatusSeeOther)
	} else {
		err := h.template.ExecuteTemplate(w, "login.html", struct{}{})
		if err != nil {
			http.Error(w, "Internal server error (template)", http.StatusInternalServerError)
		}
	}
}
