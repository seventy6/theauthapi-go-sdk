package theauthapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrKeyInvalid = errors.New("key is invalid")
)

type ApiKeysService struct {
	client *Client
}

type ApiKeysAuthResponse struct {
	Key             string            `json:"key"`
	Name            string            `json:"name"`
	CustomMetaData  map[string]string `json:"customMetaData"`
	CustomAccountID map[string]string `json:"customAccountId"`
	CustomUserID    string            `json:"customUserId"`
	Env             string            `json:"env"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	Active          bool              `json:"isActive"`
	ExpiresAt       time.Time         `json:"expiresAt"`
	RateLimitConfig interface{}       `json:"rateLimitConfig"`
	CreationContext interface{}       `json:"creationContext"`
}

func (s *ApiKeysService) IsValidKey(ctx context.Context, key string) error {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, fmt.Sprintf(PathApiKeysAuth, ""), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrKeyInvalid
	}

	return nil
}
