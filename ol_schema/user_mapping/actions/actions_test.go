package usermappingactionsschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a user mapping Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["action"])
		assert.NotNil(t, provSchema["value"])
	})
}

func TestInflate(t *testing.T) {
	// Set up test variables
	action := "test"

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.UserMappingActions
	}{
		"creates and returns the address of an user mappings action struct": {
			ResourceData: map[string]interface{}{
				"action": action,
				"value":  []interface{}{"test"},
			},
			ExpectedOutput: models.UserMappingActions{
				Action: &action,
				Value:  []string{"test"},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			if subj.Action != nil && test.ExpectedOutput.Action != nil {
				assert.Equal(t, *test.ExpectedOutput.Action, *subj.Action)
			}
			assert.Equal(t, test.ExpectedOutput.Value, subj.Value)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the user mapping actions Struct", func(t *testing.T) {
		// Set up test variables
		action1 := "test"
		action2 := "test2"

		appConditionStruct := []models.UserMappingActions{
			{
				Action: &action1,
				Value:  []string{"test"},
			},
			{
				Action: &action2,
				Value:  []string{"test2"},
			},
		}
		subj := Flatten(appConditionStruct)
		expected := []map[string]interface{}{
			{
				"action": &action1,
				"value":  []string{"test"},
			},
			{
				"action": &action2,
				"value":  []string{"test2"},
			},
		}
		assert.Equal(t, len(expected), len(subj))
		// Test specific field values
		for i, exp := range expected {
			assert.Equal(t, exp["action"], subj[i]["action"])
			assert.Equal(t, exp["value"], subj[i]["value"])
		}
	})
}
