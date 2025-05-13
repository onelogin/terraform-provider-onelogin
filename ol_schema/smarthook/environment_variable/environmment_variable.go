package smarthookenvironmentvariablesschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up the Environment Variables for a OneLogin SmartHook.
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
// an EnvVar struct, used by the OneLogin SmartHooks API.
func Inflate(s map[string]interface{}) models.EnvVar {
	out := models.EnvVar{}
	if id, notNil := s["id"].(string); notNil {
		out.ID = &id
	}
	if name, notNil := s["name"].(string); notNil {
		out.Name = &name
	}
	if value, notNil := s["value"].(string); notNil {
		out.Value = &value
	}
	return out
}
