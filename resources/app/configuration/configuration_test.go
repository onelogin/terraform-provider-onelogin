package configuration

import (
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
	"github.com/stretchr/testify/assert"
)

func TestOIDCConfigurationSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := OIDCConfigurationSchema()
		assert.NotNil(t, schema["redirect_uri"])
		assert.NotNil(t, schema["refresh_token_expiration_minutes"])
		assert.NotNil(t, schema["login_url"])
		assert.NotNil(t, schema["oidc_application_type"])
		assert.NotNil(t, schema["token_endpoint_auth_method"])
		assert.NotNil(t, schema["access_token_expiration_minutes"])
	})
}

func TestSAMLConfigurationSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := SAMLConfigurationSchema()
		assert.NotNil(t, schema["provider_arn"])
		assert.NotNil(t, schema["signature_algorithm"])
	})
}

func TestAddConfigurationSchema(t *testing.T) {
	t.Run("adds configuration schema to given resrouce schema", func(t *testing.T) {
		appSchema := app.AppSchema()
		AddConfigurationSchema(&appSchema, SAMLConfigurationSchema)
		assert.NotNil(t, appSchema["configuration"])
	})
}

func TestInflateOIDCConfiguration(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput *models.AppConfiguration
	}{
		"creates and returns the address of an AppConfiguration struct": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": 2,
				"login_url":                        "test",
				"oidc_application_type":            2,
				"token_endpoint_auth_method":       2,
				"access_token_expiration_minutes":  2,
			},
			ExpectedOutput: &models.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := InflateOIDCConfiguration(&test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}

func TestInflateSAMLConfiguration(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput *models.AppConfiguration
	}{
		"creates and returns the address of an AppConfiguration struct": {
			ResourceData: map[string]interface{}{
				"provider_arn":        "test",
				"signature_algorithm": "test",
			},
			ExpectedOutput: &models.AppConfiguration{
				ProviderArn:        oltypes.String("test"),
				SignatureAlgorithm: oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := InflateSAMLConfiguration(&test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
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
