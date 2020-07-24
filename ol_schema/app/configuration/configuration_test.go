package appconfigurationschema

import (
	"errors"
	"fmt"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInflateConfiguration(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput apps.AppConfiguration
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
			ExpectedOutput: apps.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
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
			},
			ExpectedOutput: apps.AppConfiguration{
				ProviderArn:        oltypes.String("test"),
				SignatureAlgorithm: oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, err := Inflate(test.ResourceData)
			if (test.ExpectedOutput != apps.AppConfiguration{}) {
				assert.Equal(t, subj, test.ExpectedOutput)
			}
			if test.ExpectedError != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestFlattenConfiguration(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppConfiguration
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns the address of an AppConfiguration struct for a OIDC app": {
			InputData: apps.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
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
		InputData      apps.AppConfiguration
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns the address of an AppConfiguration struct for a SAML app": {
			InputData: apps.AppConfiguration{
				ProviderArn:        oltypes.String("test"),
				SignatureAlgorithm: oltypes.String("test"),
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
