package scopes

type ScopesQuery struct {
	AuthServerID string
}

type Scope struct {
	ID           *int32  `json:"id,omitempty"`
	AuthServerID *int32  `json:"auth_server_id,omitempty"`
	Value        *string `json:"value,omitempty"`
	Description  *string `json:"description,omitempty"`
}
