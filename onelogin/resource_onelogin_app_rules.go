package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/rules"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/actions"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/conditions"
)

// AppRules returns a resource with the CRUD methods and Terraform Schema defined
func AppRules() *schema.Resource {
	return &schema.Resource{
		Create:   appRuleCreate,
		Read:     appRuleRead,
		Update:   appRuleUpdate,
		Delete:   appRuleDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   apprulesschema.Schema(),
	}
}

// appRuleCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an App with its sub-resources
func appRuleCreate(d *schema.ResourceData, m interface{}) error {
	appRule := apprulesschema.Inflate(map[string]interface{}{
		"app_id":     d.Get("app_id"),
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"position":   d.Get("position"),
		"enabled":    d.Get("enabled"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	client := m.(*client.APIClient)
	err := client.Services.AppRulesV2.Create(&appRule)
	if err != nil {
		log.Println("[ERROR] There was a problem creating the app rule!", err)
		return err
	}
	log.Printf("[CREATED] Created app rule with %d", *(appRule.ID))

	d.SetId(fmt.Sprintf("%d", *(appRule.ID)))
	return appRuleRead(d, m)
}

// appRuleRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an App with its sub-resources
func appRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	id, _ := strconv.Atoi(d.Id())
	appId, _ := strconv.Atoi(d.Get("app_id").(string))
	app, err := client.Services.AppRulesV2.GetOne(int32(appId), int32(id))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the app rule!")
		log.Println(err)
		return err
	}
	if app == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading app rule with %d", *(app.ID))

	d.Set("name", app.Name)
	d.Set("match", app.Match)
	d.Set("position", app.Position)
	d.Set("enabled", app.Enabled)

	d.Set("conditions", appruleconditionsschema.Flatten(app.Conditions))
	d.Set("actions", appruleactionsschema.Flatten(app.Actions))

	return nil
}

// appRuleUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an App and its sub-resources
func appRuleUpdate(d *schema.ResourceData, m interface{}) error {
	appRule := apprulesschema.Inflate(map[string]interface{}{
		"id":         d.Id(),
		"app_id":     d.Get("app_id"),
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"position":   d.Get("position"),
		"enabled":    d.Get("enabled"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	client := m.(*client.APIClient)

	err := client.Services.AppRulesV2.Update(&appRule)
	if err != nil {
		log.Println("[ERROR] There was a problem updating the app rule!", err)
		return err
	}
	if appRule.ID == nil { // app must be deleted in api so remove from tf state
		d.SetId("")
		return nil
	}
	log.Printf("[UPDATED] Updated app rule with %d", *(appRule.ID))
	d.SetId(fmt.Sprintf("%d", *(appRule.ID)))
	return appRuleRead(d, m)
}

// appRuleDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an App and its sub-resources
func appRuleDelete(d *schema.ResourceData, m interface{}) error {
	id, _ := strconv.Atoi(d.Id())
	appId, _ := strconv.Atoi(d.Get("app_id").(string))
	client := m.(*client.APIClient)

	err := client.Services.AppRulesV2.Destroy(int32(appId), int32(id))
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the app rule!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted app rule with %d", id)
		d.SetId("")
	}

	return nil
}
