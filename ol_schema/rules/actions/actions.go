package appruleactionsschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps/app_rules"
)

// Schema returns a key/value map of the various fields that make up the Actions of a OneLogin Rule.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"expression": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) apprules.AppRuleActions {
	out := apprules.AppRuleActions{}
	if act, notNil := s["action"].(string); notNil {
		if act == "set_role_from_existing" {
			act = "set_role"
			out.Expression = nil
		} else {
			if exp, notNil := s["expression"].(string); notNil {
				out.Expression = oltypes.String(exp)
			}
		}
		out.Action = oltypes.String(act)
	}
	if val, notNil := s["value"].([]interface{}); notNil {
		out.Value = make([]string, len(val))
		for i, str := range val {
			out.Value[i] = str.(string)
		}
	}
	return out
}

// Flatten takes a AppRuleActions instance and converts it to an array of maps
func Flatten(acts []apprules.AppRuleActions) []map[string]interface{} {
	out := make([]map[string]interface{}, len(acts))
	for i, action := range acts {
		if action.Expression == nil && *action.Action == "set_role" {
			out[i] = map[string]interface{}{
				"action":     "set_role_from_existing",
				"expression": action.Expression,
				"value":      action.Value,
			}
		} else {
			out[i] = map[string]interface{}{
				"action":     action.Action,
				"expression": action.Expression,
				"value":      action.Value,
			}
		}
	}
	return out
}
