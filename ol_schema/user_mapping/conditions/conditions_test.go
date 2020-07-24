package usermappingconditionsschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/user_mappings"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a user mappings condition Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["source"])
		assert.NotNil(t, provSchema["operator"])
		assert.NotNil(t, provSchema["value"])
	})
}

func Test(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput usermappings.UserMappingConditions
	}{
		"creates and returns the address of an user mapping conditions struct": {
			ResourceData: map[string]interface{}{
				"source":   "test",
				"operator": "=",
				"value":    "test",
			},
			ExpectedOutput: usermappings.UserMappingConditions{
				Source:   oltypes.String("test"),
				Operator: oltypes.String("="),
				Value:    oltypes.String("test"),
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
	t.Run("It flattens the user mappings condition Struct", func(t *testing.T) {
		appConditionStruct := []usermappings.UserMappingConditions{
			usermappings.UserMappingConditions{
				Source:   oltypes.String("test"),
				Operator: oltypes.String("="),
				Value:    oltypes.String("test"),
			},
			usermappings.UserMappingConditions{
				Source:   oltypes.String("test2"),
				Operator: oltypes.String("<"),
				Value:    oltypes.String("test2"),
			},
		}
		subj := Flatten(appConditionStruct)
		expected := []map[string]interface{}{
			map[string]interface{}{
				"source":   oltypes.String("test"),
				"operator": oltypes.String("="),
				"value":    oltypes.String("test"),
			},
			map[string]interface{}{
				"source":   oltypes.String("test2"),
				"operator": oltypes.String("<"),
				"value":    oltypes.String("test2"),
			},
		}
		assert.Equal(t, expected, subj)
	})
}
