package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/handler"

	"github.com/YourBCABus/bcabusd/api"
)

func index(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("Method %v not allowed", req.Method), http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello, world!\n")
}

func main() {
	schema, err := api.MakeSchema()
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	fmt.Println("Starting server...")
	http.HandleFunc("/", index)
	http.Handle("/api", handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}))
	http.ListenAndServe(":3000", nil)
}
