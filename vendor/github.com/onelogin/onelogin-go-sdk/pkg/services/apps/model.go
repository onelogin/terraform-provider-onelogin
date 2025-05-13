package apps

import "time"

const (
	AuthMethodPassword int32 = iota
	AuthMethodOpenID
	AuthMethodSAML
	AuthMethodAPI
	AuthMethodGoogle
	authMethodUnused5 // There is not auth method with the number 5
	AuthMethodForemsBasedApp
	AuthMethodWSFED
	AuthMethodOpenIDConnect
)

type AppsQuery struct {
	Limit       string
	Page        string
	Name        string
	ConnectorID string
	AuthMethod  string
	Cursor      string
}

// App is the contract for apps api v2.
type App struct {
	ID                 *int32                   `json:"id,omitempty"`
	Name               *string                  `json:"name,omitempty"`
	Visible            *bool                    `json:"visible,omitempty"`
	Description        *string                  `json:"description,omitempty"`
	Notes              *string                  `json:"notes,omitempty"`
	IconURL            *string                  `json:"icon_url,omitempty"`
	AuthMethod         *int32                   `json:"auth_method,omitempty"`
	PolicyID           *int32                   `json:"policy_id,omitempty"`
	AllowAssumedSignin *bool                    `json:"allow_assumed_signin,omitempty"`
	TabID              *int32                   `json:"tab_id,omitempty"`
	BrandID            *int32                   `json:"brand_id,omitempty"`
	ConnectorID        *int32                   `json:"connector_id,omitempty"`
	CreatedAt          *time.Time               `json:"created_at,omitempty"`
	UpdatedAt          *time.Time               `json:"updated_at,omitempty"`
	Provisioning       *AppProvisioning         `json:"provisioning"`
	Sso                *AppSso                  `json:"sso"`
	Configuration      *AppConfiguration        `json:"configuration"`
	Parameters         map[string]AppParameters `json:"parameters"`
	RoleIDs            []int                    `json:"role_ids"`
}

// AppProvisioning is the contract for provisioning.
type AppProvisioning struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// AppSso is the contract for apps sso.
type AppSso struct {
	ClientID     *string            `json:"client_id,omitempty"`
	ClientSecret *string            `json:"client_secret,omitempty"`
	MetadataURL  *string            `json:"metadata_url,omitempty"`
	AcsURL       *string            `json:"acs_url,omitempty"`
	SlsURL       *string            `json:"sls_url,omitempty"`
	Issuer       *string            `json:"issuer,omitempty"`
	Certificate  *AppSsoCertificate `json:"certificate"`
}

// AppSsoCertificate is the contract for sso certificate.
type AppSsoCertificate struct {
	ID    *int32  `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

// AppConfiguration is the contract for configuration.
type AppConfiguration struct {
	RedirectURI                   *string `json:"redirect_uri,omitempty"`
	RefreshTokenExpirationMinutes *int32  `json:"refresh_token_expiration_minutes,omitempty"`
	LoginURL                      *string `json:"login_url,omitempty"`
	OidcApplicationType           *int32  `json:"oidc_application_type,omitempty"`
	TokenEndpointAuthMethod       *int32  `json:"token_endpoint_auth_method,omitempty"`
	AccessTokenExpirationMinutes  *int32  `json:"access_token_expiration_minutes,omitempty"`
	ProviderArn                   *string `json:"provider_arn,omitempty"`
	IdpList                       *string `json:"idp_list,omitempty"`
	SignatureAlgorithm            *string `json:"signature_algorithm,omitempty"`

	LogoutURL                    *string `json:"logout_url,omitempty"`
	PostLogoutRedirectURI        *string `json:"post_logout_redirect_uri,omitempty"`
	Audience                     *string `json:"audience,omitempty"`
	ConsumerURL                  *string `json:"consumer_url,omitempty"`
	Login                        *string `json:"login,omitempty"`
	Recipient                    *string `json:"recipient,omitempty"`
	Validator                    *string `json:"validator,omitempty"`
	RelayState                   *string `json:"relaystate,omitempty"`
	Relay                        *string `json:"relay,omitempty"`
	SAMLNotValidOnOrAafter       *string `json:"saml_notonorafter,omitempty"`
	GenerateAttributeValueTags   *string `json:"generate_attribute_value_tags,omitempty"`
	SAMLInitiaterID              *string `json:"saml_initiater_id,omitempty"`
	SAMLNotValidBefore           *string `json:"saml_notbefore,omitempty"`
	SAMLIssuerType               *string `json:"saml_issuer_type,omitempty"`
	SAMLSignElement              *string `json:"saml_sign_element,omitempty"`
	EncryptAssertion             *string `json:"encrypt_assertion,omitempty"`
	SAMLSessionNotValidOnOrAfter *string `json:"saml_sessionnotonorafter,omitempty"`
	SAMLEncryptionMethodID       *string `json:"saml_encryption_method_id,omitempty"`
	SAMLNameIDFormatID           *string `json:"saml_nameid_format_id,omitempty"`
}

// AppParameters is the contract for parameters.
type AppParameters struct {
	ID                        *int32  `json:"id,omitempty"`
	Label                     *string `json:"label,omitempty"`
	UserAttributeMappings     *string `json:"user_attribute_mappings,omitempty"`
	UserAttributeMacros       *string `json:"user_attribute_macros,omitempty"`
	AttributesTransformations *string `json:"attributes_transformations,omitempty"`
	SkipIfBlank               *bool   `json:"skip_if_blank,omitempty"`
	Values                    *string `json:"values,omitempty,omitempty"`
	DefaultValues             *string `json:"default_values,omitempty"`
	ParamKeyName              *string `json:"param_key_name,omitempty"`
	ProvisionedEntitlements   *bool   `json:"provisioned_entitlements,omitempty"`
	SafeEntitlementsEnabled   *bool   `json:"safe_entitlements_enabled,omitempty"`
	IncludeInSamlAssertion    *bool   `json:"include_in_saml_assertion,omitempty"`
}

// AppUser is the contract for users of an app.
type AppUser struct {
	ID        *int32  `json:"id,omitempty"`
	Firstname *string `json:"firstname,omitempty"`
	Lastname  *string `json:"lastname,omitempty"`
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
}
