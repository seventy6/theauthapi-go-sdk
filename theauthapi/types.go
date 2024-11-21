package theauthapi

import "time"

// Common response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ApiKey represents an API key in the system
type ApiKey struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key,omitempty"`      // Only present when creating new keys
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project represents a project in the system
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Account represents a user account
type Account struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request types
type CreateApiKeyRequest struct {
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// List response types
type ListApiKeysResponse struct {
	Response
	Data []ApiKey `json:"data"`
}

type ListProjectsResponse struct {
	Response
	Data []Project `json:"data"`
}

type ListAccountsResponse struct {
	Response
	Data []Account `json:"data"`
}