package apprulesschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	appruleactionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/actions"
	appruleconditionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/conditions"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
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
				Schema: appruleconditionsschema.Schema(),
			},
		},
		"actions": {
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
// an AppRule struct for a OneLogin App.
func Inflate(s map[string]interface{}) models.AppRule {
	out := models.AppRule{}

	// Store rule ID in a variable for later use if needed
	// But it's not part of the AppRule struct
	if s["id"] != nil {
		if _, err := strconv.Atoi(s["id"].(string)); err == nil {
			// Rule ID is not directly stored in the AppRule struct
		}
	}

	if s["app_id"] != nil {
		if id, err := strconv.Atoi(s["app_id"].(string)); err == nil {
			out.AppID = id
		}
	}

	if n, notNil := s["name"].(string); notNil {
		out.Name = n
	}

	if m, notNil := s["match"].(string); notNil {
		out.Match = m
	}

	if pos, notNil := s["position"].(int); notNil {
		out.Position = pos
	}

	if en, notNil := s["enabled"].(bool); notNil {
		out.Enabled = en
	}

	if s["conditions"] != nil {
		out.Conditions = []models.Condition{}
		for _, val := range s["conditions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := appruleconditionsschema.Inflate(valMap)
			out.Conditions = append(out.Conditions, cond)
		}
	}

	if s["actions"] != nil {
		out.Actions = []models.Action{}
		for _, val := range s["actions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			action := appruleactionsschema.Inflate(valMap)
			out.Actions = append(out.Actions, action)
		}
	}

	return out
}
