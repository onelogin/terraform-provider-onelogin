package appprovisioningschema

import (
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
)

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) apps.AppProvisioning {
	out := apps.AppProvisioning{}
	if enb, notNil := s["enabled"].(bool); notNil {
		out.Enabled = oltypes.Bool(enb)
	}
	return out
}

// Flatten takes a AppProvisioning instance and converts it to an array of maps
func Flatten(prov apps.AppProvisioning) map[string]interface{} {
	return map[string]interface{}{
		"enabled": *prov.Enabled,
	}
}
