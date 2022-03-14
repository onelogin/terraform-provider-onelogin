package onelogin

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	smarthooksschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook"
	smarthookconditionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/conditions"
	smarthookoptions "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/options"
)

// SmartHooks attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func SmartHooks() *schema.Resource {
	smarthookSchema := smarthooksschema.Schema()

	return &schema.Resource{
		Create:   smartHookCreate,
		Read:     smartHookRead,
		Update:   smartHookUpdate,
		Delete:   smartHookDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   smarthookSchema,
	}
}

// smartHookCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an samlApp with its sub-resources
func smartHookCreate(d *schema.ResourceData, m interface{}) error {
	smarthook := smarthooksschema.Inflate(map[string]interface{}{
		"type":            d.Get("type"),
		"disabled":        d.Get("disabled"),
		"timeout":         d.Get("timeout"),
		"env_vars":        d.Get("env_vars"),
		"runtime":         d.Get("runtime"),
		"context_version": d.Get("context_version"),
		"retries":         d.Get("retries"),
		"options":         d.Get("options"),
		"packages":        d.Get("packages"),
		"function":        d.Get("function"),
		"conditions":      d.Get("conditions"),
		"status":          d.Get("status"),
	})
	client := m.(*client.APIClient)
	smarthook.EncodeFunction()
	fullSmarthook, err := client.Services.SmartHooksV1.Create(&smarthook)
	if err != nil {
		log.Println("[ERROR] There was a problem creating the smart hooks!", err)
		return err
	}
	log.Printf("[CREATED] Created smart hook with %s", *(fullSmarthook.ID))

	d.SetId(fmt.Sprintf("%s", *(fullSmarthook.ID)))
	return smartHookRead(d, m)
}

// SmartHookRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an samlApp with its sub-resources
func smartHookRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	smarthook, err := client.Services.SmartHooksV1.GetOne(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the smarthook!")
		log.Println(err)
		return err
	}
	if smarthook == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[READ] Reading hook with %s", *(smarthook.ID))
	d.Set("type", smarthook.Type)
	d.Set("disabled", smarthook.Disabled)
	d.Set("timeout", smarthook.Timeout)
	d.Set("env_vars", smarthooksschema.FlattenEnvVars(smarthook.EnvVars))
	d.Set("runtime", smarthook.Runtime)
	d.Set("context_version", smarthook.ContextVersion)
	d.Set("retries", smarthook.Retries)
	d.Set("options", smarthookoptions.Flatten(*smarthook.Options))
	d.Set("packages", smarthook.Packages)
	d.Set("function", smarthook.Function)
	d.Set("conditions", smarthookconditionsschema.Flatten(smarthook.Conditions))
	d.Set("status", smarthook.Status)
	d.Set("created_at", smarthook.CreatedAt.String())
	d.Set("updated_at", smarthook.UpdatedAt.String())

	return nil
}

// SmartHookUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an samlApp and its sub-resources
func smartHookUpdate(d *schema.ResourceData, m interface{}) error {
	smartHook := smarthooksschema.Inflate(map[string]interface{}{
		"id":              d.Id(),
		"type":            d.Get("type"),
		"disabled":        d.Get("disabled"),
		"timeout":         d.Get("timeout"),
		"env_vars":        d.Get("env_vars"),
		"runtime":         d.Get("runtime"),
		"context_version": d.Get("context_version"),
		"retries":         d.Get("retries"),
		"options":         d.Get("options"),
		"packages":        d.Get("packages"),
		"function":        d.Get("function"),
		"conditions":      d.Get("conditions"),
		"status":          d.Get("status"),
	})

	client := m.(*client.APIClient)
	smartHook.EncodeFunction()
	fullSmartHook, err := client.Services.SmartHooksV1.Update(&smartHook)
	if err != nil {
		log.Println("[ERROR] There was a problem Updating the smart hooks!", err)
		return err
	}
	if fullSmartHook.ID == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[UPDATED] Updated smart hook with %s", *(fullSmartHook.ID))
	d.SetId(fmt.Sprintf("%s", *(fullSmartHook.ID)))
	return smartHookRead(d, m)
}

// smartHookDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete a smart hooks
func smartHookDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	err := client.Services.SmartHooksV1.Destroy(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the smart hooks!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted smart hooks with %s", d.Id())
		d.SetId("")
	}

	return nil
}
