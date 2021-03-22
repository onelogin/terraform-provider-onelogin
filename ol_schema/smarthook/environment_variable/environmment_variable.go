package smarthookenvironmentvariablesschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks/envs"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) smarthookenvs.EnvVar {
	out := smarthookenvs.EnvVar{}
	if id, notNil := s["id"].(string); notNil {
		out.ID = oltypes.String(id)
	}
	if hookType, notNil := s["name"].(string); notNil {
		out.Name = oltypes.String(hookType)
	}
	if runtime, notNil := s["value"].(string); notNil {
		out.Value = oltypes.String(runtime)
	}
	return out
}
