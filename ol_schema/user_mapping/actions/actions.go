package usermappingactionsschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up the Actions of a OneLogin Rule.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a UserMappingActions struct, a sub-field of a OneLogin UserMapping.
func Inflate(s map[string]interface{}) models.UserMappingActions {
	out := models.UserMappingActions{}
	if act, notNil := s["action"].(string); notNil {
		out.Action = &act
	}
	if val, notNil := s["value"].([]interface{}); notNil {
		out.Value = make([]string, len(val))
		for i, str := range val {
			out.Value[i] = str.(string)
		}
	}
	return out
}

// Flatten takes a UserMappingActions instance and converts it to an array of maps
func Flatten(acts []models.UserMappingActions) []map[string]interface{} {
	out := make([]map[string]interface{}, len(acts))
	for i, action := range acts {
		out[i] = map[string]interface{}{
			"action": action.Action,
			"value":  action.Value,
		}
	}
	return out
}
