package onelogin

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	smarthookenvironmentvariablesschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/environment_variable"
)

// SmarthookEnvironmentVariables returns a resource with the CRUD methods and Terraform Schema defined
func SmarthookEnvironmentVariables() *schema.Resource {
	return &schema.Resource{
		CreateContext: environmentVariablesCreate,
		ReadContext:   environmentVariablesRead,
		UpdateContext: environmentVariablesUpdate,
		DeleteContext: environmentVariablesDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        smarthookenvironmentvariablesschema.Schema(),
	}
}

func environmentVariablesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	envVar := smarthookenvironmentvariablesschema.Inflate(map[string]interface{}{
		"name":  d.Get("name"),
		"value": d.Get("value"),
	})

	result, err := client.CreateEnvironmentVariable(envVar)
	if err != nil {
		return diag.Errorf("error creating environment variable: %v", err)
	}

	// Extract the variable ID from the response
	envVarMap, ok := result.(map[string]interface{})
	if !ok || envVarMap["id"] == nil {
		return diag.Errorf("failed to parse environment variable creation response or variable ID not found in response")
	}

	envVarID := envVarMap["id"].(string)
	d.SetId(envVarID)
	log.Printf("[CREATED] Created environment variable with id %s", envVarID)

	return environmentVariablesRead(ctx, d, m)
}

func environmentVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	result, err := client.GetEnvironmentVariable(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the environment variable: %v", err)
		return diag.FromErr(err)
	}

	// Check if the resource was not found
	if result == nil {
		d.SetId("")
		return nil
	}

	// Parse the response
	envVarMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse environment variable response")
	}

	log.Printf("[READ] Reading environment variable with id %s", d.Id())

	if envVarMap["name"] != nil {
		d.Set("name", envVarMap["name"])
	}

	if envVarMap["value"] != nil {
		d.Set("value", envVarMap["value"])
	}

	// Handle created_at and updated_at if they exist in the response
	if envVarMap["created_at"] != nil {
		createdAt, ok := envVarMap["created_at"].(string)
		if ok {
			d.Set("created_at", createdAt)
		}
	}

	if envVarMap["updated_at"] != nil {
		updatedAt, ok := envVarMap["updated_at"].(string)
		if ok {
			d.Set("updated_at", updatedAt)
		}
	}

	return nil
}

func environmentVariablesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	envVar := smarthookenvironmentvariablesschema.Inflate(map[string]interface{}{
		"id":    d.Id(),
		"value": d.Get("value"),
	})

	_, err := client.UpdateEnvironmentVariable(d.Id(), envVar)
	if err != nil {
		return diag.Errorf("error updating environment variable: %v", err)
	}

	log.Printf("[UPDATED] Updated environment variable with id %s", d.Id())

	// Wait a second for the update to propagate
	time.Sleep(1 * time.Second)

	return environmentVariablesRead(ctx, d, m)
}

func environmentVariablesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	_, err := client.DeleteEnvironmentVariable(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the environment variable: %v", err)
		return diag.FromErr(err)
	}

	log.Printf("[DELETED] Deleted environment variable with id %s", d.Id())
	d.SetId("")

	return nil
}
