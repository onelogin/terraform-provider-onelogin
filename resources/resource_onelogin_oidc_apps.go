package resources

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/configuration"
)

func OneloginOIDCApps() *schema.Resource {
	appSchema := app.AppSchema()
	configuration.AddConfigurationSchema(&appSchema, configuration.OIDCConfigurationSchema)
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
	oidcApp := app.InflateApp(d)

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		oidcApp.Configuration = configuration.InflateOIDCConfiguration(&valMap)
	}

	client := m.(*client.APIClient)
	resp, oidcApp, err := client.Services.AppsV2.CreateApp(oidcApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the oidcApp!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created oidcApp with %d", *(oidcApp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(oidcApp.ID)))
	return oidcAppRead(d, m)
}

// oidcAppRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an oidcApp with its sub-resources
func oidcAppRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// oidcAppUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an oidcApp and its sub-resources
func oidcAppUpdate(d *schema.ResourceData, m interface{}) error {
	oidcApp := app.InflateApp(d)

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		oidcApp.Configuration = configuration.InflateOIDCConfiguration(&valMap)
	}

	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, oidcApp, err := client.Services.AppsV2.UpdateAppByID(int32(aid), oidcApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the oidcApp!")
		log.Println(err)
		return err
	}
	log.Printf("[UPDATED] Updated oidcApp with %d", *(oidcApp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(oidcApp.ID)))
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
