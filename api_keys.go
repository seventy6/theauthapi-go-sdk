package theauthapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (s *ApiKeysService) IsValidKey(ctx context.Context, key string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api-keys/auth/%s", s.client.BaseURL, key), nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("x-api-key", s.client.AccessToken)

    resp, err := s.client.HTTPClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the body once
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Use the bytes for both JSON decoding and string conversion
    var result struct {
        Valid bool `json:"valid"`
    }
    if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
        return "", err
    }

    bodyString := string(bodyBytes)

    return bodyString, err
}