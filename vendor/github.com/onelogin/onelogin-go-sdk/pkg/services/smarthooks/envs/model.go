package smarthookenvs

import (
	"time"
)

// SmartHookEnvVarQuery represents available query parameters
type SmartHookEnvVarQuery struct {
	Limit  string
	Page   string
	Cursor string
	Type   string
}

// EnvVar represents an Environment Variable to be associated with a SmartHook
type EnvVar struct {
	ID        *string    `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Value     *string    `json:"value,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
