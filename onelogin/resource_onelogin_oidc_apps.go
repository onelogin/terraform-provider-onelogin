package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/configuration"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/parameters"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/provisioning"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/sso"
)

func OneloginOIDCApps() *schema.Resource {
	appSchema := app.AppSchema()
	app.AddSubSchema("configuration", &appSchema, configuration.OIDCConfigurationSchema)
	app.AddSubSchema("sso", &appSchema, sso.OIDCSSOSchema)

	return &schema.Resource{
		Create: oidcAppCreate,
		Read:   oidcAppRead,
		Update: oidcAppUpdate,
		Delete: oidcAppDelete,
		Schema: appSchema,
	}
}

// oidcAppCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an oidcApp with its sub-resources
func oidcAppCreate(d *schema.ResourceData, m interface{}) error {
	appData := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
	}

	oidcApp := app.InflateApp(&appData)

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		oidcApp.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap := val.(map[string]interface{})
			oidcApp.Parameters[valMap["param_key_name"].(string)] = parameters.InflateParameter(&valMap)
		}
	}

	for _, val := range d.Get("provisioning").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		oidcApp.Provisioning = provisioning.InflateProvisioning(&valMap)
	}

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		oidcApp.Configuration = configuration.InflateOIDCConfiguration(&valMap)
	}

	client := m.(*client.APIClient)
	resp, oidcAppResp, err := client.Services.AppsV2.CreateApp(&oidcApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the oidcApp!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created oidcApp with %d", *(oidcAppResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(oidcAppResp.ID)))
	return oidcAppRead(d, m)
}

// oidcAppRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an oidcApp with its sub-resources
func oidcAppRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	aid, _ := strconv.Atoi(d.Id())
	resp, oidcApp, err := client.Services.AppsV2.GetAppByID(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[READ] Reading app with %d", *(oidcApp.ID))
	log.Println(resp)

	d.Set("name", oidcApp.Name)
	d.Set("visible", oidcApp.Visible)
	d.Set("description", oidcApp.Description)
	d.Set("notes", oidcApp.Notes)
	d.Set("icon_url", oidcApp.IconURL)
	d.Set("auth_method", oidcApp.AuthMethod)
	d.Set("policy_id", oidcApp.PolicyID)
	d.Set("allow_assumed_signin", oidcApp.AllowAssumedSignin)
	d.Set("tab_id", oidcApp.TabID)
	d.Set("connector_id", oidcApp.ConnectorID)
	d.Set("created_at", oidcApp.CreatedAt.String())
	d.Set("updated_at", oidcApp.UpdatedAt.String())
	d.Set("provisioning", provisioning.Flatten(oidcApp.Provisioning))
	d.Set("parameters", parameters.Flatten(oidcApp.Parameters))
	d.Set("configuration", configuration.FlattenOIDCConfiguration(oidcApp.Configuration))
	return nil
}

// oidcAppUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an oidcApp and its sub-resources
func oidcAppUpdate(d *schema.ResourceData, m interface{}) error {
	appData := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
	}

	oidcApp := app.InflateApp(&appData)

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		oidcApp.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap := val.(map[string]interface{})
			oidcApp.Parameters[valMap["param_key_name"].(string)] = parameters.InflateParameter(&valMap)
		}
	}

	for _, val := range d.Get("provisioning").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		oidcApp.Provisioning = provisioning.InflateProvisioning(&valMap)
	}

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		oidcApp.Configuration = configuration.InflateOIDCConfiguration(&valMap)
	}

	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, oidcAppResp, err := client.Services.AppsV2.UpdateAppByID(int32(aid), &oidcApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the oidcApp!")
		log.Println(err)
		return err
	}
	log.Printf("[UPDATED] Updated oidcApp with %d", *(oidcAppResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(oidcAppResp.ID)))
	return oidcAppRead(d, m)
}

// oidcAppDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an oidcApp and its sub-resources
func oidcAppDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, err := client.Services.AppsV2.DeleteApp(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the oidcApp!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted oidcApp with %d", aid)
		log.Println(resp)
		d.SetId("")
	}

	return nil
}
