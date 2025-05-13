package usermappingconditionsschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
// a UserMappingConditions struct, a sub-field of a OneLogin Rule.
func Inflate(s map[string]interface{}) models.UserMappingConditions {
	out := models.UserMappingConditions{}
	if src, notNil := s["source"].(string); notNil {
		out.Source = &src
	}
	if op, notNil := s["operator"].(string); notNil {
		out.Operator = &op
	}
	if val, notNil := s["value"].(string); notNil {
		out.Value = &val
	}
	return out
}

// Flatten takes a UserMappingConditions instance and converts it to an array of maps
func Flatten(conds []models.UserMappingConditions) []map[string]interface{} {
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
