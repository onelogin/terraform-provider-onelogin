package app

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/stretchr/testify/assert"
)

func TestParameterSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppParameter Schema", func(t *testing.T) {
		schema := ParameterSchema()
		assert.NotNil(t, schema["param_key_name"])
		assert.NotNil(t, schema["param_id"])
		assert.NotNil(t, schema["label"])
		assert.NotNil(t, schema["user_attribute_mappings"])
		assert.NotNil(t, schema["user_attribute_macros"])
		assert.NotNil(t, schema["attributes_transformations"])
		assert.NotNil(t, schema["default_values"])
		assert.NotNil(t, schema["skip_if_blank"])
		assert.NotNil(t, schema["values"])
		assert.NotNil(t, schema["provisioned_entitlements"])
		assert.NotNil(t, schema["safe_entitlements_enabled"])
	})
}

func TestInflateParameter(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput *models.AppParameters
	}{
		"creates and returns the address of an AppParameters struct": {
			ResourceData: map[string]interface{}{
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
			ExpectedOutput: &models.AppParameters{
				ID:                        oltypes.Int32(123),
				Label:                     oltypes.String("test"),
				UserAttributeMappings:     oltypes.String("test"),
				UserAttributeMacros:       oltypes.String("test"),
				AttributesTransformations: oltypes.String("test"),
				SkipIfBlank:               oltypes.Bool(true),
				Values:                    oltypes.String("test"),
				DefaultValues:             oltypes.String("test"),
				ProvisionedEntitlements:   oltypes.Bool(true),
				SafeEntitlementsEnabled:   oltypes.Bool(true),
			},
		},
		"ignores unsupplied fields": {
			ResourceData: map[string]interface{}{
				"param_key_name":             "test",
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
			ExpectedOutput: &models.AppParameters{
				Label:                     oltypes.String("test"),
				UserAttributeMappings:     oltypes.String("test"),
				UserAttributeMacros:       oltypes.String("test"),
				AttributesTransformations: oltypes.String("test"),
				SkipIfBlank:               oltypes.Bool(true),
				Values:                    oltypes.String("test"),
				DefaultValues:             oltypes.String("test"),
				ProvisionedEntitlements:   oltypes.Bool(true),
				SafeEntitlementsEnabled:   oltypes.Bool(true),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := InflateParameter(&test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}
