package smarthookoptions

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	"github.com/stretchr/testify/assert"
)

func TestSmartHookOptionsSchema(t *testing.T) {
	t.Run("creates and returns a map of a Smarthook Options Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["risk_enabled"])
		assert.NotNil(t, provSchema["mfa_device_info_enabled"])
		assert.NotNil(t, provSchema["location_enabled"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput smarthooks.SmartHookOptions
	}{
		"creates and returns the address of a SmartHookOptions struct": {
			ResourceData: map[string]interface{}{
				"risk_enabled":            true,
				"mfa_device_info_enabled": true,
			},
			ExpectedOutput: smarthooks.SmartHookOptions{
				RiskEnabled:          oltypes.Bool(true),
				MFADeviceInfoEnabled: oltypes.Bool(true),
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
	tests := map[string]struct {
		Input          smarthooks.SmartHookOptions
		ExpectedOutput map[string]interface{}
	}{
		"converts an instance of SmartHookOptions to a map interfaces with string keys": {
			Input: smarthooks.SmartHookOptions{
				RiskEnabled:          oltypes.Bool(true),
				MFADeviceInfoEnabled: oltypes.Bool(true),
				LocationEnabled:      oltypes.Bool(true),
			},
			ExpectedOutput: map[string]interface{}{
				"risk_enabled":            oltypes.Bool(true),
				"mfa_device_info_enabled": oltypes.Bool(true),
				"location_enabled":        oltypes.Bool(true),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSmartHookOptions(test.Input)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
