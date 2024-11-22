package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/seventy6/theauthapi-go-sdk"
)

type ApiResponse struct {
    StatusCode int `json:"statusCode"`
    // Add other fields as needed
}

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(".env"); err != nil {
        panic("Error loading .env file")
    }
}

func apiKeyMiddleware(client *theauthapi.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            apiKey := r.Header.Get("x-api-key")
            if apiKey == "" {
                http.Error(w, "No API key provided", http.StatusUnauthorized)
                return
            }
            
            resp, err := client.ApiKeys.IsValidKey(r.Context(), apiKey)
            log.Printf("API Response: %+v", resp)
            if err != nil {
                log.Printf("Error: %v", err)
                http.Error(w, "Error validating key", http.StatusInternalServerError)
                return
            }

            var apiResponse ApiResponse
            err = json.Unmarshal([]byte(resp), &apiResponse)
            if err != nil {
                log.Printf("Error unmarshalling response: %v", err)
                http.Error(w, "Error processing response", http.StatusInternalServerError)
                return
            }
            if apiResponse.StatusCode != http.StatusOK {
                http.Error(w, "Invalid API key", http.StatusUnauthorized)
                return
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
    client := theauthapi.NewClient(func(c *theauthapi.Client) { c.AccessToken = accessToken })

    mux := http.NewServeMux()
    mux.Handle("/", handlers.LoggingHandler(os.Stdout, apiKeyMiddleware(client)(http.HandlerFunc(mainHandler))))
    log.Println("Starting server on :8080")

    if err := http.ListenAndServe(":8080", handlers.CompressHandler(mux)); err != nil {
        log.Fatal(err)
    }
}