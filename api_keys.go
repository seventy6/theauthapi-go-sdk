package theauthapi

import (
	"bytes"
	"context"
	"encoding/json"
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

type HTTPResponse struct {
	StatusCode int
}

type ApiKeysAuthResponse struct {
	Key             string          `json:"key"`
	Name            string          `json:"name"`
	CustomMetaData  json.RawMessage `json:"customMetaData"`
	CustomAccountID json.RawMessage `json:"customAccountId"`
	CustomUserID    string          `json:"customUserId"`
	Env             string          `json:"env"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Active          bool            `json:"isActive"`
	ExpiresAt       time.Time       `json:"expiresAt"`
	RateLimitConfig json.RawMessage `json:"rateLimitConfig"`
	CreationContext json.RawMessage `json:"creationContext"`

	HTTPResponse HTTPResponse
}

func (s *ApiKeysService) GetValidKey(ctx context.Context, key string) (*ApiKeysAuthResponse, error) {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, fmt.Sprintf(PathApiKeysAuth, key), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))

	fakebody := bytes.NewBuffer(body)

	log.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode, http.StatusOK)
		return nil, ErrKeyInvalid
	}

	apiResp := &ApiKeysAuthResponse{}
	err = json.NewDecoder(fakebody).Decode(apiResp)
	if err != nil {
		return nil, err
	}

	apiResp.HTTPResponse = HTTPResponse{
		StatusCode: resp.StatusCode,
	}

	return apiResp, nil

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
