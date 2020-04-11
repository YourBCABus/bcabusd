package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth things")
}

// ApplyRoutes applies authentication-related routes
// to a given router.
func ApplyRoutes(router *mux.Router) {
	router.HandleFunc("", index)
}
