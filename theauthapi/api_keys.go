package theauthapi

import (
    "context"
)

type ApiKeysService struct {
    client *Client
}

func (s *ApiKeysService) IsValidKey(ctx context.Context, key string) (bool, error) {
    // Implementation
}

func (s *ApiKeysService) CreateKey(ctx context.Context, opts CreateKeyOptions) (*ApiKey, error) {
    // Implementation
}

// Add other API key related methods