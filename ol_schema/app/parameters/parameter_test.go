package appparametersschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
		ExpectedOutput models.Parameter
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
			ExpectedOutput: models.Parameter{
				ID:                        123,
				Label:                     "test",
				UserAttributeMappings:     "test",
				UserAttributeMacros:       "test",
				AttributesTransformations: "test",
				SkipIfBlank:               true,
				Values:                    "test",
				DefaultValues:             "test",
				ProvisionedEntitlements:   true,
				IncludeInSamlAssertion:    true,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, subj.ID, test.ExpectedOutput.ID)
			assert.Equal(t, subj.Label, test.ExpectedOutput.Label)
			assert.Equal(t, subj.UserAttributeMappings, test.ExpectedOutput.UserAttributeMappings)
			assert.Equal(t, subj.UserAttributeMacros, test.ExpectedOutput.UserAttributeMacros)
			assert.Equal(t, subj.AttributesTransformations, test.ExpectedOutput.AttributesTransformations)
			assert.Equal(t, subj.SkipIfBlank, test.ExpectedOutput.SkipIfBlank)
			assert.Equal(t, subj.Values, test.ExpectedOutput.Values)
			assert.Equal(t, subj.DefaultValues, test.ExpectedOutput.DefaultValues)
			assert.Equal(t, subj.ProvisionedEntitlements, test.ExpectedOutput.ProvisionedEntitlements)
			assert.Equal(t, subj.IncludeInSamlAssertion, test.ExpectedOutput.IncludeInSamlAssertion)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the AppParameters Struct", func(t *testing.T) {
		appParamStruct := map[string]models.Parameter{
			"test": {
				ID:                        123,
				Label:                     "test",
				UserAttributeMappings:     "test",
				UserAttributeMacros:       "test",
				AttributesTransformations: "test",
				SkipIfBlank:               true,
				Values:                    "test",
				DefaultValues:             "test",
				ProvisionedEntitlements:   true,
				IncludeInSamlAssertion:    true,
			},
		}
		subj := Flatten(appParamStruct)
		expected := []map[string]interface{}{
			{
				"param_key_name":             "test",
				"param_id":                   123,
				"label":                      "test",
				"user_attribute_mappings":    "test",
				"user_attribute_macros":      "test",
				"attributes_transformations": "test",
				"skip_if_blank":              true,
				"values":                     "test",
				"default_values":             "test",
				"provisioned_entitlements":   true,
				"include_in_saml_assertion":  true,
			},
		}
		assert.Equal(t, expected, subj)
	})
}
