package appruleactionsschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	apprules "github.com/onelogin/onelogin-go-sdk/pkg/services/apps/app_rules"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a Rules Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["action"])
		assert.NotNil(t, provSchema["expression"])
		assert.NotNil(t, provSchema["value"])
	})
}

func mockSetFn(interface{}) int {
	return 0
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput apprules.AppRuleActions
	}{
		"creates and returns the address of an AppParameters struct": {
			ResourceData: map[string]interface{}{
				"action":     "test",
				"expression": ".*",
				"value":      schema.NewSet(mockSetFn, []interface{}{"test"}),
			},
			ExpectedOutput: apprules.AppRuleActions{
				Action:     oltypes.String("test"),
				Expression: oltypes.String(".*"),
				Value:      []string{"test"},
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
		appConditionStruct := []apprules.AppRuleActions{
			{
				Action:     oltypes.String("test"),
				Expression: oltypes.String(".*"),
				Value:      []string{"test"},
			},
			{
				Action:     oltypes.String("test2"),
				Expression: oltypes.String(".*"),
				Value:      []string{"test2"},
			},
		}
		subj := Flatten(appConditionStruct)
		expected := []map[string]interface{}{
			{
				"action":     oltypes.String("test"),
				"expression": oltypes.String(".*"),
				"value":      []string{"test"},
			},
			{
				"action":     oltypes.String("test2"),
				"expression": oltypes.String(".*"),
				"value":      []string{"test2"},
			},
		}
		assert.Equal(t, expected, subj)
	})
}
