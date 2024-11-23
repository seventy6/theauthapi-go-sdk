package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/seventy6/theauthapi-go-sdk"
)

func apiKeyMiddleware(client *theauthapi.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("x-api-key")
			if apiKey == "" {
				http.Error(w, "No API key provided", http.StatusUnauthorized)
				return
			}

			err := client.ApiKeys.IsValidKey(r.Context(), apiKey)
			if err != nil {
				if !errors.Is(err, theauthapi.ErrKeyInvalid) {
					log.Printf("Error: %v", err)

					http.Error(w, "Error validating key", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
			}

			// Store the key in the request context
			ctx := context.WithValue(r.Context(), "key", apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value("key").(string)
	w.Write([]byte(key))
}

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	client := theauthapi.NewClient(
		theauthapi.WithAccessToken(accessToken),
	)

	mux := http.NewServeMux()

	apiMiddleware := apiKeyMiddleware(client)
	mux.Handle("/", handlers.LoggingHandler(os.Stdout, apiMiddleware(http.HandlerFunc(mainHandler))))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.CompressHandler(mux)); err != nil {
		log.Fatal(err)
	}
}
