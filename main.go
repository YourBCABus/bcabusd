package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"

	"github.com/YourBCABus/bcabusd/api"
	"github.com/YourBCABus/bcabusd/auth"
	"github.com/YourBCABus/bcabusd/db"
	"github.com/YourBCABus/bcabusd/legacy"
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

	authRouter := router.PathPrefix("/auth").Subrouter()
	auth.ApplyRoutes(authRouter, db, auth.Config{
		Providers: map[string]auth.OAuthProvider{
			"google": &auth.GoogleProvider{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
			},
		},
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
	})

	router.Handle("/api", handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}))
	router.PathPrefix("/schools").Handler(legacy.Handler())
	router.HandleFunc("/teapot", teapot)
	router.HandleFunc("/", index)

	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)
}
