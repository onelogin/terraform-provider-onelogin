package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/rules"
)

// Apps returns a resource with the CRUD methods and Terraform Schema defined
func Apps() *schema.Resource {
	return &schema.Resource{
		Create:   appCreate,
		Read:     appRead,
		Update:   appUpdate,
		Delete:   appDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   appschema.Schema(),
	}
}

// appCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an App with its sub-resources
func appCreate(d *schema.ResourceData, m interface{}) error {
	basicApp, _ := appschema.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"rules":                d.Get("rules"),
	})
	client := m.(*client.APIClient)
	appResp, err := client.Services.AppsV2.Create(&basicApp)
	if err != nil {
		if appResp.ID != nil {
			log.Println("[ERROR] There was a problem setting sub-resources!", err)
			d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
			return appRead(d, m)
		}
		log.Println("[ERROR] There was a problem creating the app!", err)
		return err
	}
	log.Printf("[CREATED] Created app with %d", *(appResp.ID))

	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	return appRead(d, m)
}

// appRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an App with its sub-resources
func appRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	aid, _ := strconv.Atoi(d.Id())
	app, err := client.Services.AppsV2.GetOne(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the app!")
		log.Println(err)
		return err
	}
	if app == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading app with %d", *(app.ID))

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
	d.Set("parameters", appparametersschema.Flatten(app.Parameters))
	d.Set("provisioning", appprovisioningschema.Flatten(*app.Provisioning))
	d.Set("rules", apprulesschema.Flatten(app.Rules))

	return nil
}

// appUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an App and its sub-resources
func appUpdate(d *schema.ResourceData, m interface{}) error {
	basicApp, _ := appschema.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"rules":                d.Get("rules"),
	})

	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	appResp, err := client.Services.AppsV2.Update(int32(aid), &basicApp)
	if err != nil {
		if appResp.ID != nil {
			log.Println("[ERROR] There was a problem setting sub-resources!", err)
			d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
			return appRead(d, m)
		}
		log.Println("[ERROR] There was a problem updating the app!", err)
		return err
	}
	if appResp == nil { // app must be deleted in api so remove from tf state
		d.SetId("")
		return nil
	}
	log.Printf("[UPDATED] Updated app with %d", *(appResp.ID))
	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	return appRead(d, m)
}

// appDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an App and its sub-resources
func appDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	err := client.Services.AppsV2.Destroy(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the app!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted app with %d", aid)
		d.SetId("")
	}

	return nil
}
