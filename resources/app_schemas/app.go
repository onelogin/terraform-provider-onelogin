package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// App returns a key/value map of the various fields that make up an App at OneLogin.
func App() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"visible": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"policy_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"allow_assumed_signin": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"tab_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"connector_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"provisioning": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: AppProvisioning(),
			},
		},
		"configuration": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: AppConfiguration(),
			},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: AppParameter(),
			},
		},
	}
}

// InflateApp takes a pointer to a ResourceData struct and uses it to construct a
// OneLogin App struct to be used in requests to OneLogin.
func InflateApp(d *schema.ResourceData) *models.App {
	var val interface{}
	var valMap map[string]interface{}
	var isSet bool

	app := models.App{
		Name:        oltypes.String(d.Get("name").(string)),
		Description: oltypes.String(d.Get("description").(string)),
		Notes:       oltypes.String(d.Get("notes").(string)),
		IconURL:     oltypes.String(d.Get("icon_url").(string)),
	}

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		app.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap = val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = *InflateAppParameter(&valMap) // dereference appParameters due to field constraint on App struct
		}
	}

	for _, val = range d.Get("provisioning").(*schema.Set).List() {
		valMap = val.(map[string]interface{})
		app.Provisioning = InflateAppProvisioning(&valMap)
	}

	for _, val = range d.Get("configuration").(*schema.Set).List() {
		valMap = val.(map[string]interface{})
		app.Configuration = InflateAppConfiguration(&valMap)
	}

	if val, isSet = d.GetOk("visible"); isSet {
		app.Visible = oltypes.Bool(val.(bool))
	}

	if val, isSet = d.GetOk("allow_assumed_signin"); isSet {
		app.AllowAssumedSignin = oltypes.Bool(val.(bool))
	}

	if val, isSet = d.GetOk("connector_id"); isSet {
		app.ConnectorID = oltypes.Int32(int32(val.(int)))
	}

	if val, isSet = d.GetOk("auth_method"); isSet {
		app.AuthMethod = oltypes.Int32(int32(val.(int)))
	}

	if val, isSet = d.GetOk("policy_id"); isSet {
		app.PolicyID = oltypes.Int32(int32(val.(int)))
	}

	if val, isSet = d.GetOk("tab_id"); isSet {
		app.TabID = oltypes.Int32(int32(val.(int)))
	}

	return &app
}
