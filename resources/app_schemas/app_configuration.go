package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// AppConfiguration returns a key/value map of the various fields that make up
// the AppConfiguration field for a OneLogin App.
func AppConfiguration() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"redirect_uri": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"refresh_token_expiration_minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"login_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"oidc_application_type": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"token_endpoint_auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"access_token_expiration_minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"provider_arn": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"signature_algorithm": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// InflateAppConfiguration takes a key/value map of interfaces and uses the fields to construct
// an AppConfiguration struct, a sub-field of a OneLogin App.
func InflateAppConfiguration(s map[string]interface{}) *models.AppConfiguration {
	return &models.AppConfiguration{
		RedirectURI:                   oltypes.String(s["redirect_uri"].(string)),
		RefreshTokenExpirationMinutes: oltypes.Int32(int32(s["refresh_token_expiration_minutes"].(int))),
		LoginURL:                      oltypes.String(s["login_url"].(string)),
		OidcApplicationType:           oltypes.Int32(int32(s["oidc_application_type"].(int))),
		TokenEndpointAuthMethod:       oltypes.Int32(int32(s["token_endpoint_auth_method"].(int))),
		AccessTokenExpirationMinutes:  oltypes.Int32(int32(s["access_token_expiration_minutes"].(int))),
		ProviderArn:                   oltypes.String(s["provider_arn"].(string)),
		SignatureAlgorithm:            oltypes.String(s["signature_algorithm"].(string)),
	}
}
