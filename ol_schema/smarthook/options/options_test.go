package smarthookoptions

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	"github.com/stretchr/testify/assert"
)

func TestOptionsSchema(t *testing.T) {
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
		ExpectedOutput smarthooks.Options
	}{
		"creates and returns the address of a Options struct": {
			ResourceData: map[string]interface{}{
				"risk_enabled":            true,
				"mfa_device_info_enabled": true,
			},
			ExpectedOutput: smarthooks.Options{
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
		Input          smarthooks.Options
		ExpectedOutput map[string]interface{}
	}{
		"converts an instance of Options to a map interfaces with string keys": {
			Input: smarthooks.Options{
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
			subj := Flatten(test.Input)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
