package appschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"

	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/rules"
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
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeBool},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: appparametersschema.Schema(),
			},
		},
		"rules": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: apprulesschema.Schema(),
			},
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin App.
func Inflate(s map[string]interface{}) (apps.App, error) {
	var err error
	app := apps.App{
		Name:               oltypes.String(s["name"].(string)),
		Description:        oltypes.String(s["description"].(string)),
		Notes:              oltypes.String(s["notes"].(string)),
		ConnectorID:        oltypes.Int32(int32(s["connector_id"].(int))),
		Visible:            oltypes.Bool(s["visible"].(bool)),
		AllowAssumedSignin: oltypes.Bool(s["allow_assumed_signin"].(bool)),
	}
	if s["parameters"] != nil {
		p := s["parameters"].(*schema.Set).List()
		app.Parameters = make(map[string]apps.AppParameters, len(p))
		for _, val := range p {
			valMap := val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = appparametersschema.Inflate(valMap)
		}
	}
	if s["provisioning"] != nil {
		prov := appprovisioningschema.Inflate(s["provisioning"].(map[string]interface{}))
		app.Provisioning = &prov
	}
	if s["configuration"] != nil {
		var conf apps.AppConfiguration
		conf, err = appconfigurationschema.Inflate(s["configuration"].(map[string]interface{}))
		app.Configuration = &conf
	}
	if s["rules"] != nil {
		appRules := make([]apps.AppRule, len(s["rules"].([]interface{})))
		for i, val := range s["rules"].([]interface{}) {
			valMap := val.(map[string]interface{})
			appRules[i] = apprulesschema.Inflate(valMap)
			appRules[i].Position = oltypes.Int32(int32(i + 1))
		}
		app.Rules = appRules
	}
	return app, err
}
