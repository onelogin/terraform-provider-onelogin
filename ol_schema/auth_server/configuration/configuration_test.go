package authserverconfigurationschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
	// Setup test variables
	resourceID := "test.com"
	atMinutes := int32(2)
	rtMinutes := int32(2)

	tests := map[string]struct {
		ResourceData   []interface{}
		ExpectedOutput models.AuthServerConfiguration
	}{
		"creates and returns the address of an AuthServerConfiguration": {
			ResourceData: []interface{}{
				map[string]interface{}{
					"resource_identifier":              resourceID,
					"audiences":                        []string{"aud_1", "aud_2"},
					"refresh_token_expiration_minutes": int(rtMinutes),
					"access_token_expiration_minutes":  int(atMinutes),
				},
			},
			ExpectedOutput: models.AuthServerConfiguration{
				ResourceIdentifier:            &resourceID,
				Audiences:                     []string{"aud_1", "aud_2"},
				AccessTokenExpirationMinutes:  &atMinutes,
				RefreshTokenExpirationMinutes: &rtMinutes,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			// Compare pointer values
			if subj.ResourceIdentifier != nil && test.ExpectedOutput.ResourceIdentifier != nil {
				assert.Equal(t, *subj.ResourceIdentifier, *test.ExpectedOutput.ResourceIdentifier)
			}
			assert.Equal(t, subj.Audiences, test.ExpectedOutput.Audiences)
			if subj.AccessTokenExpirationMinutes != nil && test.ExpectedOutput.AccessTokenExpirationMinutes != nil {
				assert.Equal(t, *subj.AccessTokenExpirationMinutes, *test.ExpectedOutput.AccessTokenExpirationMinutes)
			}
			if subj.RefreshTokenExpirationMinutes != nil && test.ExpectedOutput.RefreshTokenExpirationMinutes != nil {
				assert.Equal(t, *subj.RefreshTokenExpirationMinutes, *test.ExpectedOutput.RefreshTokenExpirationMinutes)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	// Setup test variables
	resourceID := "test.com"
	atMinutes := int32(2)
	rtMinutes := int32(2)

	tests := map[string]struct {
		Input  models.AuthServerConfiguration
		Output map[string]interface{}
	}{
		"converts the AuthServerConfiguration to a map": {
			Input: models.AuthServerConfiguration{
				ResourceIdentifier:            &resourceID,
				Audiences:                     []string{"aud_1", "aud_2"},
				AccessTokenExpirationMinutes:  &atMinutes,
				RefreshTokenExpirationMinutes: &rtMinutes,
			},
			Output: map[string]interface{}{
				"resource_identifier":              resourceID,
				"audiences":                        []string{"aud_1", "aud_2"},
				"refresh_token_expiration_minutes": rtMinutes,
				"access_token_expiration_minutes":  atMinutes,
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
