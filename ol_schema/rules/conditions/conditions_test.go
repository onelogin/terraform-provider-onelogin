package appruleconditionsschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
		ExpectedOutput models.Condition
	}{
		"creates and returns the address of an Condition struct": {
			ResourceData: map[string]interface{}{
				"source":   "test",
				"operator": "=",
				"value":    "test",
			},
			ExpectedOutput: models.Condition{
				Source:   "test",
				Operator: "=",
				Value:    "test",
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
	t.Run("It flattens the Condition Struct", func(t *testing.T) {
		conditionStruct := []models.Condition{
			{
				Source:   "test",
				Operator: "=",
				Value:    "test",
			},
			{
				Source:   "test2",
				Operator: "<",
				Value:    "test2",
			},
		}
		subj := Flatten(conditionStruct)
		expected := []map[string]interface{}{
			{
				"source":   "test",
				"operator": "=",
				"value":    "test",
			},
			{
				"source":   "test2",
				"operator": "<",
				"value":    "test2",
			},
		}
		assert.Equal(t, expected, subj)
	})
}
