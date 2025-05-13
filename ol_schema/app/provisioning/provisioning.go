package appprovisioningschema

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a Provisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) models.Provisioning {
	out := models.Provisioning{}
	if enb, notNil := s["enabled"].(bool); notNil {
		out.Enabled = enb
	}
	return out
}

// Flatten takes a Provisioning instance and converts it to an array of maps
func Flatten(prov models.Provisioning) map[string]interface{} {
	return map[string]interface{}{
		"enabled": prov.Enabled,
	}
}

// FlattenMap takes a map[string]interface{} instance and converts it to a map
func FlattenMap(prov map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	if val, ok := prov["enabled"].(bool); ok {
		out["enabled"] = val
	}
	return out
}
