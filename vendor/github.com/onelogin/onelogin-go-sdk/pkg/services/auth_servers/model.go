package authservers

type AuthServerQuery struct {
	Name   string `json:"name,omitempty"`
	Limit  string
	Page   string
	Cursor string
}

type AuthServer struct {
	ID            *int32                   `json:"id,omitempty"`
	Name          *string                  `json:"name,omitempty"`
	Description   *string                  `json:"description,omitempty"`
	Configuration *AuthServerConfiguration `json:"configuration,omitempty"`
}

type AuthServerConfiguration struct {
	ResourceIdentifier            *string  `json:"resource_identifier,omitempty"`
	Audiences                     []string `json:"audiences,omitempty"`
	AccessTokenExpirationMinutes  *int32   `json:"access_token_expiration_minutes,omitempty"`
	RefreshTokenExpirationMinutes *int32   `json:"refresh_token_expiration_minutes,omitempty"`
}
