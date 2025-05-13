package authserverschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
	name := "name"
	description := "description"
	resourceIdentifier := "test.com"
	audiences := []string{"aud_1", "aud_2"}
	accessTokenExpirationMinutes := int32(2)
	refreshTokenExpirationMinutes := int32(2)

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.AuthServer
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
			ExpectedOutput: models.AuthServer{
				Name:        &name,
				Description: &description,
				Configuration: &models.AuthServerConfiguration{
					ResourceIdentifier:            &resourceIdentifier,
					Audiences:                     audiences,
					AccessTokenExpirationMinutes:  &accessTokenExpirationMinutes,
					RefreshTokenExpirationMinutes: &refreshTokenExpirationMinutes,
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)

			// Check basic fields
			if subj.Name != nil && test.ExpectedOutput.Name != nil {
				assert.Equal(t, *test.ExpectedOutput.Name, *subj.Name)
			}
			if subj.Description != nil && test.ExpectedOutput.Description != nil {
				assert.Equal(t, *test.ExpectedOutput.Description, *subj.Description)
			}

			// Check Configuration
			if subj.Configuration != nil && test.ExpectedOutput.Configuration != nil {
				confSubj := subj.Configuration
				confExp := test.ExpectedOutput.Configuration

				if confSubj.ResourceIdentifier != nil && confExp.ResourceIdentifier != nil {
					assert.Equal(t, *confExp.ResourceIdentifier, *confSubj.ResourceIdentifier)
				}

				assert.Equal(t, confExp.Audiences, confSubj.Audiences)

				if confSubj.AccessTokenExpirationMinutes != nil && confExp.AccessTokenExpirationMinutes != nil {
					assert.Equal(t, *confExp.AccessTokenExpirationMinutes, *confSubj.AccessTokenExpirationMinutes)
				}

				if confSubj.RefreshTokenExpirationMinutes != nil && confExp.RefreshTokenExpirationMinutes != nil {
					assert.Equal(t, *confExp.RefreshTokenExpirationMinutes, *confSubj.RefreshTokenExpirationMinutes)
				}
			}
		})
	}
}
