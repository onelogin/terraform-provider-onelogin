package usermappingschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	usermappingactionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/actions"
	usermappingconditionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/conditions"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin UserMapping.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"match": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validMatch,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"position": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"conditions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: usermappingconditionsschema.Schema(),
			},
		},
		"actions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: usermappingactionsschema.Schema(),
			},
		},
	}
}

func validMatch(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"all", "any"})
}

// Inflate takes a key/value map of interfaces and uses the fields to construct a user mapping struct
func Inflate(s map[string]interface{}) (models.UserMapping, error) {
	out := models.UserMapping{}
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			id32 := int32(id)
			out.ID = &id32
		}
	}
	if name, notNil := s["name"].(string); notNil {
		out.Name = &name
	}
	if match, notNil := s["match"].(string); notNil {
		out.Match = &match
	}
	if pos, notNil := s["position"].(int); notNil {
		pos32 := int32(pos)
		out.Position = &pos32
	}
	if enabled, notNil := s["enabled"].(bool); notNil {
		out.Enabled = &enabled
	}
	if s["conditions"] != nil {
		out.Conditions = []models.UserMappingConditions{}
		for _, val := range s["conditions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := usermappingconditionsschema.Inflate(valMap)
			out.Conditions = append(out.Conditions, cond)
		}
	}
	if s["actions"] != nil {
		out.Actions = []models.UserMappingActions{}
		for _, val := range s["actions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			action := usermappingactionsschema.Inflate(valMap)
			out.Actions = append(out.Actions, action)
		}
	}
	return out, nil
}

// Flatten takes a UserMappings array and converts it to an array of maps
func Flatten(UserMappings []models.UserMapping) []map[string]interface{} {
	out := make([]map[string]interface{}, len(UserMappings))
	for i, UserMapping := range UserMappings {
		out[i] = map[string]interface{}{
			"id":         UserMapping.ID,
			"name":       UserMapping.Name,
			"match":      UserMapping.Match,
			"enabled":    UserMapping.Enabled,
			"position":   UserMapping.Position,
			"conditions": usermappingconditionsschema.Flatten(UserMapping.Conditions),
			"actions":    usermappingactionsschema.Flatten(UserMapping.Actions),
		}
	}
	return out
}

// FlattenConditions converts an array of interface{} to an array of maps for conditions
func FlattenConditions(conditions []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(conditions))
	for i, condition := range conditions {
		if condMap, ok := condition.(map[string]interface{}); ok {
			result[i] = map[string]interface{}{
				"source":   condMap["source"],
				"operator": condMap["operator"],
				"value":    condMap["value"],
			}
		}
	}
	return result
}

// FlattenActions converts an array of interface{} to an array of maps for actions
func FlattenActions(actions []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(actions))
	for i, action := range actions {
		if actMap, ok := action.(map[string]interface{}); ok {
			var vals []interface{}
			if v, ok := actMap["value"].([]interface{}); ok {
				vals = v
			}
			result[i] = map[string]interface{}{
				"action": actMap["action"],
				"value":  vals,
			}
		}
	}
	return result
}
