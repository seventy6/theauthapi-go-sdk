// theauthapi/client.go
package theauthapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the main TheAuthAPI client
type Client struct {
	BaseURL    string
	AccessKey  string
	HTTPClient *http.Client
	ApiKeys    *ApiKeysService
	Projects   *ProjectsService
	Accounts   *AccountsService
}

// ClientOption allows customizing the client
type ClientOption func(*Client)

// NewClient creates a new TheAuthAPI client with optional configurations
func NewClient(accessKey string, opts ...ClientOption) *Client {
	// Default configuration
	client := &Client{
		BaseURL: "https://api.theauthapi.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		AccessKey: accessKey,
	}

	// Apply custom options
	for _, opt := range opts {
		opt(client)
	}

	// Initialize services
	client.ApiKeys = &ApiKeysService{client: client}
	client.Projects = &ProjectsService{client: client}
	client.Accounts = &AccountsService{client: client}

	return client
}

// WithBaseURL allows overriding the base API URL (useful for testing)
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = url
	}
}

// WithHTTPClient allows providing a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// sendRequest is a helper method to send HTTP requests
func (c *Client) sendRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Serialize request body
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s%s", c.BaseURL, path),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessKey))

	// Send request
	return c.HTTPClient.Do(req)
}