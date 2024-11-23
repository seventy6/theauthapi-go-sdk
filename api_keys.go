package theauthapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ErrKeyInvalid = errors.New("key is invalid")
)

type ApiKeysService struct {
	client *Client
	debug  bool
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
	resp, err := s.client.sendRequest(ctx, http.MethodGet, fmt.Sprintf(PathApiKeysAuth, key), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if s.debug {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading body: %v \n", err)
		}
		log.Println("body: ", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return ErrKeyInvalid
	}

	return nil
}
