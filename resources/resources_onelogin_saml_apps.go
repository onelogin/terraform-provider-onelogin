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

func OneloginSAMLApps() *schema.Resource {
	appSchema := app.AppSchema()
	configuration.AddConfigurationSchema(&appSchema, configuration.SAMLConfigurationSchema)
	return &schema.Resource{
		Create: samlAppCreate,
		Read:   samlAppRead,
		Update: samlAppUpdate,
		Delete: samlAppDelete,
		Schema: appSchema,
	}
}

// samlAppCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an samlApp with its sub-resources
func samlAppCreate(d *schema.ResourceData, m interface{}) error {
	samlApp := app.InflateApp(d)

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		samlApp.Configuration = configuration.InflateSAMLConfiguration(&valMap)
	}

	client := m.(*client.APIClient)
	resp, samlApp, err := client.Services.AppsV2.CreateApp(samlApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the samlApp!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created samlApp with %d", *(samlApp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(samlApp.ID)))
	return samlAppRead(d, m)
}

// samlAppRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an samlApp with its sub-resources
func samlAppRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// samlAppUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an samlApp and its sub-resources
func samlAppUpdate(d *schema.ResourceData, m interface{}) error {
	samlApp := app.InflateApp(d)

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		samlApp.Configuration = configuration.InflateSAMLConfiguration(&valMap)
	}

	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, samlApp, err := client.Services.AppsV2.UpdateAppByID(int32(aid), samlApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the samlApp!")
		log.Println(err)
		return err
	}
	log.Printf("[UPDATED] Updated samlApp with %d", *(samlApp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(samlApp.ID)))
	return samlAppRead(d, m)
}

// samlAppDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an samlApp and its sub-resources
func samlAppDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, err := client.Services.AppsV2.DeleteApp(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the samlApp!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted samlApp with %d", aid)
		log.Println(resp)
		d.SetId("")
	}

	return nil
}
