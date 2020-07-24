package usermappingschema

import (
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/user_mappings"
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
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput usermappings.UserMapping
	}{
		"creates and returns the address of an user mapping struct": {
			ResourceData: map[string]interface{}{
				"id":       123,
				"name":     "test",
				"match":    "test",
				"enabled":  true,
				"position": 1,
				"conditions": []interface{}{
					map[string]interface{}{
						"source":   "test",
						"operator": "=",
						"value":    "test",
					},
				},
				"actions": []interface{}{
					map[string]interface{}{
						"action":     "test",
						"expression": ".*",
						"value":      []interface{}{"test"},
					},
				},
			},
			ExpectedOutput: usermappings.UserMapping{
				ID:       oltypes.Int32(int32(123)),
				Name:     oltypes.String("test"),
				Match:    oltypes.String("test"),
				Enabled:  oltypes.Bool(true),
				Position: oltypes.Int32(int32(1)),
				Conditions: []usermappings.UserMappingConditions{
					usermappings.UserMappingConditions{
						Source:   oltypes.String("test"),
						Operator: oltypes.String("="),
						Value:    oltypes.String("test"),
					},
				},
				Actions: []usermappings.UserMappingActions{
					usermappings.UserMappingActions{
						Action: oltypes.String("test"),
						Value:  []string{"test"},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the user mapping Struct", func(t *testing.T) {
		UserMappingStruct := []usermappings.UserMapping{
			usermappings.UserMapping{
				ID:       oltypes.Int32(int32(123)),
				Name:     oltypes.String("test"),
				Match:    oltypes.String("test"),
				Enabled:  oltypes.Bool(true),
				Position: oltypes.Int32(int32(1)),
				Conditions: []usermappings.UserMappingConditions{
					usermappings.UserMappingConditions{
						Source:   oltypes.String("test"),
						Operator: oltypes.String("="),
						Value:    oltypes.String("test"),
					},
				},
				Actions: []usermappings.UserMappingActions{
					usermappings.UserMappingActions{
						Action: oltypes.String("test"),
						Value:  []string{"test"},
					},
				},
			},
			usermappings.UserMapping{
				ID:       oltypes.Int32(int32(456)),
				Name:     oltypes.String("test2"),
				Match:    oltypes.String("test2"),
				Enabled:  oltypes.Bool(true),
				Position: oltypes.Int32(int32(2)),
				Conditions: []usermappings.UserMappingConditions{
					usermappings.UserMappingConditions{
						Source:   oltypes.String("test2"),
						Operator: oltypes.String(">"),
						Value:    oltypes.String("test2"),
					},
				},
				Actions: []usermappings.UserMappingActions{
					usermappings.UserMappingActions{
						Action: oltypes.String("test2"),
						Value:  []string{"test2"},
					},
				},
			},
		}
		subj := Flatten(UserMappingStruct)
		expected := []map[string]interface{}{
			map[string]interface{}{
				"id":       oltypes.Int32(int32(123)),
				"name":     oltypes.String("test"),
				"match":    oltypes.String("test"),
				"enabled":  oltypes.Bool(true),
				"position": oltypes.Int32(int32(1)),
				"conditions": []map[string]interface{}{
					map[string]interface{}{
						"source":   oltypes.String("test"),
						"operator": oltypes.String("="),
						"value":    oltypes.String("test"),
					},
				},
				"actions": []map[string]interface{}{
					map[string]interface{}{
						"action": oltypes.String("test"),
						"value":  []string{"test"},
					},
				},
			},
			map[string]interface{}{
				"id":       oltypes.Int32(int32(456)),
				"name":     oltypes.String("test2"),
				"match":    oltypes.String("test2"),
				"enabled":  oltypes.Bool(true),
				"position": oltypes.Int32(int32(2)),
				"conditions": []map[string]interface{}{
					map[string]interface{}{
						"source":   oltypes.String("test2"),
						"operator": oltypes.String(">"),
						"value":    oltypes.String("test2"),
					},
				},
				"actions": []map[string]interface{}{
					map[string]interface{}{
						"action": oltypes.String("test2"),
						"value":  []string{"test2"},
					},
				},
			},
		}
		assert.Equal(t, expected, subj)
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
