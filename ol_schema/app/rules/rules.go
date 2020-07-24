package apprulesschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/rules/actions"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/rules/conditions"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
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
			Computed: true,
		},
		"conditions": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: appruleconditionsschema.Schema(),
			},
		},
		"actions": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: appruleactionsschema.Schema(),
			},
		},
	}
}

func validMatch(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"all", "any"})
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) apps.AppRule {
	out := apps.AppRule{}
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
			cond := appruleconditionsschema.Inflate(valMap)
			out.Conditions = append(out.Conditions, cond)
		}
	}
	if s["actions"] != nil {
		for _, val := range s["actions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := appruleactionsschema.Inflate(valMap)
			out.Actions = append(out.Actions, cond)
		}
	}
	return out
}

// Flatten takes a AppRules array and converts it to an array of maps
func Flatten(appRules []apps.AppRule) []map[string]interface{} {
	out := make([]map[string]interface{}, len(appRules))
	for i, appRule := range appRules {
		out[i] = map[string]interface{}{
			"id":         appRule.ID,
			"name":       appRule.Name,
			"match":      appRule.Match,
			"enabled":    appRule.Enabled,
			"position":   appRule.Position,
			"conditions": appruleconditionsschema.Flatten(appRule.Conditions),
			"actions":    appruleactionsschema.Flatten(appRule.Actions),
		}
	}
	return out
}
