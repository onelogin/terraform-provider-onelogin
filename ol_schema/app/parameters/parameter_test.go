package appparametersschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
)

func TestParameterSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppParameter Schema", func(t *testing.T) {
		schema := Schema()
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
		ExpectedOutput apps.AppParameters
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
				"include_in_saml_assertion":  true,
			},
			ExpectedOutput: apps.AppParameters{
				ID:                        oltypes.Int32(int32(123)),
				Label:                     oltypes.String("test"),
				UserAttributeMappings:     oltypes.String("test"),
				UserAttributeMacros:       oltypes.String("test"),
				AttributesTransformations: oltypes.String("test"),
				SkipIfBlank:               oltypes.Bool(true),
				Values:                    oltypes.String("test"),
				DefaultValues:             oltypes.String("test"),
				ProvisionedEntitlements:   oltypes.Bool(true),
				SafeEntitlementsEnabled:   oltypes.Bool(true),
				IncludeInSamlAssertion:    oltypes.Bool(true),
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

func TestFlatten(t *testing.T) {
	t.Run("It flattens the AppParameters Struct", func(t *testing.T) {
		appParamStruct := map[string]apps.AppParameters{
			"test": apps.AppParameters{
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
				IncludeInSamlAssertion:    oltypes.Bool(true),
			},
		}
		subj := Flatten(appParamStruct)
		expected := []map[string]interface{}{
			map[string]interface{}{
				"param_key_name":             "test",
				"param_id":                   oltypes.Int32(123),
				"label":                      oltypes.String("test"),
				"user_attribute_mappings":    oltypes.String("test"),
				"user_attribute_macros":      oltypes.String("test"),
				"attributes_transformations": oltypes.String("test"),
				"skip_if_blank":              oltypes.Bool(true),
				"values":                     oltypes.String("test"),
				"default_values":             oltypes.String("test"),
				"provisioned_entitlements":   oltypes.Bool(true),
				"safe_entitlements_enabled":  oltypes.Bool(true),
				"include_in_saml_assertion":  oltypes.Bool(true),
			},
		}
		assert.Equal(t, expected, subj)
	})
}
