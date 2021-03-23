package smarthookoptions

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
)

// Schema returns a key/value map of the various fields that make up
// the Parameters field for a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"risk_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"mfa_device_info_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"location_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func Inflate(s map[string]interface{}) smarthooks.Options {
	opts := smarthooks.Options{}

	if re, notNil := s["risk_enabled"].(bool); notNil {
		opts.RiskEnabled = oltypes.Bool(re)
	}
	if mdie, notNil := s["mfa_device_info_enabled"].(bool); notNil {
		opts.MFADeviceInfoEnabled = oltypes.Bool(mdie)
	}
	if le, notNil := s["location_enabled"].(bool); notNil {
		opts.LocationEnabled = oltypes.Bool(le)
	}
	return opts
}

// Flatten takes a SmartHook Options instance and creates a map
func Flatten(smarthookOptions smarthooks.Options) map[string]interface{} {
	return map[string]interface{}{
		"risk_enabled":            smarthookOptions.RiskEnabled,
		"mfa_device_info_enabled": smarthookOptions.MFADeviceInfoEnabled,
		"location_enabled":        smarthookOptions.LocationEnabled,
	}
}
