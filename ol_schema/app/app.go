package app

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/onelogin-terraform-provider/ol_schema/app/configuration"
	"github.com/onelogin/onelogin-terraform-provider/ol_schema/app/parameters"
	"github.com/onelogin/onelogin-terraform-provider/ol_schema/app/provisioning"
	"github.com/onelogin/onelogin-terraform-provider/ol_schema/shared/rules"
)

// Schema returns a key/value map of the various fields that make up an App at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"visible": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"policy_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"allow_assumed_signin": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tab_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"connector_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: provisioning.Schema(),
			},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: parameters.Schema(),
			},
		},
		"rules": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rules.Schema(),
			},
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin App.
func Inflate(s map[string]interface{}) apps.App {
	app := apps.App{
		Name:               oltypes.String(s["name"].(string)),
		Description:        oltypes.String(s["description"].(string)),
		Notes:              oltypes.String(s["notes"].(string)),
		ConnectorID:        oltypes.Int32(int32(s["connector_id"].(int))),
		Visible:            oltypes.Bool(s["visible"].(bool)),
		AllowAssumedSignin: oltypes.Bool(s["allow_assumed_signin"].(bool)),
	}
	if s["parameters"] != nil {
		app.Parameters = make(map[string]apps.AppParameters, len(s["parameters"].([]interface{})))
		for _, val := range s["parameters"].([]interface{}) {
			valMap := val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = parameters.Inflate(valMap)
		}
	}
	if s["provisioning"] != nil {
		for _, val := range s["provisioning"].([]interface{}) {
			valMap := val.(map[string]interface{})
			prov := provisioning.Inflate(valMap)
			app.Provisioning = &prov
		}
	}
	if s["configuration"] != nil {
		for _, val := range s["configuration"].([]interface{}) {
			valMap := val.(map[string]interface{})
			config := configuration.Inflate(valMap)
			app.Configuration = &config
		}
	}
	if s["rules"] != nil {
		appRules := make([]apps.AppRule, len(s["rules"].([]interface{})))
		for i, val := range s["rules"].([]interface{}) {
			valMap := val.(map[string]interface{})
			appRules[i] = rules.Inflate(valMap)
		}
		app.Rules = appRules
	}
	return app
}
