package resources

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-terraform-provider/resources/app"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/configuration"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/parameters"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/provisioning"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/sso"
)

func OneloginSAMLApps() *schema.Resource {
	appSchema := app.AppSchema()
	app.AddSubSchema("configuration", &appSchema, configuration.SAMLConfigurationSchema)
	app.AddSubSchema("sso", &appSchema, sso.SAMLSSOSchema)

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
	appData := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
	}

	samlApp := app.InflateApp(&appData)

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		samlApp.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap := val.(map[string]interface{})
			samlApp.Parameters[valMap["param_key_name"].(string)] = parameters.InflateParameter(&valMap)
		}
	}

	for _, val := range d.Get("provisioning").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		samlApp.Provisioning = provisioning.InflateProvisioning(&valMap)
	}

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		samlApp.Configuration = configuration.InflateSAMLConfiguration(&valMap)
	}

	client := m.(*client.APIClient)
	resp, samlAppResp, err := client.Services.AppsV2.CreateApp(&samlApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the samlApp!")
		log.Println(err)
		return err
	}
	log.Printf("[CREATED] Created samlApp with %d", *(samlAppResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(samlAppResp.ID)))
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
	appData := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
	}

	samlApp := app.InflateApp(&appData)

	if paramsList, isSet := d.GetOk("parameters"); isSet {
		samlApp.Parameters = make(map[string]models.AppParameters, len(paramsList.(*schema.Set).List()))
		for _, val := range paramsList.(*schema.Set).List() {
			valMap := val.(map[string]interface{})
			samlApp.Parameters[valMap["param_key_name"].(string)] = parameters.InflateParameter(&valMap)
		}
	}

	for _, val := range d.Get("provisioning").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		samlApp.Provisioning = provisioning.InflateProvisioning(&valMap)
	}

	for _, val := range d.Get("configuration").(*schema.Set).List() {
		valMap := val.(map[string]interface{})
		samlApp.Configuration = configuration.InflateSAMLConfiguration(&valMap)
	}

	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, samlAppResp, err := client.Services.AppsV2.UpdateAppByID(int32(aid), &samlApp)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the samlApp!")
		log.Println(err)
		return err
	}
	log.Printf("[UPDATED] Updated samlApp with %d", *(samlAppResp.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(samlAppResp.ID)))
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
