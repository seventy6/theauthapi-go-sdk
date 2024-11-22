package theauthapi

import (
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
        panic("Error loading .env file")
    }
}

func TestIsValidKey(t *testing.T) {
    // Get environment variables
    testAPIKeySuccess := os.Getenv("TEST_API_KEY_SUCCESS")
    testAPIKeyFail := os.Getenv("TEST_API_KEY_FAIL")
    accessToken := os.Getenv("ACCESS_TOKEN")
    //client := Client()
    client := NewClient(func(c *Client) { c.AccessToken = accessToken })
    
    // Create a mock server
    // Create a mock server
mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    apiKey := r.Header.Get("x-api-key")
    resp, err := client.ApiKeys.IsValidKey(r.Context(), apiKey)

    t.Logf("API Response - keyMatch: %v, response: %s", testAPIKeySuccess == apiKey, resp)
    
    if err != nil {
        t.Logf("Error validating key: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]bool{"valid": false})
        return
    }

    // Handle different cases without multiple WriteHeader calls
    if testAPIKeySuccess == apiKey {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]bool{"valid": true})
    } else if testAPIKeyFail == apiKey {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]bool{"valid": false})
    } else {
        t.Logf("reached here")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]bool{"valid": false})
    }
}))
        
    
    defer mockServer.Close()

    
    tests := []struct {
        name         string
        apiKey       string
        expectedCode int
    }{
        {
            name:         "valid key",
            apiKey:       testAPIKeySuccess,
            expectedCode: http.StatusOK,
        },
        {
            name:         "invalid key",
            apiKey:       testAPIKeyFail,
            expectedCode: http.StatusUnauthorized,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            req, err := http.NewRequest("GET", mockServer.URL, nil)
            if err != nil {
                t.Fatalf("failed to create request: %v", err)
            }

            req.Header.Set("x-api-key", test.apiKey)
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                t.Fatalf("failed to make request: %v", err)
            }
            defer resp.Body.Close()

            if resp.StatusCode != test.expectedCode {
                t.Errorf("expected status code %d, got %d", test.expectedCode, resp.StatusCode)
            }
        })
    }
}