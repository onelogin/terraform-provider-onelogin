package olhttp

type OLHTTPRequest struct {
	URL        string
	AuthMethod string
	Headers    map[string]string
	Payload    interface{}
}

// AuthBody is the request payload for authorization.
type AuthBody struct {
	GrantType string `json:"grant_type"`
}

// ClientCredential is the authorization response payload.
type ClientCredential struct {
	AccessToken  *string `json:"access_token,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	ExpiresIn    *int32  `json:"expires_in,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	TokenType    *string `json:"token_type,omitempty"`
	AccountID    *int32  `json:"account_id,omitempty"`
}
