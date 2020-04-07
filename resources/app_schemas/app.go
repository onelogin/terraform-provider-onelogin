package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
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
	nam := d.Get("name").(string)
	des := d.Get("description").(string)
	not := d.Get("notes").(string)
	iur := d.Get("icon_url").(string)

	app := models.App{
		Name:        &nam,
		Description: &des,
		Notes:       &not,
		IconURL:     &iur,
	}

	if paramsList, paramsGiven := d.GetOk("parameters"); paramsGiven {
		app.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, s := range paramsList.(*schema.Set).List() {
			sMap := s.(map[string]interface{})
			key := sMap["param_key_name"].(string)
			app.Parameters[key] = InflateAppParameter(sMap)
		}
	}

	for _, s := range d.Get("provisioning").(*schema.Set).List() {
		app.Provisioning = InflateAppProvisioning(s.(map[string]interface{}))
	}

	for _, s := range d.Get("configuration").(*schema.Set).List() {
		app.Configuration = InflateAppConfiguration(s.(map[string]interface{}))
	}

	if vis, visSet := d.GetOk("visible"); visSet {
		vis := vis.(bool)
		app.Visible = &vis
	}

	if aas, aasSet := d.GetOk("allow_assumed_signin"); aasSet {
		aas := aas.(bool)
		app.AllowAssumedSignin = &aas
	}

	if cid, cidSet := d.GetOk("connector_id"); cidSet {
		cid := int32(cid.(int))
		app.ConnectorID = &cid
	}

	if aum, aumSet := d.GetOk("auth_method"); aumSet {
		aum := int32(aum.(int))
		app.AuthMethod = &aum
	}

	if pid, pidSet := d.GetOk("policy_id"); pidSet {
		pid := int32(pid.(int))
		app.PolicyID = &pid
	}

	if tid, tidSet := d.GetOk("tab_id"); tidSet {
		tid := int32(tid.(int))
		app.TabID = &tid
	}

	return &app
}
