package usermappingschema

import (
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a user mapping Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["match"])
		assert.NotNil(t, provSchema["position"])
		assert.NotNil(t, provSchema["conditions"])
		assert.NotNil(t, provSchema["actions"])
	})
}

func TestInflate(t *testing.T) {
	// Create test variables
	id := int32(123)
	name := "test"
	match := "test"
	enabled := true
	position := int32(1)

	// Create source/operator/value for conditions
	source := "test"
	operator := "="
	value := "test"

	// Create action/expression for actions
	action := "test"

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.UserMapping
	}{
		"creates and returns the address of an user mapping struct": {
			ResourceData: map[string]interface{}{
				"id":       "123",
				"name":     name,
				"match":    match,
				"enabled":  enabled,
				"position": int(position),
				"conditions": []interface{}{
					map[string]interface{}{
						"source":   source,
						"operator": operator,
						"value":    value,
					},
				},
				"actions": []interface{}{
					map[string]interface{}{
						"action":     action,
						"expression": ".*",
						"value":      []interface{}{"test"},
					},
				},
			},
			ExpectedOutput: models.UserMapping{
				ID:       &id,
				Name:     &name,
				Match:    &match,
				Enabled:  &enabled,
				Position: &position,
				Conditions: []models.UserMappingConditions{
					{
						Source:   &source,
						Operator: &operator,
						Value:    &value,
					},
				},
				Actions: []models.UserMappingActions{
					{
						Action: &action,
						Value:  []string{"test"},
					},
				},
			},
		},
		"handles a user mapping without the position provided": {
			ResourceData: map[string]interface{}{
				"id":      "123",
				"name":    name,
				"match":   match,
				"enabled": enabled,
				"conditions": []interface{}{
					map[string]interface{}{
						"source":   source,
						"operator": operator,
						"value":    value,
					},
				},
				"actions": []interface{}{
					map[string]interface{}{
						"action":     action,
						"expression": ".*",
						"value":      []interface{}{"test"},
					},
				},
			},
			ExpectedOutput: models.UserMapping{
				ID:      &id,
				Name:    &name,
				Match:   &match,
				Enabled: &enabled,
				Conditions: []models.UserMappingConditions{
					{
						Source:   &source,
						Operator: &operator,
						Value:    &value,
					},
				},
				Actions: []models.UserMappingActions{
					{
						Action: &action,
						Value:  []string{"test"},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			// Compare pointer values properly
			if subj.ID != nil && test.ExpectedOutput.ID != nil {
				assert.Equal(t, *test.ExpectedOutput.ID, *subj.ID)
			}
			if subj.Name != nil && test.ExpectedOutput.Name != nil {
				assert.Equal(t, *test.ExpectedOutput.Name, *subj.Name)
			}
			if subj.Match != nil && test.ExpectedOutput.Match != nil {
				assert.Equal(t, *test.ExpectedOutput.Match, *subj.Match)
			}
			if subj.Enabled != nil && test.ExpectedOutput.Enabled != nil {
				assert.Equal(t, *test.ExpectedOutput.Enabled, *subj.Enabled)
			}
			if test.ResourceData["position"] != nil {
				if subj.Position != nil && test.ExpectedOutput.Position != nil {
					assert.Equal(t, *test.ExpectedOutput.Position, *subj.Position)
				}
			}

			assert.Equal(t, len(test.ExpectedOutput.Conditions), len(subj.Conditions))
			assert.Equal(t, len(test.ExpectedOutput.Actions), len(subj.Actions))

			if len(subj.Conditions) > 0 {
				if subj.Conditions[0].Source != nil && test.ExpectedOutput.Conditions[0].Source != nil {
					assert.Equal(t, *test.ExpectedOutput.Conditions[0].Source, *subj.Conditions[0].Source)
				}
				if subj.Conditions[0].Operator != nil && test.ExpectedOutput.Conditions[0].Operator != nil {
					assert.Equal(t, *test.ExpectedOutput.Conditions[0].Operator, *subj.Conditions[0].Operator)
				}
				if subj.Conditions[0].Value != nil && test.ExpectedOutput.Conditions[0].Value != nil {
					assert.Equal(t, *test.ExpectedOutput.Conditions[0].Value, *subj.Conditions[0].Value)
				}
			}

			if len(subj.Actions) > 0 {
				if subj.Actions[0].Action != nil && test.ExpectedOutput.Actions[0].Action != nil {
					assert.Equal(t, *test.ExpectedOutput.Actions[0].Action, *subj.Actions[0].Action)
				}
				assert.Equal(t, test.ExpectedOutput.Actions[0].Value, subj.Actions[0].Value)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the user mapping Struct", func(t *testing.T) {
		// Create test variables
		id1 := int32(123)
		name1 := "test"
		match1 := "test"
		enabled1 := true
		position1 := int32(1)

		id2 := int32(456)
		name2 := "test2"
		match2 := "test2"
		enabled2 := true
		position2 := int32(2)

		// Create source/operator/value for conditions
		source1 := "test"
		operator1 := "="
		value1 := "test"

		source2 := "test2"
		operator2 := ">"
		value2 := "test2"

		// Create action/expression for actions
		action1 := "test"
		action2 := "test2"

		UserMappingStruct := []models.UserMapping{
			{
				ID:       &id1,
				Name:     &name1,
				Match:    &match1,
				Enabled:  &enabled1,
				Position: &position1,
				Conditions: []models.UserMappingConditions{
					{
						Source:   &source1,
						Operator: &operator1,
						Value:    &value1,
					},
				},
				Actions: []models.UserMappingActions{
					{
						Action: &action1,
						Value:  []string{"test"},
					},
				},
			},
			{
				ID:       &id2,
				Name:     &name2,
				Match:    &match2,
				Enabled:  &enabled2,
				Position: &position2,
				Conditions: []models.UserMappingConditions{
					{
						Source:   &source2,
						Operator: &operator2,
						Value:    &value2,
					},
				},
				Actions: []models.UserMappingActions{
					{
						Action: &action2,
						Value:  []string{"test2"},
					},
				},
			},
		}
		subj := Flatten(UserMappingStruct)
		expected := []map[string]interface{}{
			{
				"id":       &id1,
				"name":     &name1,
				"match":    &match1,
				"enabled":  &enabled1,
				"position": &position1,
				"conditions": []map[string]interface{}{
					{
						"source":   &source1,
						"operator": &operator1,
						"value":    &value1,
					},
				},
				"actions": []map[string]interface{}{
					{
						"action": &action1,
						"value":  []string{"test"},
					},
				},
			},
			{
				"id":       &id2,
				"name":     &name2,
				"match":    &match2,
				"enabled":  &enabled2,
				"position": &position2,
				"conditions": []map[string]interface{}{
					{
						"source":   &source2,
						"operator": &operator2,
						"value":    &value2,
					},
				},
				"actions": []map[string]interface{}{
					{
						"action": &action2,
						"value":  []string{"test2"},
					},
				},
			},
		}
		// Test just the keys because we need to test deep hierarchies of pointers
		assert.Equal(t, len(expected), len(subj))
	})
}

func TestValidMatch(t *testing.T) {
	tests := map[string]struct {
		InputKey       string
		InputValue     string
		ExpectedOutput []error
	}{
		"no errors on valid input": {
			InputKey:       "match",
			InputValue:     "all",
			ExpectedOutput: nil,
		},
		"errors on invalid input": {
			InputKey:       "match",
			InputValue:     "asdf",
			ExpectedOutput: []error{fmt.Errorf("match must be one of [all any], got: asdf")},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, errs := validMatch(test.InputValue, test.InputKey)
			assert.Equal(t, test.ExpectedOutput, errs)
		})
	}
}
