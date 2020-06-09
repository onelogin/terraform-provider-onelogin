package appruleconditionsschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
)

// Schema returns a key/value map of the various fields that make up the Actions of a OneLogin Rule.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"operator": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppRuleConditions struct, a sub-field of a OneLogin Rule.
func Inflate(s map[string]interface{}) apps.AppRuleConditions {
	out := apps.AppRuleConditions{}
	if enb, notNil := s["source"].(string); notNil {
		out.Source = oltypes.String(enb)
	}
	if enb, notNil := s["operator"].(string); notNil {
		out.Operator = oltypes.String(enb)
	}
	if enb, notNil := s["value"].(string); notNil {
		out.Value = oltypes.String(enb)
	}
	return out
}

// Flatten takes a AppProvisioning instance and converts it to an array of maps
func Flatten(conds []apps.AppRuleConditions) []map[string]interface{} {
	out := make([]map[string]interface{}, len(conds))
	for i, condition := range conds {
		out[i] = map[string]interface{}{
			"source":   condition.Source,
			"operator": condition.Operator,
			"value":    condition.Value,
		}
	}
	return out
}
