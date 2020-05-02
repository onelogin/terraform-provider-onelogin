package configuration

import (
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/stretchr/testify/assert"
)

func TestOIDCSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := OIDCSchema()
		assert.NotNil(t, schema["redirect_uri"])
		assert.NotNil(t, schema["refresh_token_expiration_minutes"])
		assert.NotNil(t, schema["login_url"])
		assert.NotNil(t, schema["oidc_application_type"])
		assert.NotNil(t, schema["token_endpoint_auth_method"])
		assert.NotNil(t, schema["access_token_expiration_minutes"])
	})
}

func TestSAMLSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := SAMLSchema()
		assert.NotNil(t, schema["provider_arn"])
		assert.NotNil(t, schema["signature_algorithm"])
	})
}

func TestInflateConfiguration(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.AppConfiguration
	}{
		"creates and returns the address of an AppConfiguration struct for a OIDC app": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": 2,
				"login_url":                        "test",
				"oidc_application_type":            2,
				"token_endpoint_auth_method":       2,
				"access_token_expiration_minutes":  2,
			},
			ExpectedOutput: models.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
			},
		},
		"creates and returns the address of an AppConfiguration struct for a SAML app": {
			ResourceData: map[string]interface{}{
				"provider_arn":        "test",
				"signature_algorithm": "test",
			},
			ExpectedOutput: models.AppConfiguration{
				ProviderArn:        oltypes.String("test"),
				SignatureAlgorithm: oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}

func TestFlattenConfiguration(t *testing.T) {
	tests := map[string]struct {
		InputData      models.AppConfiguration
		ExpectedOutput []map[string]interface{}
	}{
		"creates and returns the address of an AppConfiguration struct for a OIDC app": {
			InputData: models.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
			},
			ExpectedOutput: []map[string]interface{}{
				map[string]interface{}{
					"redirect_uri":                     oltypes.String("test"),
					"refresh_token_expiration_minutes": oltypes.Int32(2),
					"login_url":                        oltypes.String("test"),
					"oidc_application_type":            oltypes.Int32(2),
					"token_endpoint_auth_method":       oltypes.Int32(2),
					"access_token_expiration_minutes":  oltypes.Int32(2),
				},
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
		InputData      models.AppConfiguration
		ExpectedOutput []map[string]interface{}
	}{
		"creates and returns the address of an AppConfiguration struct for a SAML app": {
			InputData: models.AppConfiguration{
				ProviderArn:        oltypes.String("test"),
				SignatureAlgorithm: oltypes.String("test"),
			},
			ExpectedOutput: []map[string]interface{}{
				map[string]interface{}{
					"provider_arn":        oltypes.String("test"),
					"signature_algorithm": oltypes.String("test"),
				},
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

func TestValidSignatureAlgo(t *testing.T) {
	validOpts := []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"}
	tests := map[string]struct {
		InputData      string
		ExpectedOutput []error
	}{
		"no errors on valid input": {
			InputData:      "SHA-1",
			ExpectedOutput: nil,
		},
		"errors on invalid input": {
			InputData:      "asdf",
			ExpectedOutput: []error{fmt.Errorf("signature_algorithm must be one of %v, got: %s", validOpts, "asdf")},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, errs := validSignatureAlgo(test.InputData, "signature_algorithm")
			assert.Equal(t, errs, test.ExpectedOutput)
		})
	}
}
