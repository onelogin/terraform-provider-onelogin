package appruleactionsschema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
		"scriplet": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"macro": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a Action struct, a sub-field of a OneLogin Rule.
func Inflate(s map[string]interface{}) models.Action {
	out := models.Action{}
	if act, notNil := s["action"].(string); notNil {
		if strings.HasSuffix(act, NO_EXPRESSION_SUFFIX) {
			act = strings.TrimSuffix(act, NO_EXPRESSION_SUFFIX)
			// Expression is empty by default
		} else {
			if exp, notNil := s["expression"].(string); notNil {
				out.Expression = exp
			}
		}
		out.Action = act
	}
	if s["value"] != nil {
		v := s["value"].(*schema.Set).List()
		out.Value = make([]string, len(v))
		for i, val := range v {
			out.Value[i] = val.(string)
		}
	}
	if scriplet, notNil := s["scriplet"].(string); notNil {
		out.Scriplet = scriplet
	}
	if macro, notNil := s["macro"].(string); notNil {
		out.Macro = macro
	}
	return out
}

// Flatten takes a Action instance and converts it to an array of maps
func Flatten(acts []models.Action) []map[string]interface{} {
	out := make([]map[string]interface{}, len(acts))
	for i, action := range acts {
		if action.Expression == "" {
			out[i] = map[string]interface{}{
				"action":   fmt.Sprintf("%s%s", action.Action, NO_EXPRESSION_SUFFIX),
				"value":    action.Value,
				"scriplet": action.Scriplet,
				"macro":    action.Macro,
			}
		} else {
			out[i] = map[string]interface{}{
				"action":     action.Action,
				"expression": action.Expression,
				"value":      action.Value,
				"scriplet":   action.Scriplet,
				"macro":      action.Macro,
			}
		}
	}
	return out
}
