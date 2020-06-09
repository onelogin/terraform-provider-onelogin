package usermappingschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/user_mappings"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/actions"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/conditions"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"match": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validMatch,
		},
		"enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"position": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"conditions": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: usermappingconditionsschema.Schema(),
			},
		},
		"actions": &schema.Schema{
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
func Inflate(s map[string]interface{}) usermappings.UserMapping {
	out := usermappings.UserMapping{}
	if id, notNil := s["id"].(int); id != 0 && notNil {
		out.ID = oltypes.Int32(int32(id))
	}
	if name, notNil := s["name"].(string); notNil {
		out.Name = oltypes.String(name)
	}
	if match, notNil := s["match"].(string); notNil {
		out.Match = oltypes.String(match)
	}
	if pos, notNil := s["position"].(int); notNil {
		out.Position = oltypes.Int32(int32(pos))
	}
	if pos, notNil := s["enabled"].(bool); notNil {
		out.Enabled = oltypes.Bool(pos)
	}
	if s["conditions"] != nil {
		for _, val := range s["conditions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := usermappingconditionsschema.Inflate(valMap)
			out.Conditions = append(out.Conditions, cond)
		}
	}
	if s["actions"] != nil {
		for _, val := range s["actions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := usermappingactionsschema.Inflate(valMap)
			out.Actions = append(out.Actions, cond)
		}
	}
	return out
}

// Flatten takes a UserMappings array and converts it to an array of maps
func Flatten(UserMappings []usermappings.UserMapping) []map[string]interface{} {
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
