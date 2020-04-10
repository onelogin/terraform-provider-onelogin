package app

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// AppConfiguration returns a key/value map of the various fields that make up
// the AppConfiguration field for a OneLogin App.
func ConfigurationSchema() map[string]*schema.Schema {
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
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				validOpts := []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"}
				v := val.(string)
				isValid := false
				for _, o := range validOpts {
					isValid = v == o
					if isValid {
						break
					}
				}
				if !isValid {
					errs = append(errs, fmt.Errorf("signature_algorithm must be one of %v, got: %s", validOpts, v))
				}
				return
			},
		},
	}
}

// InflateAppConfiguration takes a key/value map of interfaces and uses the fields to construct
// an AppConfiguration struct, a sub-field of a OneLogin App.
func InflateConfiguration(s *map[string]interface{}) *models.AppConfiguration {
	out := models.AppConfiguration{}
	var st string
	var n int
	var notNil bool

	if st, notNil = (*s)["redirect_uri"].(string); notNil {
		out.RedirectURI = oltypes.String(st)
	}
	if st, notNil = (*s)["login_url"].(string); notNil {
		out.LoginURL = oltypes.String(st)
	}
	if st, notNil = (*s)["provider_arn"].(string); notNil {
		out.ProviderArn = oltypes.String(st)
	}
	if st, notNil = (*s)["signature_algorithm"].(string); notNil {
		out.SignatureAlgorithm = oltypes.String(st)
	}

	if n, notNil = (*s)["refresh_token_expiration_minutes"].(int); notNil {
		out.RefreshTokenExpirationMinutes = oltypes.Int32(int32(n))
	}
	if n, notNil = (*s)["oidc_application_type"].(int); notNil {
		out.OidcApplicationType = oltypes.Int32(int32(n))
	}
	if n, notNil = (*s)["token_endpoint_auth_method"].(int); notNil {
		out.TokenEndpointAuthMethod = oltypes.Int32(int32(n))
	}
	if n, notNil = (*s)["access_token_expiration_minutes"].(int); notNil {
		out.AccessTokenExpirationMinutes = oltypes.Int32(int32(n))
	}
	return &out
}
