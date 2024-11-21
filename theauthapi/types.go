package theauthapi

import "time"

// Common response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
// RateLimitConfiguration defines rate limiting settings
type RateLimitConfiguration struct {
	RateLimit    int `json:"rateLimit"`
	RateLimitTtl int `json:"rateLimitTtl"`
}

// ApiKey represents an API key and its associated data
type ApiKey struct {
	Key              string                 `json:"key"`
	Name             string                 `json:"name"`
	CustomMetaData   map[string]interface{} `json:"customMetaData"`
	CustomAccountID  string                 `json:"customAccountId"`
	CustomUserID     string                 `json:"customUserId"`
	Env              Environment            `json:"env"`
	CreatedAt        time.Time             `json:"createdAt"`
	UpdatedAt        time.Time             `json:"updatedAt"`
	IsActive         bool                   `json:"isActive"`
	RateLimitConfigs RateLimitConfiguration `json:"rateLimitConfigs"`
	Expiry           time.Time             `json:"expiry"`
}

// ApiKeyInput represents the input for creating a new API key
type ApiKeyInput struct {
	Name             string                  `json:"name"`
	ProjectID        *string                 `json:"projectId,omitempty"`
	Key              *string                 `json:"key,omitempty"`
	CustomMetaData   map[string]interface{}  `json:"customMetaData,omitempty"`
	CustomAccountID  *string                 `json:"customAccountId,omitempty"`
	CustomUserID     *string                 `json:"customUserId,omitempty"`
	RateLimitConfigs *RateLimitConfiguration `json:"rateLimitConfigs,omitempty"`
	Expiry           *time.Time             `json:"expiry,omitempty"`
}

// ApiKeyFilter represents filtering options for API keys
type ApiKeyFilter struct {
	ProjectID      *string `json:"projectId,omitempty"`
	Name           *string `json:"name,omitempty"`
	CustomAccountID *string `json:"customAccountId,omitempty"`
	CustomUserID    *string `json:"customUserId,omitempty"`
	IsActive       *bool   `json:"isActive,omitempty"`
}

// UpdateApiKeyInput represents the input for updating an API key
type UpdateApiKeyInput struct {
	Name             string                  `json:"name"`
	Key              *string                 `json:"key,omitempty"`
	CustomMetaData   map[string]interface{}  `json:"customMetaData,omitempty"`
	CustomAccountID  *string                 `json:"customAccountId,omitempty"`
	CustomUserID     *string                 `json:"customUserId,omitempty"`
	Expiry           *time.Time             `json:"expiry,omitempty"`
	RateLimitConfigs *RateLimitConfiguration `json:"rateLimitConfigs,omitempty"`
}

// AuthedEntityType represents the type of authenticated entity
type AuthedEntityType string

const (
	AuthedEntityTypeUser      AuthedEntityType = "USER"
	AuthedEntityTypeAccessKey AuthedEntityType = "ACCESS_KEY"
)

// AuthBaseEntity contains common fields for authenticated entities
type AuthBaseEntity struct {
	IsActive           bool              `json:"isActive"`
	CreatedBy         string            `json:"createdBy"`
	CreatedByType     *AuthedEntityType `json:"createdByType,omitempty"`
	CreatedIn         string            `json:"createdIn"`
	LastChangedBy     string            `json:"lastChangedBy"`
	LastChangedByType *AuthedEntityType `json:"lastChangedByType,omitempty"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	CreatedAt         time.Time         `json:"createdAt"`
}

// Project represents a project entity
type Project struct {
	AuthBaseEntity
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	AccountID string      `json:"accountId"`
	Env       Environment `json:"env"`
}

// Environment represents the environment type
type Environment string

const (
	EnvironmentLive Environment = "live"
	EnvironmentTest Environment = "test"
)

// CreateProjectInput represents the input for creating a new project
type CreateProjectInput struct {
	Name      string      `json:"name"`
	AccountID string      `json:"accountId"`
	Env       Environment `json:"env"`
}

// UpdateProjectInput represents the input for updating a project
type UpdateProjectInput struct {
	Name string `json:"name"`
}

// Account represents an account entity
type Account struct {
	AuthBaseEntity
	ID   string `json:"id"`
	Name string `json:"name"`
}
