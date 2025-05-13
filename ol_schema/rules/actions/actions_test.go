package appruleactionsschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a Rules Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["action"])
		assert.NotNil(t, provSchema["expression"])
		assert.NotNil(t, provSchema["value"])
		assert.NotNil(t, provSchema["scriplet"])
		assert.NotNil(t, provSchema["macro"])
	})
}

func mockSetFn(interface{}) int {
	return 0
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.Action
	}{
		"creates and returns the address of an Action struct": {
			ResourceData: map[string]interface{}{
				"action":     "test",
				"expression": ".*",
				"value":      schema.NewSet(mockSetFn, []interface{}{"test"}),
				"scriplet":   "",
				"macro":      "",
			},
			ExpectedOutput: models.Action{
				Action:     "test",
				Expression: ".*",
				Value:      []string{"test"},
				Scriplet:   "",
				Macro:      "",
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
	t.Run("It flattens the Action Struct", func(t *testing.T) {
		appActionStruct := []models.Action{
			{
				Action:     "test",
				Expression: ".*",
				Value:      []string{"test"},
				Scriplet:   "",
				Macro:      "",
			},
			{
				Action:     "test2",
				Expression: ".*",
				Value:      []string{"test2"},
				Scriplet:   "",
				Macro:      "",
			},
		}
		subj := Flatten(appActionStruct)
		expected := []map[string]interface{}{
			{
				"action":     "test",
				"expression": ".*",
				"value":      []string{"test"},
				"scriplet":   "",
				"macro":      "",
			},
			{
				"action":     "test2",
				"expression": ".*",
				"value":      []string{"test2"},
				"scriplet":   "",
				"macro":      "",
			},
		}
		assert.Equal(t, expected, subj)
	})
}
