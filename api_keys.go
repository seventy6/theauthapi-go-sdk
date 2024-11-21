package theauthapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(".env"); err != nil {
        panic("Error 1 loading .env file")
    }
}
type ApiKeysService struct {
    client *Client
}

func (s *ApiKeysService) IsValidKey(ctx context.Context, key string) (bool, error) {
    accessToken := os.Getenv("ACCESS_TOKEN")
    if accessToken == "" {
        return false, fmt.Errorf("access token is not set")
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api-keys/auth/%s", s.client.BaseURL, key), nil)
    if err != nil {
        return false, err
    }
    req.Header.Set("x-api-key", accessToken)

    resp, err := s.client.HTTPClient.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var result struct {
        Valid bool `json:"valid"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false, err
    }

    return result.Valid, nil
}

func (s *ApiKeysService) CreateKey(ctx context.Context, opts ApiKeyInput) (*ApiKey, error) {
    accessToken := os.Getenv("ACCESS_TOKEN")
    if accessToken == "" {
        return nil, fmt.Errorf("access token is not set")
    }

    body, err := json.Marshal(opts)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api-keys/", s.client.BaseURL), bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", accessToken)

    resp, err := s.client.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var apiKey ApiKey
    if err := json.NewDecoder(resp.Body).Decode(&apiKey); err != nil {
        return nil, err
    }

    return &apiKey, nil
}