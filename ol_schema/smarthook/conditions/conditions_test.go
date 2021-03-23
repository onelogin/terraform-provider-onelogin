package smarthookconditionsschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a Rules Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["source"])
		assert.NotNil(t, provSchema["operator"])
		assert.NotNil(t, provSchema["value"])
	})
}

func Test(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput smarthooks.Condition
	}{
		"creates and returns the address of an AppParameters struct": {
			ResourceData: map[string]interface{}{
				"source":   "test",
				"operator": "=",
				"value":    "test",
			},
			ExpectedOutput: smarthooks.Condition{
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
	t.Run("It flattens the AppParameters Struct", func(t *testing.T) {
		appConditionStruct := []smarthooks.Condition{
			smarthooks.Condition{
				Source:   oltypes.String("test"),
				Operator: oltypes.String("="),
				Value:    oltypes.String("test"),
			},
			smarthooks.Condition{
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
