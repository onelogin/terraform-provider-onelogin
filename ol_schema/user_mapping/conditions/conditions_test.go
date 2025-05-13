package usermappingconditionsschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
	// Set up test variables
	source := "test"
	operator := "="
	value := "test"

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.UserMappingConditions
	}{
		"creates and returns the address of an user mapping conditions struct": {
			ResourceData: map[string]interface{}{
				"source":   source,
				"operator": operator,
				"value":    value,
			},
			ExpectedOutput: models.UserMappingConditions{
				Source:   &source,
				Operator: &operator,
				Value:    &value,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			if subj.Source != nil && test.ExpectedOutput.Source != nil {
				assert.Equal(t, *test.ExpectedOutput.Source, *subj.Source)
			}
			if subj.Operator != nil && test.ExpectedOutput.Operator != nil {
				assert.Equal(t, *test.ExpectedOutput.Operator, *subj.Operator)
			}
			if subj.Value != nil && test.ExpectedOutput.Value != nil {
				assert.Equal(t, *test.ExpectedOutput.Value, *subj.Value)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the user mappings condition Struct", func(t *testing.T) {
		// Set up test variables
		source1 := "test"
		operator1 := "="
		value1 := "test"

		source2 := "test2"
		operator2 := "<"
		value2 := "test2"

		appConditionStruct := []models.UserMappingConditions{
			{
				Source:   &source1,
				Operator: &operator1,
				Value:    &value1,
			},
			{
				Source:   &source2,
				Operator: &operator2,
				Value:    &value2,
			},
		}
		subj := Flatten(appConditionStruct)
		expected := []map[string]interface{}{
			{
				"source":   &source1,
				"operator": &operator1,
				"value":    &value1,
			},
			{
				"source":   &source2,
				"operator": &operator2,
				"value":    &value2,
			},
		}
		assert.Equal(t, len(expected), len(subj))
		// Test specific fields
		for i, exp := range expected {
			assert.Equal(t, exp["source"], subj[i]["source"])
			assert.Equal(t, exp["operator"], subj[i]["operator"])
			assert.Equal(t, exp["value"], subj[i]["value"])
		}
	})
}
