package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/rules"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/app/sso"
)

// OIDCApps attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func OIDCApps() *schema.Resource {
	appSchema := appschema.Schema()
	appSchema["configuration"] = &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	appSchema["sso"] = &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}

	return &schema.Resource{
		Create:   oidcAppCreate,
		Read:     oidcAppRead,
		Update:   oidcAppUpdate,
		Delete:   oidcAppDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   appSchema,
	}
}

// oidcAppCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an oidcApp with its sub-resources
func oidcAppCreate(d *schema.ResourceData, m interface{}) error {
	oidcApp := appschema.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
		"rules":                d.Get("rules"),
	})
	if err != nil {
		log.Println("Unable to convert string in plan to required value type", err)
		return err
	}
	client := m.(*client.APIClient)
	appResp, err := client.Services.AppsV2.Create(&oidcApp)
	if err != nil {
		if appResp.ID != nil {
			log.Println("[ERROR] There was a problem setting sub-resources!", err)
			d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
			return oidcAppRead(d, m)
		}
		log.Println("[ERROR] There was a problem creating the app!", err)
		return err
	}
	log.Printf("[CREATED] Created app with %d", *(appResp.ID))

	d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
	return oidcAppRead(d, m)
}

// oidcAppRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an oidcApp with its sub-resources
func oidcAppRead(d *schema.ResourceData, m interface{}) error {
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
	d.Set("configuration", appconfigurationschema.FlattenOIDC(*app.Configuration))
	d.Set("sso", appssoschema.FlattenOIDC(*app.Sso))
	d.Set("rules", apprulesschema.Flatten(app.Rules))

	return nil
}

// oidcAppUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an oidcApp and its sub-resources
func oidcAppUpdate(d *schema.ResourceData, m interface{}) error {
	oidcApp := appschema.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
		"rules":                d.Get("rules"),
	})
	if err != nil {
		log.Println("Unable to convert string in plan to required value type", err)
		return err
	}
	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	appResp, err := client.Services.AppsV2.Update(int32(aid), &oidcApp)
	if err != nil {
		if appResp.ID != nil {
			log.Println("[ERROR] There was a problem setting sub-resources!", err)
			d.SetId(fmt.Sprintf("%d", *(appResp.ID)))
			return oidcAppRead(d, m)
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
	return oidcAppRead(d, m)
}

// oidcAppDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an oidcApp and its sub-resources
func oidcAppDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	err := client.Services.AppsV2.Destroy(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the oidcApp!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted oidcApp with %d", aid)
		d.SetId("")
	}

	return nil
}
