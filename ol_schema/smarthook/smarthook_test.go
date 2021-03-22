package smarthooksschema

import (
	"fmt"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks/envs"
	"github.com/stretchr/testify/assert"
)

func TestSmartHookSchema(t *testing.T) {
	t.Run("creates and returns a map of a Smarthooks Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["type"])
		assert.NotNil(t, provSchema["status"])
		assert.NotNil(t, provSchema["disabled"])
		assert.NotNil(t, provSchema["runtime"])
		assert.NotNil(t, provSchema["retries"])
		assert.NotNil(t, provSchema["timeout"])
		assert.NotNil(t, provSchema["packages"])
		assert.NotNil(t, provSchema["env_vars"])
		assert.NotNil(t, provSchema["options"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput smarthooks.SmartHook
	}{
		"creates and returns the address of a SmartHook": {
			ResourceData: map[string]interface{}{
				"id":       "32f9dfee-a02c-4932-98ec-37838ce62ba0",
				"type":     "pre-authentication",
				"function": "function myFunc(){...}",
				"packages": map[string]interface{}{"mysql": "^2.18.1"},
				"retries":  0,
				"timeout":  2,
				"disabled": false,
				"env_vars": []interface{}{"API_KEY"},
				"options": map[string]interface{}{
					"risk_enabled": false,
				},
			},
			ExpectedOutput: smarthooks.SmartHook{
				ID:       oltypes.String("32f9dfee-a02c-4932-98ec-37838ce62ba0"),
				Type:     oltypes.String("pre-authentication"),
				Function: oltypes.String("function myFunc(){...}"),
				Packages: map[string]string{"mysql": "^2.18.1"},
				Retries:  oltypes.Int32(int32(0)),
				Timeout:  oltypes.Int32(int32(2)),
				Disabled: oltypes.Bool(false),
				EnvVars:  []smarthookenvs.EnvVar{smarthookenvs.EnvVar{Name: oltypes.String("API_KEY")}},
				Options: &smarthooks.SmartHookOptions{
					RiskEnabled: oltypes.Bool(false),
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

func TestValidTypes(t *testing.T) {
	tests := map[string]struct {
		InputKey       string
		InputValue     string
		ExpectedOutput []error
	}{
		"no errors on valid input": {
			InputKey:       "type",
			InputValue:     "pre-authentication",
			ExpectedOutput: nil,
		},
		"errors on invalid input": {
			InputKey:       "type",
			InputValue:     "asdf",
			ExpectedOutput: []error{fmt.Errorf("type must be one of [pre-authentication user-migration], got: asdf")},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, errs := validTypes(test.InputValue, test.InputKey)
			assert.Equal(t, test.ExpectedOutput, errs)
		})
	}
}
