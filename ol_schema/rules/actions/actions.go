package appruleactionsschema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	apprules "github.com/onelogin/onelogin-go-sdk/pkg/services/apps/app_rules"
)

const NO_EXPRESSION_SUFFIX = "_from_existing"

// Schema returns a key/value map of the various fields that make up the Actions of a OneLogin Rule.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:     schema.TypeString,
			Required: true,
		},
		"expression": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": {
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
		if strings.HasSuffix(act, NO_EXPRESSION_SUFFIX) {
			act = strings.TrimSuffix(act, NO_EXPRESSION_SUFFIX)
			out.Expression = nil
		} else {
			if exp, notNil := s["expression"].(string); notNil {
				out.Expression = oltypes.String(exp)
			}
		}
		out.Action = oltypes.String(act)
	}
	if s["value"] != nil {
		v := s["value"].(*schema.Set).List()
		out.Value = make([]string, len(v))
		for i, val := range v {
			out.Value[i] = val.(string)
		}
	}
	return out
}

// Flatten takes a AppRuleActions instance and converts it to an array of maps
func Flatten(acts []apprules.AppRuleActions) []map[string]interface{} {
	out := make([]map[string]interface{}, len(acts))
	for i, action := range acts {
		if action.Expression == nil && action.Action != nil {
			out[i] = map[string]interface{}{
				"action":     fmt.Sprintf("%s%s", *action.Action, NO_EXPRESSION_SUFFIX),
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
