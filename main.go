package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/ory/hydra-client-go/client"

	"github.com/YourBCABus/bcabusd/api"
	"github.com/YourBCABus/bcabusd/auth"
	"github.com/YourBCABus/bcabusd/db"
)

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("Method %v not allowed", req.Method), http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello, world!\n")
}

func teapot(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "I'm a teapot")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env: %v\n", err)
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	defer db.Close()

	schema, err := api.MakeSchema(db)
	if err != nil {
		log.Fatalf("failed to create schema: %v\n", err)
	}

	fmt.Println("Starting server...")

	router := mux.NewRouter()

	adminURL, err := url.Parse(os.Getenv("HYDRA_URL"))
	if err != nil {
		log.Fatalf("failed to parse Hydra URL: %v\n", err)
	}

	tmpl, err := template.ParseFiles("auth/login.html")
	if err != nil {
		panic(err)
	}

	rememberFor, _ := strconv.ParseInt(os.Getenv("HYDRA_REMEMBER"), 10, 64)

	authRouter := router.PathPrefix("/auth").Subrouter()
	auth.ApplyRoutes(authRouter, db, auth.Config{
		Providers: map[string]auth.OAuthProvider{
			"google": &auth.GoogleProvider{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
			},
		},
		JWTSecret:   []byte(os.Getenv("JWT_SECRET")),
		HydraClient: client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path}),
		Template:    tmpl,
		Remember:    os.Getenv("HYDRA_REMEMBER") != "",
		RememberFor: rememberFor,
	})

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.Handle("/api", handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}))
	router.HandleFunc("/teapot", teapot)
	router.HandleFunc("/", index)

	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)
}
