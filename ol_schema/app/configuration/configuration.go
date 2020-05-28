package configuration

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/onelogin-terraform-provider/utils"
)

// OIDCSchema returns a key/value map of the various fields that make up
// the Configuration field for a OneLogin OIDC App.
func OIDCSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"redirect_uri": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"refresh_token_expiration_minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"login_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"oidc_application_type": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"token_endpoint_auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"access_token_expiration_minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
	}
}

// SAMLConfigurationSchema returns a key/value map of the various fields that make up
// the Configuration field for a OneLogin SAML App.
func SAMLSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"provider_arn": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"signature_algorithm": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validSignatureAlgorithm,
		},
	}
}

func validSignatureAlgorithm(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"})
}

// Inflate takes a map of interfaces and uses the fields to construct
// an AppConfiguration instance.
func Inflate(s map[string]interface{}) apps.AppConfiguration {
	out := apps.AppConfiguration{}
	var st string
	var n int
	var notNil bool

	// oidc fields
	if st, notNil = s["redirect_uri"].(string); notNil {
		out.RedirectURI = oltypes.String(st)
	}
	if st, notNil = s["login_url"].(string); notNil {
		out.LoginURL = oltypes.String(st)
	}

	if n, notNil = s["refresh_token_expiration_minutes"].(int); n != 0 && notNil {
		out.RefreshTokenExpirationMinutes = oltypes.Int32(int32(n))
	}

	if n, notNil = s["oidc_application_type"].(int); n != 0 && notNil {
		out.OidcApplicationType = oltypes.Int32(int32(n))
	}

	if n, notNil = s["token_endpoint_auth_method"].(int); n != 0 && notNil {
		out.TokenEndpointAuthMethod = oltypes.Int32(int32(n))
	}

	if n, notNil = s["access_token_expiration_minutes"].(int); n != 0 && notNil {
		out.AccessTokenExpirationMinutes = oltypes.Int32(int32(n))
	}

	// saml fields
	if st, notNil = s["provider_arn"].(string); notNil {
		out.ProviderArn = oltypes.String(st)
	}
	if st, notNil = s["signature_algorithm"].(string); notNil {
		out.SignatureAlgorithm = oltypes.String(st)
	}
	return out
}

// FlattenOIDC takes an instance of AppConfiguration and return an array of
// maps. Fields differ depending on if the app is a SAML or OIDC app.
func FlattenOIDC(config apps.AppConfiguration) []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"redirect_uri":                     config.RedirectURI,
			"login_url":                        config.LoginURL,
			"refresh_token_expiration_minutes": config.RefreshTokenExpirationMinutes,
			"oidc_application_type":            config.OidcApplicationType,
			"token_endpoint_auth_method":       config.TokenEndpointAuthMethod,
			"access_token_expiration_minutes":  config.AccessTokenExpirationMinutes,
		},
	}
}

// FlattenSAML takes an instance of AppConfiguration and return an array of
// maps. Fields differ depending on if the app is a SAML or OIDC app.
func FlattenSAML(config apps.AppConfiguration) []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"provider_arn":        config.ProviderArn,
			"signature_algorithm": config.SignatureAlgorithm,
		},
	}
}
