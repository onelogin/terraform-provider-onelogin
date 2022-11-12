package authserverconfigurationschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of an AuthServerConfiguration Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["resource_identifier"])
		assert.NotNil(t, provSchema["audiences"])
		assert.NotNil(t, provSchema["access_token_expiration_minutes"])
		assert.NotNil(t, provSchema["refresh_token_expiration_minutes"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   []interface{}
		ExpectedOutput authservers.AuthServerConfiguration
	}{
		"creates and returns the address of an AuthServerConfiguration": {
			ResourceData: []interface{}{
				map[string]interface{}{
					"resource_identifier":              "test.com",
					"audiences":                        []string{"aud_1", "aud_2"},
					"refresh_token_expiration_minutes": 2,
					"access_token_expiration_minutes":  2,
				},
			},
			ExpectedOutput: authservers.AuthServerConfiguration{
				ResourceIdentifier:            oltypes.String("test.com"),
				Audiences:                     []string{"aud_1", "aud_2"},
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := map[string]struct {
		Input  authservers.AuthServerConfiguration
		Output []map[string]interface{}
	}{
		"converts the AuthServerConfiguration to a map": {
			Input: authservers.AuthServerConfiguration{
				ResourceIdentifier:            oltypes.String("test.com"),
				Audiences:                     []string{"aud_1", "aud_2"},
				AccessTokenExpirationMinutes:  oltypes.Int32(2),
				RefreshTokenExpirationMinutes: oltypes.Int32(2),
			},
			Output: []map[string]interface{}{
				map[string]interface{}{
					"resource_identifier":              "test.com",
					"audiences":                        []string{"aud_1", "aud_2"},
					"refresh_token_expiration_minutes": int32(2),
					"access_token_expiration_minutes":  int32(2),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Flatten(test.Input)
			assert.Equal(t, test.Output, subj)
		})
	}
}
