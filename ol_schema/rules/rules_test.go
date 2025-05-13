package apprulesschema

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a Rules Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["app_id"])
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["match"])
		assert.NotNil(t, provSchema["position"])
		assert.NotNil(t, provSchema["conditions"])
		assert.NotNil(t, provSchema["actions"])
	})
}

func mockSetFn(interface{}) int {
	return 0
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.AppRule
	}{
		"creates and returns the address of a Rule struct": {
			ResourceData: map[string]interface{}{
				"id":       "123",
				"app_id":   "123",
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
						"value":      schema.NewSet(mockSetFn, []interface{}{"test"}),
					},
				},
			},
			ExpectedOutput: models.AppRule{
				AppID:    123,
				Name:     "test",
				Match:    "test",
				Enabled:  true,
				Position: 1,
				Conditions: []models.Condition{
					{
						Source:   "test",
						Operator: "=",
						Value:    "test",
					},
				},
				Actions: []models.Action{
					{
						Action:     "test",
						Expression: ".*",
						Value:      []string{"test"},
						Scriplet:   "",
						Macro:      "",
					},
				},
			},
		},
		"handles a rule without the position provided": {
			ResourceData: map[string]interface{}{
				"id":      "123",
				"app_id":  "123",
				"name":    "test",
				"match":   "test",
				"enabled": true,
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
						"value":      schema.NewSet(mockSetFn, []interface{}{"test"}),
					},
				},
			},
			ExpectedOutput: models.AppRule{
				AppID:    123,
				Name:     "test",
				Match:    "test",
				Enabled:  true,
				Position: 0,
				Conditions: []models.Condition{
					{
						Source:   "test",
						Operator: "=",
						Value:    "test",
					},
				},
				Actions: []models.Action{
					{
						Action:     "test",
						Expression: ".*",
						Value:      []string{"test"},
						Scriplet:   "",
						Macro:      "",
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
