package models

import (
	"time"
)

const (
	TypePreAuthentication string = "pre-authentication"
	TypeUserMigration     string = "user-migration"
)

const (
	ContextPreAuthentication1_0_0 string = "1.0.0"
	ContextPreAuthentication1_1_0 string = "1.1.0"

	ContextUserMigration1_0_0 string = "1.0.0"
)

const (
	StatusReady         string = "ready"
	StatusCreateQueued  string = "create-queued"
	StatusCreateRunning string = "create-running"
	StatusCreateFailed  string = "create-failed"
	StatusUpdateQueued  string = "update-queued"
	StatusUpdateRunning string = "update-running"
	StatusUpdateFailed  string = "update-failed"
)

// SmartHookQuery represents available query parameters
type SmartHookQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Type   string `json:"type,omitempty"`
}

// SmartHook represents a OneLogin SmartHook with associated resource data
type SmartHook struct {
	ID             *string           `json:"id,omitempty"`
	Type           *string           `json:"type,omitempty"`
	Disabled       *bool             `json:"disabled,omitempty"`
	Timeout        *int32            `json:"timeout,omitempty"`
	EnvVars        []EnvVar          `json:"env_vars"`
	Runtime        *string           `json:"runtime,omitempty"`
	ContextVersion *string           `json:"context_version,omitempty"`
	Retries        *int32            `json:"retries,omitempty"`
	Options        *Options          `json:"options,omitempty"`
	Packages       map[string]string `json:"packages"`
	Function       *string           `json:"function,omitempty"`
	Status         *string           `json:"status,omitempty"`
	CreatedAt      *time.Time        `json:"created_at,omitempty"`
	UpdatedAt      *time.Time        `json:"updated_at,omitempty"`
	Conditions     []Condition       `json:"conditions,omitempty"`
}

// SmartHookOptions represents the options to be associated with a SmartHook
type Options struct {
	RiskEnabled          *bool `json:"risk_enabled,omitempty"`
	MFADeviceInfoEnabled *bool `json:"mfa_device_info_enabled,omitempty"`
	LocationEnabled      *bool `json:"location_enabled,omitempty"`
}

// SmartHookEnvVarQuery represents available query parameters
type SmartHookEnvVarQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Type   string `json:"type,omitempty"`
}

// EnvVar represents an Environment Variable to be associated with a SmartHook
type EnvVar struct {
	ID        *string    `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Value     *string    `json:"value,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (s *SmartHook) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":  validateString,
		"page":   validateString,
		"cursor": validateString,
		"type":   validateString,
	}
}
