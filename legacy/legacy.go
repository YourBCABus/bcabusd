package legacy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the legacy API! Please use the GraphQL API instead.")
}

// Handler returns a handler for
// the legacy YourBCABus API.
func Handler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/schools", index)
	return r
}
