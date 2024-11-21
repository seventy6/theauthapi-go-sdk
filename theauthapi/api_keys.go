package theauthapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiKeysService struct {
    client *Client
}


func (s *ApiKeysService) IsValidKey(ctx context.Context, key string) (bool, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/validate?key=%s", s.client.BaseURL, key), nil)
    if err != nil {
        return false, err
    }

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
    body, err := json.Marshal(opts)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/keys", s.client.BaseURL), bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

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