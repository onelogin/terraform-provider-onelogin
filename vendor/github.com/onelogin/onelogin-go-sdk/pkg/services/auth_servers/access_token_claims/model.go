package accesstokenclaims

type AccessTokenClaimsQuery struct {
	AuthServerID string
}

type AccessTokenClaim struct {
	ID                       *int32   `json:"id,omitempty"`
	AuthServerID             *int32   `json:"auth_server_id,omitempty"`
	Label                    *string  `json:"label,omitempty"`
	UserAttributeMappings    *string  `json:"user_attribute_mappings,omitempty"`
	UserAttributeMacros      *string  `json:"user_attribute_macros,omitempty"`
	AttributeTransformations *string  `json:"attribute_transformations,omitempty"`
	SkipIfBlank              *bool    `json:"skip_if_blank,omitempty"`
	Values                   []string `json:"values,omitempty"`
	DefaultValues            *string  `json:"default_values,omitempty"`
	ProvisionedEntitlements  *bool    `json:"provisioned_entitlements,omitempty"`
}
