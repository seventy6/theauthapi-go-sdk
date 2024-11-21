package theauthapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(".env"); err != nil {
        panic("Error 2 loading .env file")
    }
}

func TestIsValidKey(t *testing.T) {
    // Get environment variables
    testAPIKeySuccess := os.Getenv("TEST_API_KEY_SUCCESS")
    testAPIKeyFail := os.Getenv("TEST_API_KEY_FAIL")
    accessToken := os.Getenv("ACCESS_TOKEN")

    if testAPIKeySuccess == "" || testAPIKeyFail == "" || accessToken == "" {
        t.Fatal("Required environment variables not set")
    }

    // Create a mock server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify access token
        if r.Header.Get("x-api-key") != accessToken {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        // Extract key from URL path
        key := r.URL.Path[len("/api-keys/auth/"):]
        
        if key == testAPIKeySuccess {
            json.NewEncoder(w).Encode(map[string]bool{"valid": true})
        } else {
            json.NewEncoder(w).Encode(map[string]bool{"valid": false})
        }
    }))
    defer mockServer.Close()

    client := &Client{BaseURL: mockServer.URL, HTTPClient: mockServer.Client()}
    service := &ApiKeysService{client: client}

    tests := []struct {
        name     string
        key      string
        expected bool
    }{
        {
            name:     "valid key",
            key:      testAPIKeySuccess,
            expected: true,
        },
        {
            name:     "invalid key", 
            key:      testAPIKeyFail,
            expected: false,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            valid, err := service.IsValidKey(context.Background(), test.key)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if valid != test.expected {
                t.Errorf("expected %v, got %v", test.expected, valid)
            }
        })
    }
}

// func TestCreateKey(t *testing.T) {
//     // Create a mock server
//     mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         if r.URL.Path == "/keys" && r.Method == http.MethodPost {
//             var opts ApiKeyInput
//             if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
//                 http.Error(w, "bad request", http.StatusBadRequest)
//                 return
//             }
//             json.NewEncoder(w).Encode(ApiKey{Key: "new-key"})
//         } else {
//             http.Error(w, "not found", http.StatusNotFound)
//         }
//     }))
//     defer mockServer.Close()

//     client := &Client{BaseURL: mockServer.URL, HTTPClient: mockServer.Client()}
//     service := &ApiKeysService{client: client}

//     opts := ApiKeyInput{Name: "test-key"}
//     apiKey, err := service.CreateKey(context.Background(), opts)
//     if err != nil {
//         t.Fatalf("unexpected error: %v", err)
//     }
//     if apiKey.Key != "new-key" {
//         t.Errorf("expected new-key, got %v", apiKey.Key)
//     }
// }