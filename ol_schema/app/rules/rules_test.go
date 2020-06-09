package apprulesschema

import (
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a Rules Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["id"])
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
		ExpectedOutput apps.AppRule
	}{
		"creates and returns the address of an AppParameters struct": {
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
			ExpectedOutput: apps.AppRule{
				ID:       oltypes.Int32(int32(123)),
				Name:     oltypes.String("test"),
				Match:    oltypes.String("test"),
				Enabled:  oltypes.Bool(true),
				Position: oltypes.Int32(int32(1)),
				Conditions: []apps.AppRuleConditions{
					apps.AppRuleConditions{
						Source:   oltypes.String("test"),
						Operator: oltypes.String("="),
						Value:    oltypes.String("test"),
					},
				},
				Actions: []apps.AppRuleActions{
					apps.AppRuleActions{
						Action:     oltypes.String("test"),
						Expression: oltypes.String(".*"),
						Value:      []string{"test"},
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
	t.Run("It flattens the AppParameters Struct", func(t *testing.T) {
		appRuleStruct := []apps.AppRule{
			apps.AppRule{
				ID:       oltypes.Int32(int32(123)),
				Name:     oltypes.String("test"),
				Match:    oltypes.String("test"),
				Enabled:  oltypes.Bool(true),
				Position: oltypes.Int32(int32(1)),
				Conditions: []apps.AppRuleConditions{
					apps.AppRuleConditions{
						Source:   oltypes.String("test"),
						Operator: oltypes.String("="),
						Value:    oltypes.String("test"),
					},
				},
				Actions: []apps.AppRuleActions{
					apps.AppRuleActions{
						Action:     oltypes.String("test"),
						Expression: oltypes.String(".*"),
						Value:      []string{"test"},
					},
				},
			},
			apps.AppRule{
				ID:       oltypes.Int32(int32(456)),
				Name:     oltypes.String("test2"),
				Match:    oltypes.String("test2"),
				Enabled:  oltypes.Bool(true),
				Position: oltypes.Int32(int32(2)),
				Conditions: []apps.AppRuleConditions{
					apps.AppRuleConditions{
						Source:   oltypes.String("test2"),
						Operator: oltypes.String(">"),
						Value:    oltypes.String("test2"),
					},
				},
				Actions: []apps.AppRuleActions{
					apps.AppRuleActions{
						Action:     oltypes.String("test2"),
						Expression: oltypes.String(".*"),
						Value:      []string{"test2"},
					},
				},
			},
		}
		subj := Flatten(appRuleStruct)
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
						"action":     oltypes.String("test"),
						"expression": oltypes.String(".*"),
						"value":      []string{"test"},
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
						"action":     oltypes.String("test2"),
						"expression": oltypes.String(".*"),
						"value":      []string{"test2"},
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
