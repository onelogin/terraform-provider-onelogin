package appconfigurationschema

import (
	"errors"
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestInflateConfiguration(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput interface{}
		ExpectedError  error
	}{
		"creates and returns the address of an AppConfiguration struct for a OIDC app": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": "2",
				"login_url":                        "test",
				"oidc_application_type":            "2",
				"token_endpoint_auth_method":       "2",
				"access_token_expiration_minutes":  "2",
			},
			ExpectedOutput: CustomConfigurationOpenId{
				RedirectURI:                   "test",
				RefreshTokenExpirationMinutes: &[]int{2}[0],
				LoginURL:                      "test",
				OidcApplicationType:           2,
				TokenEndpointAuthMethod:       2,
				AccessTokenExpirationMinutes:  &[]int{2}[0],
			},
		},
		"handles OIDC app config with only redirect_uri (no timeout fields)": {
			ResourceData: map[string]interface{}{
				"redirect_uri":               "https://example.com/callback",
				"oidc_application_type":      "0",
				"token_endpoint_auth_method": "1",
			},
			ExpectedOutput: CustomConfigurationOpenId{
				RedirectURI:             "https://example.com/callback",
				OidcApplicationType:     0,
				TokenEndpointAuthMethod: 1,
				// timeout fields should be nil (not set)
				RefreshTokenExpirationMinutes: nil,
				AccessTokenExpirationMinutes:  nil,
			},
		},
		"handles OIDC app config with empty string timeout fields": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "https://example.com/callback",
				"oidc_application_type":            "0",
				"token_endpoint_auth_method":       "1",
				"refresh_token_expiration_minutes": "",
				"access_token_expiration_minutes":  "",
			},
			ExpectedOutput: CustomConfigurationOpenId{
				RedirectURI:             "https://example.com/callback",
				OidcApplicationType:     0,
				TokenEndpointAuthMethod: 1,
				// timeout fields should be nil when empty strings
				RefreshTokenExpirationMinutes: nil,
				AccessTokenExpirationMinutes:  nil,
			},
		},
		"returns an error if invalid refresh_token_expiration_minutes given": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": "asdf",
				"login_url":                        "test",
				"oidc_application_type":            "2",
				"token_endpoint_auth_method":       "2",
				"access_token_expiration_minutes":  "2",
			},
			ExpectedError: errors.New(`strconv.Atoi: parsing "asdf": invalid syntax`),
		},
		"returns an error if invalid oidc_application_type given": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": "2",
				"login_url":                        "test",
				"oidc_application_type":            "asdf",
				"token_endpoint_auth_method":       "2",
				"access_token_expiration_minutes":  "2",
			},
			ExpectedError: errors.New(`strconv.Atoi: parsing "asdf": invalid syntax`),
		},
		"returns an error if invalid token_endpoint_auth_method given": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": "2",
				"login_url":                        "test",
				"oidc_application_type":            "2",
				"token_endpoint_auth_method":       "asdf",
				"access_token_expiration_minutes":  "2",
			},
			ExpectedError: errors.New(`strconv.Atoi: parsing "asdf": invalid syntax`),
		},
		"returns an error if invalid access_token_expiration_minutes given": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": "2",
				"login_url":                        "test",
				"oidc_application_type":            "2",
				"token_endpoint_auth_method":       "2",
				"access_token_expiration_minutes":  "asdf",
			},
			ExpectedError: errors.New(`strconv.Atoi: parsing "asdf": invalid syntax`),
		},
		"creates and returns the address of an AppConfiguration struct for a SAML app": {
			ResourceData: map[string]interface{}{
				"provider_arn":        "test",
				"signature_algorithm": "test",
				"idp_list":            "test",
			},
			ExpectedOutput: models.ConfigurationSAML{
				ProviderArn:        "test",
				SignatureAlgorithm: "test",
			},
		},
		"creates and returns the address of an AppConfiguration struct for a SAML app with exra fields": {
			ResourceData: map[string]interface{}{
				"provider_arn":        "test",
				"signature_algorithm": "test",
				"idp_list":            "test",
				"encrypt_assertion":   "1",
			},
			ExpectedOutput: models.ConfigurationSAML{
				ProviderArn:        "test",
				SignatureAlgorithm: "test",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, err := Inflate(test.ResourceData)
			if test.ExpectedOutput != nil {
				if oidcConfig, ok := test.ExpectedOutput.(models.ConfigurationOpenId); ok {
					if oidcResult, ok := subj.(models.ConfigurationOpenId); ok {
						assert.Equal(t, oidcConfig.RedirectURI, oidcResult.RedirectURI)
						assert.Equal(t, oidcConfig.LoginURL, oidcResult.LoginURL)
						assert.Equal(t, oidcConfig.RefreshTokenExpirationMinutes, oidcResult.RefreshTokenExpirationMinutes)
						assert.Equal(t, oidcConfig.OidcApplicationType, oidcResult.OidcApplicationType)
						assert.Equal(t, oidcConfig.TokenEndpointAuthMethod, oidcResult.TokenEndpointAuthMethod)
						assert.Equal(t, oidcConfig.AccessTokenExpirationMinutes, oidcResult.AccessTokenExpirationMinutes)
					} else {
						t.Errorf("Expected ConfigurationOpenId but got different type: %T", subj)
					}
				} else if customOidcConfig, ok := test.ExpectedOutput.(CustomConfigurationOpenId); ok {
					if customOidcResult, ok := subj.(CustomConfigurationOpenId); ok {
						assert.Equal(t, customOidcConfig.RedirectURI, customOidcResult.RedirectURI)
						assert.Equal(t, customOidcConfig.LoginURL, customOidcResult.LoginURL)
						assert.Equal(t, customOidcConfig.OidcApplicationType, customOidcResult.OidcApplicationType)
						assert.Equal(t, customOidcConfig.TokenEndpointAuthMethod, customOidcResult.TokenEndpointAuthMethod)

						// Handle pointer comparisons for timeout fields
						if customOidcConfig.RefreshTokenExpirationMinutes == nil {
							assert.Nil(t, customOidcResult.RefreshTokenExpirationMinutes)
						} else {
							assert.NotNil(t, customOidcResult.RefreshTokenExpirationMinutes)
							assert.Equal(t, *customOidcConfig.RefreshTokenExpirationMinutes, *customOidcResult.RefreshTokenExpirationMinutes)
						}

						if customOidcConfig.AccessTokenExpirationMinutes == nil {
							assert.Nil(t, customOidcResult.AccessTokenExpirationMinutes)
						} else {
							assert.NotNil(t, customOidcResult.AccessTokenExpirationMinutes)
							assert.Equal(t, *customOidcConfig.AccessTokenExpirationMinutes, *customOidcResult.AccessTokenExpirationMinutes)
						}
					} else {
						t.Errorf("Expected CustomConfigurationOpenId but got different type: %T", subj)
					}
				} else if samlConfig, ok := test.ExpectedOutput.(models.ConfigurationSAML); ok {
					if samlResult, ok := subj.(models.ConfigurationSAML); ok {
						assert.Equal(t, samlConfig.ProviderArn, samlResult.ProviderArn)
						assert.Equal(t, samlConfig.SignatureAlgorithm, samlResult.SignatureAlgorithm)
					} else {
						t.Errorf("Expected ConfigurationSAML but got different type")
					}
				}
			}
			if test.ExpectedError != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestFlattenConfiguration(t *testing.T) {
	tests := map[string]struct {
		InputData      models.ConfigurationOpenId
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns the address of an AppConfiguration struct for a OIDC app": {
			InputData: models.ConfigurationOpenId{
				RedirectURI:                   "test",
				RefreshTokenExpirationMinutes: 2,
				LoginURL:                      "test",
				OidcApplicationType:           2,
				TokenEndpointAuthMethod:       2,
				AccessTokenExpirationMinutes:  2,
			},
			ExpectedOutput: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": "2",
				"login_url":                        "test",
				"oidc_application_type":            "2",
				"token_endpoint_auth_method":       "2",
				"access_token_expiration_minutes":  "2",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenOIDC(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAMLConfiguration(t *testing.T) {
	tests := map[string]struct {
		InputData      models.ConfigurationSAML
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns the address of an AppConfiguration struct for a SAML app": {
			InputData: models.ConfigurationSAML{
				ProviderArn:        "test",
				SignatureAlgorithm: "test",
			},
			ExpectedOutput: map[string]interface{}{
				"provider_arn":        "test",
				"signature_algorithm": "test",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAML(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestValidSignatureAlgorithm(t *testing.T) {
	tests := map[string]struct {
		InputKey       string
		InputValue     string
		ExpectedOutput []error
	}{
		"no errors on valid input": {
			InputKey:       "signature_algorithm",
			InputValue:     "SHA-1",
			ExpectedOutput: nil,
		},
		"errors on invalid input": {
			InputKey:       "signature_algorithm",
			InputValue:     "asdf",
			ExpectedOutput: []error{fmt.Errorf("signature_algorithm must be one of [SHA-1 SHA-256 SHA-348 SHA-512], got: asdf")},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, errs := validSignatureAlgorithm(test.InputValue, test.InputKey)
			assert.Equal(t, test.ExpectedOutput, errs)
		})
	}
}
