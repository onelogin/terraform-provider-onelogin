package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
)

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

func InflateAppConfiguration(s map[string]interface{}) *models.AppConfiguration {
	rui := s["redirect_uri"].(string)
	rte := int32(s["refresh_token_expiration_minutes"].(int))
	lur := s["login_url"].(string)
	oat := int32(s["oidc_application_type"].(int))
	tea := int32(s["token_endpoint_auth_method"].(int))
	ate := int32(s["access_token_expiration_minutes"].(int))
	par := s["provider_arn"].(string)
	sal := s["signature_algorithm"].(string)

	return &models.AppConfiguration{
		RedirectURI:                   &rui,
		RefreshTokenExpirationMinutes: &rte,
		LoginURL:                      &lur,
		OidcApplicationType:           &oat,
		TokenEndpointAuthMethod:       &tea,
		AccessTokenExpirationMinutes:  &ate,
		ProviderArn:                   &par,
		SignatureAlgorithm:            &sal,
	}
}
