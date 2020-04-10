package app

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := ConfigurationSchema()
		assert.NotNil(t, schema["redirect_uri"])
		assert.NotNil(t, schema["refresh_token_expiration_minutes"])
		assert.NotNil(t, schema["login_url"])
		assert.NotNil(t, schema["oidc_application_type"])
		assert.NotNil(t, schema["token_endpoint_auth_method"])
		assert.NotNil(t, schema["access_token_expiration_minutes"])
		assert.NotNil(t, schema["provider_arn"])
		assert.NotNil(t, schema["signature_algorithm"])
	})
}

func TestInflateConfiguration(t *testing.T) {
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
				"provider_arn":                     "test",
				"signature_algorithm":              "test",
			},
			ExpectedOutput: &models.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
				ProviderArn:                   oltypes.String("test"),
				SignatureAlgorithm:            oltypes.String("test"),
			},
		},
		"ignores unsupplied fields": {
			ResourceData: map[string]interface{}{
				"redirect_uri":                     "test",
				"refresh_token_expiration_minutes": 2,
				"login_url":                        "test",
				"oidc_application_type":            2,
				"token_endpoint_auth_method":       2,
				"access_token_expiration_minutes":  2,
				"provider_arn":                     "test",
			},
			ExpectedOutput: &models.AppConfiguration{
				RedirectURI:                   oltypes.String("test"),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
				LoginURL:                      oltypes.String("test"),
				OidcApplicationType:           oltypes.Int32(2),
				TokenEndpointAuthMethod:       oltypes.Int32(2),
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
				ProviderArn:                   oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := InflateConfiguration(&test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}
