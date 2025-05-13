package models

type App struct {
	ID                 *int32                `json:"id,omitempty"`
	ConnectorID        *int32                `json:"connector_id"`
	Name               *string               `json:"name"`
	Description        *string               `json:"description,omitempty"`
	Notes              *string               `json:"notes,omitempty"`
	PolicyID           *int                  `json:"policy_id,omitempty"`
	BrandID            *int                  `json:"brand_id,omitempty"`
	IconURL            *string               `json:"icon_url,omitempty"`
	Visible            *bool                 `json:"visible,omitempty"`
	AuthMethod         *int                  `json:"auth_method,omitempty"`
	TabID              *int                  `json:"tab_id,omitempty"`
	CreatedAt          *string               `json:"created_at,omitempty"`
	UpdatedAt          *string               `json:"updated_at,omitempty"`
	RoleIDs            *[]int                `json:"role_ids,omitempty"`
	AllowAssumedSignin *bool                 `json:"allow_assumed_signin,omitempty"`
	Provisioning       *Provisioning         `json:"provisioning,omitempty"`
	SSO                interface{}           `json:"sso,omitempty"`
	Configuration      interface{}           `json:"configuration,omitempty"`
	Parameters         *map[string]Parameter `json:"parameters,omitempty"`
	EnforcementPoint   *EnforcementPoint     `json:"enforcement_point,omitempty"`
}

type Provisioning struct {
	Enabled bool `json:"enabled"`
}

type SSO interface {
	ValidateSSO() error
}

type SSOOpenId struct {
	ClientID string `json:"client_id"`
}

type SSOSAML struct {
	MetadataURL string      `json:"metadata_url"`
	AcsURL      string      `json:"acs_url"`
	SlsURL      string      `json:"sls_url"`
	Issuer      string      `json:"issuer"`
	Certificate Certificate `json:"certificate"`
}

type Certificate struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ConfigurationOpenId struct {
	RedirectURI                   string `json:"redirect_uri"`
	LoginURL                      string `json:"login_url"`
	OidcApplicationType           int    `json:"oidc_application_type"`
	TokenEndpointAuthMethod       int    `json:"token_endpoint_auth_method"`
	AccessTokenExpirationMinutes  int    `json:"access_token_expiration_minutes"`
	RefreshTokenExpirationMinutes int    `json:"refresh_token_expiration_minutes"`
}

type ConfigurationSAML struct {
	ProviderArn        interface{} `json:"provider_arn"`
	SignatureAlgorithm string      `json:"signature_algorithm"`
	CertificateID      int         `json:"certificate_id"`
}

type Parameter struct {
	Values                    interface{} `json:"values,omitempty"`
	UserAttributeMappings     interface{} `json:"user_attribute_mappings,omitempty"`
	ProvisionedEntitlements   bool        `json:"provisioned_entitlements,omitempty"`
	SkipIfBlank               bool        `json:"skip_if_blank,omitempty"`
	ID                        int         `json:"id,omitempty"`
	DefaultValues             interface{} `json:"default_values"`
	AttributesTransformations interface{} `json:"attributes_transformations,omitempty"`
	Label                     string      `json:"label,omitempty"`
	UserAttributeMacros       interface{} `json:"user_attribute_macros,omitempty"`
	IncludeInSamlAssertion    bool        `json:"include_in_saml_assertion,omitempty"`
}

type EnforcementPoint struct {
	RequireSitewideAuthentication bool        `json:"require_sitewide_authentication"`
	Conditions                    *Conditions `json:"conditions,omitempty"`
	SessionExpiryFixed            Duration    `json:"session_expiry_fixed"`
	SessionExpiryInactivity       Duration    `json:"session_expiry_inactivity"`
	Permissions                   string      `json:"permissions"`
	Token                         string      `json:"token,omitempty"`
	Target                        string      `json:"target"`
	Resources                     []Resource  `json:"resources"`
	ContextRoot                   string      `json:"context_root"`
	UseTargetHostHeader           bool        `json:"use_target_host_header"`
	Vhost                         string      `json:"vhost"`
	LandingPage                   string      `json:"landing_page"`
	CaseSensitive                 bool        `json:"case_sensitive"`
}

type Conditions struct {
	Type  string   `json:"type"`
	Roles []string `json:"roles"`
}

type Duration struct {
	Value int `json:"value"`
	Unit  int `json:"unit"`
}

type Resource struct {
	Path        string  `json:"path"`
	RequireAuth string  `json:"require_authentication"`
	Permissions string  `json:"permissions"`
	Conditions  *string `json:"conditions,omitempty"`
	IsPathRegex *bool   `json:"is_path_regex,omitempty"`
	ResourceID  int     `json:"resource_id,omitempty"`
}

const (
	UnitSeconds = 0
	UnitMinutes = 1
	UnitHours   = 2
)

type AppQuery struct {
	Limit       string  `json:"limit,omitempty"`
	Page        string  `json:"page,omitempty"`
	Cursor      string  `json:"cursor,omitempty"`
	Name        *string `json:"name,omitempty"`
	ConnectorID *int    `json:"connector_id,omitempty"`
	AuthMethod  *int    `json:"auth_method,omitempty"`
}

func (q *AppQuery) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":        validateString,
		"page":         validateString,
		"cursor":       validateString,
		"name":         validateString,
		"connector_id": validateInt,
		"auth_method":  validateInt,
	}
}
