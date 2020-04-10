package app

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
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: ProvisioningSchema(),
			},
		},
		"configuration": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: ConfigurationSchema(),
			},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: ParameterSchema(),
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
	}

	// if val, isSet = d.GetOkExists("description"); isSet {
	// 	log.Println("SET THE DESC FOOL!")
	// 	app.Description = oltypes.String(val.(string))
	// }

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		app.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap = val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = *InflateParameter(&valMap) // dereference appParameters due to field constraint on App struct
		}
	}

	for _, val = range d.Get("provisioning").(*schema.Set).List() {
		valMap = val.(map[string]interface{})
		app.Provisioning = InflateProvisioning(&valMap)
	}

	for _, val = range d.Get("configuration").(*schema.Set).List() {
		valMap = val.(map[string]interface{})
		app.Configuration = InflateConfiguration(&valMap)
	}

	if val, isSet = d.GetOkExists("visible"); isSet {
		app.Visible = oltypes.Bool(val.(bool))
	}

	if val, isSet = d.GetOkExists("allow_assumed_signin"); isSet {
		app.AllowAssumedSignin = oltypes.Bool(val.(bool))
	}

	if val, isSet = d.GetOkExists("connector_id"); isSet {
		app.ConnectorID = oltypes.Int32(int32(val.(int)))
	}

	return &app
}
