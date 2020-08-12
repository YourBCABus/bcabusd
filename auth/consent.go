package auth

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

type consentHandler struct {
	client      admin.ClientService
	template    *template.Template
	remember    bool
	rememberFor int64
}

type consentTemplateArgs struct {
	CSRFToken interface{}
}

func (h consentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("consent_challenge")
	if challenge == "" {
		http.Error(w, "Consent challenge required", http.StatusBadRequest)
		return
	}

	consentRequest, err := h.client.GetConsentRequest(&admin.GetConsentRequestParams{ConsentChallenge: challenge, Context: r.Context()})
	if err != nil {
		http.Error(w, "Bad consent challenge", http.StatusBadRequest)
		return
	}

	if consentRequest.Payload.Skip {
		accept, err := h.client.AcceptConsentRequest(&admin.AcceptConsentRequestParams{ConsentChallenge: challenge, Context: r.Context(), Body: &models.AcceptConsentRequest{
			GrantScope:  consentRequest.Payload.RequestedScope,
			Remember:    h.remember,
			RememberFor: h.rememberFor,
		}})
		if err != nil {
			http.Error(w, "Internal auth server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, *accept.Payload.RedirectTo, http.StatusSeeOther)
	} else {
		err := h.template.ExecuteTemplate(w, "consent.html", &consentTemplateArgs{csrf.TemplateField(r)})
		if err != nil {
			http.Error(w, "Internal server error (template)", http.StatusInternalServerError)
		}
	}
}
