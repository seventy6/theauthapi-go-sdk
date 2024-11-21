package theauthapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsValidKey(t *testing.T) {
    // Create a mock server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/validate" && r.Method == http.MethodPost {
            key := r.URL.Query().Get("key")
            if key == "valid-key" {
                json.NewEncoder(w).Encode(map[string]bool{"valid": true})
            } else {
                json.NewEncoder(w).Encode(map[string]bool{"valid": false})
            }
        } else {
            http.Error(w, "not found", http.StatusNotFound)
        }
    }))
    defer mockServer.Close()

    client := &Client{BaseURL: mockServer.URL, HTTPClient: mockServer.Client()}
    service := &ApiKeysService{client: client}

    tests := []struct {
        key      string
        expected bool
    }{
        {"valid-key", true},
        {"invalid-key", false},
    }

    for _, test := range tests {
        t.Run(test.key, func(t *testing.T) {
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

func TestCreateKey(t *testing.T) {
    // Create a mock server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/keys" && r.Method == http.MethodPost {
            var opts ApiKeyInput
            if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
                http.Error(w, "bad request", http.StatusBadRequest)
                return
            }
            json.NewEncoder(w).Encode(ApiKey{Key: "new-key"})
        } else {
            http.Error(w, "not found", http.StatusNotFound)
        }
    }))
    defer mockServer.Close()

    client := &Client{BaseURL: mockServer.URL, HTTPClient: mockServer.Client()}
    service := &ApiKeysService{client: client}

    opts := ApiKeyInput{Name: "test-key"}
    apiKey, err := service.CreateKey(context.Background(), opts)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if apiKey.Key != "new-key" {
        t.Errorf("expected new-key, got %v", apiKey.Key)
    }
}