package main

import (
	"log"
	"net/http"

	"github.com/seventy6/theauthapi-go-sdk/theauthapi"
)

func main() {
    client := theauthapi.NewClient("your-api-key")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("x-api-key")
        if apiKey == "" {
            http.Error(w, "No API key provided", http.StatusUnauthorized)
            return
        }

        isValid, err := client.ApiKeys.IsValidKey(r.Context(), apiKey)
        if err != nil {
            log.Printf("Error: %v", err)
            http.Error(w, "Error validating key", http.StatusInternalServerError)
            return
        }

        if !isValid {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
        }
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}