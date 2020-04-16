package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/parameters"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/provisioning"
)

func OneloginApps() *schema.Resource {
	return &schema.Resource{
		Create: appCreate,
		Read:   appRead,
		Update: appUpdate,
		Delete: appDelete,
		Schema: app.AppSchema(),
	}
}

// appCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an App with its sub-resources
func appCreate(d *schema.ResourceData, m interface{}) error {
	appData := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
	}

	app := app.InflateApp(&appData)

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		app.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap := val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = parameters.InflateParameter(&valMap)
		}
	}

	for _, val := range d.Get("provisioning").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		app.Provisioning = provisioning.InflateProvisioning(&valMap)
	}

	client := m.(*client.APIClient)
	resp, appResp, err := client.Services.AppsV2.CreateApp(&app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created app with %d", *(appResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	// return appRead(d, m)
	return nil
}

// appRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an App with its sub-resources
func appRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	aid, _ := strconv.Atoi(d.Id())
	resp, app, err := client.Services.AppsV2.GetAppByID(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[READ] Reading app with %d", *(app.ID))
	log.Println(resp)

	d.Set("name", app.Name)
	d.Set("visible", app.Visible)
	d.Set("description", app.Description)
	d.Set("notes", app.Notes)
	d.Set("icon_url", app.IconURL)
	d.Set("auth_method", app.AuthMethod)
	d.Set("policy_id", app.PolicyID)
	d.Set("allow_assumed_signin", app.AllowAssumedSignin)
	d.Set("tab_id", app.TabID)
	d.Set("connector_id", app.ConnectorID)
	d.Set("created_at", app.CreatedAt.String())
	d.Set("updated_at", app.UpdatedAt.String())
	d.Set("provisioning", provisioning.Flatten(app.Provisioning))
	d.Set("parameters", parameters.Flatten(app.Parameters))
	return nil
}

// appUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an App and its sub-resources
func appUpdate(d *schema.ResourceData, m interface{}) error {
	appData := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
	}

	app := app.InflateApp(&appData)

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		app.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap := val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = parameters.InflateParameter(&valMap)
		}
	}

	for _, val := range d.Get("provisioning").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		app.Provisioning = provisioning.InflateProvisioning(&valMap)
	}
	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, appResp, err := client.Services.AppsV2.UpdateAppByID(int32(aid), &app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[UPDATED] Updated app with %d", *(appResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	return appRead(d, m)
	// return nil
}

// appDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an App and its sub-resources
func appDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, err := client.Services.AppsV2.DeleteApp(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted app with %d", aid)
		log.Println(resp)
		d.SetId("")
	}

	return nil
}
