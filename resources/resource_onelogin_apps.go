package resources

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
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
	app := app.InflateApp(d)

	client := m.(*client.APIClient)
	resp, app, err := client.Services.AppsV2.CreateApp(app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created app with %d", *(app.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(app.ID)))
	return appRead(d, m)
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

	d.SetId(fmt.Sprintf("%d", app.ID))

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
	return nil
}

// appUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an App and its sub-resources
func appUpdate(d *schema.ResourceData, m interface{}) error {
	app := app.InflateApp(d)
	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, app, err := client.Services.AppsV2.UpdateAppByID(int32(aid), app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
		return err
	}
	log.Printf("[UPDATED] Updated app with %d", *(app.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(app.ID)))
	return appRead(d, m)
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
