package authserverschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of an AuthServer Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["description"])
		assert.NotNil(t, provSchema["configuration"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput authservers.AuthServer
	}{
		"creates and returns the address of a user struct": {
			ResourceData: map[string]interface{}{
				"name":        "name",
				"description": "description",
				"firstname":   "description",
				"configuration": []interface{}{
					map[string]interface{}{
						"resource_identifier":              "test.com",
						"audiences":                        []string{"aud_1", "aud_2"},
						"refresh_token_expiration_minutes": 2,
						"access_token_expiration_minutes":  2,
					},
				},
			},
			ExpectedOutput: authservers.AuthServer{
				Name:        oltypes.String("name"),
				Description: oltypes.String("description"),
				Configuration: &authservers.AuthServerConfiguration{
					ResourceIdentifier:            oltypes.String("test.com"),
					Audiences:                     []string{"aud_1", "aud_2"},
					AccessTokenExpirationMinutes:  oltypes.Int32(2),
					RefreshTokenExpirationMinutes: oltypes.Int32(2),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
