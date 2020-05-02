package app

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/stretchr/testify/assert"
)

func TestAppSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := AppSchema()
		assert.NotNil(t, schema["name"])
		assert.NotNil(t, schema["visible"])
		assert.NotNil(t, schema["description"])
		assert.NotNil(t, schema["notes"])
		assert.NotNil(t, schema["icon_url"])
		assert.NotNil(t, schema["auth_method"])
		assert.NotNil(t, schema["policy_id"])
		assert.NotNil(t, schema["allow_assumed_signin"])
		assert.NotNil(t, schema["tab_id"])
		assert.NotNil(t, schema["connector_id"])
		assert.NotNil(t, schema["created_at"])
		assert.NotNil(t, schema["updated_at"])
		assert.NotNil(t, schema["provisioning"])
		assert.NotNil(t, schema["parameters"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.App
	}{
		"creates and returns the address of an AppParameters struct with all sub-fiekds": {
			ResourceData: map[string]interface{}{
				"name":                 "test",
				"visible":              true,
				"description":          "test",
				"notes":                "test",
				"allow_assumed_signin": true,
				"connector_id":         123,
				"parameters": []interface{}{
					map[string]interface{}{
						"param_key_name":             "test",
						"param_id":                   123,
						"label":                      "test",
						"user_attribute_mappings":    "test",
						"user_attribute_macros":      "test",
						"attributes_transformations": "test",
						"default_values":             "test",
						"skip_if_blank":              true,
						"values":                     "test",
						"provisioned_entitlements":   true,
						"safe_entitlements_enabled":  true,
					},
				},
				"provisioning": []interface{}{
					map[string]interface{}{
						"enabled": true,
					},
				},
				"configuration": []interface{}{
					map[string]interface{}{
						"provider_arn":        "test",
						"signature_algorithm": "test",
					},
				},
			},
			ExpectedOutput: models.App{
				Name:               oltypes.String("test"),
				Visible:            oltypes.Bool(true),
				Description:        oltypes.String("test"),
				Notes:              oltypes.String("test"),
				AllowAssumedSignin: oltypes.Bool(true),
				ConnectorID:        oltypes.Int32(123),
				Parameters: map[string]models.AppParameters{
					"test": models.AppParameters{
						ID:                        oltypes.Int32(123),
						Label:                     oltypes.String("test"),
						UserAttributeMappings:     oltypes.String("test"),
						UserAttributeMacros:       oltypes.String("test"),
						AttributesTransformations: oltypes.String("test"),
						DefaultValues:             oltypes.String("test"),
						SkipIfBlank:               oltypes.Bool(true),
						Values:                    oltypes.String("test"),
						ProvisionedEntitlements:   oltypes.Bool(true),
						SafeEntitlementsEnabled:   oltypes.Bool(true),
					},
				},
				Provisioning: &models.AppProvisioning{
					Enabled: oltypes.Bool(true),
				},
				Configuration: &models.AppConfiguration{
					ProviderArn:        oltypes.String("test"),
					SignatureAlgorithm: oltypes.String("test"),
				},
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
