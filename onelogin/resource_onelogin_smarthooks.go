package onelogin

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	smarthooksschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook"
)

// SmartHooks attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func SmartHooks() *schema.Resource {
	smarthookSchema := smarthooksschema.Schema()

	return &schema.Resource{
		CreateContext: smartHookCreate,
		ReadContext:   smartHookRead,
		UpdateContext: smartHookUpdate,
		DeleteContext: smartHookDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        smarthookSchema,
	}
}

// smartHookCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create a SmartHook with its sub-resources
func smartHookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	// Build the new SmartHook model directly using v4 SDK types
	hook := models.SmartHook{}

	// Set basic fields from schema
	hookType := d.Get("type").(string)
	hook.Type = &hookType

	if v, ok := d.GetOk("disabled"); ok {
		disabled := v.(bool)
		hook.Disabled = &disabled
	}

	if v, ok := d.GetOk("timeout"); ok {
		timeout := int32(v.(int))
		hook.Timeout = &timeout
	}

	if v, ok := d.GetOk("runtime"); ok {
		runtime := v.(string)
		hook.Runtime = &runtime
	}

	if v, ok := d.GetOk("context_version"); ok {
		contextVersion := v.(string)
		hook.ContextVersion = &contextVersion
	}

	if v, ok := d.GetOk("retries"); ok {
		retries := int32(v.(int))
		hook.Retries = &retries
	}

	if v, ok := d.GetOk("function"); ok {
		function := v.(string)
		hook.Function = &function
	}

	// Handle packages (map[string]string)
	if v, ok := d.GetOk("packages"); ok {
		packagesMap := v.(map[string]interface{})
		packages := make(map[string]string)
		for k, v := range packagesMap {
			packages[k] = v.(string)
		}
		hook.Packages = packages
	}

	// Handle options
	if v, ok := d.GetOk("options"); ok {
		optsList := v.(*schema.Set).List()
		if len(optsList) > 0 {
			optsMap := optsList[0].(map[string]interface{})
			opts := &models.Options{}

			if v, ok := optsMap["risk_enabled"]; ok {
				riskEnabled := v.(bool)
				opts.RiskEnabled = &riskEnabled
			}

			if v, ok := optsMap["location_enabled"]; ok {
				locationEnabled := v.(bool)
				opts.LocationEnabled = &locationEnabled
			}

			if v, ok := optsMap["mfa_device_info_enabled"]; ok {
				mfaEnabled := v.(bool)
				opts.MFADeviceInfoEnabled = &mfaEnabled
			}

			hook.Options = opts
		}
	}

	// Handle conditions
	if v, ok := d.GetOk("conditions"); ok {
		condsList := v.([]interface{})
		conditions := make([]models.Condition, 0, len(condsList))

		for _, c := range condsList {
			condMap := c.(map[string]interface{})
			condition := models.Condition{
				Source:   condMap["source"].(string),
				Operator: condMap["operator"].(string),
				Value:    condMap["value"].(string),
			}
			conditions = append(conditions, condition)
		}

		hook.Conditions = conditions
	}

	// Create the hook
	result, err := client.CreateHook(hook)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem creating the smart hook!", map[string]interface{}{"error": err})
		return diag.FromErr(err)
	}

	// Extract ID from result
	hookMap, ok := result.(map[string]interface{})
	if !ok || hookMap["id"] == nil {
		return diag.Errorf("Failed to parse smarthook creation response or hook ID not found in response")
	}

	hookID := hookMap["id"].(string)
	tflog.Info(ctx, "[CREATED] Created smart hook", map[string]interface{}{"id": hookID})

	d.SetId(hookID)
	return smartHookRead(ctx, d, m)
}

// SmartHookRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read a SmartHook with its sub-resources
func smartHookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	result, err := client.GetHook(d.Id(), nil)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem reading the smarthook!", map[string]interface{}{"error": err})
		return diag.FromErr(err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	hookMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("Failed to parse smarthook response")
	}

	tflog.Info(ctx, "[READ] Reading hook", map[string]interface{}{"id": d.Id()})

	d.Set("type", hookMap["type"])
	d.Set("disabled", hookMap["disabled"])
	d.Set("timeout", hookMap["timeout"])

	// Handle env_vars - these might be structured differently in v4 API
	if envVars, ok := hookMap["env_vars"].([]interface{}); ok {
		d.Set("env_vars", flattenHookEnvVars(envVars))
	}

	d.Set("runtime", hookMap["runtime"])
	d.Set("context_version", hookMap["context_version"])
	d.Set("retries", hookMap["retries"])

	// Handle options
	if options, ok := hookMap["options"].(map[string]interface{}); ok {
		d.Set("options", []map[string]interface{}{options})
	}

	// Handle packages
	if packages, ok := hookMap["packages"].(map[string]interface{}); ok {
		d.Set("packages", packages)
	}

	// Handle function code
	d.Set("function", hookMap["function"])

	// Handle conditions
	if conditions, ok := hookMap["conditions"].([]interface{}); ok {
		conditionsList := make([]map[string]interface{}, 0, len(conditions))
		for _, c := range conditions {
			if condMap, ok := c.(map[string]interface{}); ok {
				condition := map[string]interface{}{
					"source":   condMap["source"],
					"operator": condMap["operator"],
					"value":    condMap["value"],
				}
				conditionsList = append(conditionsList, condition)
			}
		}
		d.Set("conditions", conditionsList)
	}

	d.Set("status", hookMap["status"])
	d.Set("created_at", hookMap["created_at"])
	d.Set("updated_at", hookMap["updated_at"])

	return nil
}

// Helper function to flatten environment variables to match expected schema format
func flattenHookEnvVars(envVars []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(envVars))
	for _, env := range envVars {
		if envMap, ok := env.(map[string]interface{}); ok {
			flattened := map[string]interface{}{
				"id":    envMap["id"],
				"name":  envMap["name"],
				"value": envMap["value"],
			}
			result = append(result, flattened)
		}
	}
	return result
}

// SmartHookUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update a SmartHook and its sub-resources
func smartHookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	// Build the new SmartHook model directly using v4 SDK types
	hook := models.SmartHook{}

	// Set ID
	hookID := d.Id()
	hook.ID = &hookID

	// Set basic fields from schema
	hookType := d.Get("type").(string)
	hook.Type = &hookType

	if v, ok := d.GetOk("disabled"); ok {
		disabled := v.(bool)
		hook.Disabled = &disabled
	}

	if v, ok := d.GetOk("timeout"); ok {
		timeout := int32(v.(int))
		hook.Timeout = &timeout
	}

	if v, ok := d.GetOk("runtime"); ok {
		runtime := v.(string)
		hook.Runtime = &runtime
	}

	if v, ok := d.GetOk("context_version"); ok {
		contextVersion := v.(string)
		hook.ContextVersion = &contextVersion
	}

	if v, ok := d.GetOk("retries"); ok {
		retries := int32(v.(int))
		hook.Retries = &retries
	}

	if v, ok := d.GetOk("function"); ok {
		function := v.(string)
		hook.Function = &function
	}

	// Handle packages (map[string]string)
	if v, ok := d.GetOk("packages"); ok {
		packagesMap := v.(map[string]interface{})
		packages := make(map[string]string)
		for k, v := range packagesMap {
			packages[k] = v.(string)
		}
		hook.Packages = packages
	}

	// Handle options
	if v, ok := d.GetOk("options"); ok {
		optsList := v.(*schema.Set).List()
		if len(optsList) > 0 {
			optsMap := optsList[0].(map[string]interface{})
			opts := &models.Options{}

			if v, ok := optsMap["risk_enabled"]; ok {
				riskEnabled := v.(bool)
				opts.RiskEnabled = &riskEnabled
			}

			if v, ok := optsMap["location_enabled"]; ok {
				locationEnabled := v.(bool)
				opts.LocationEnabled = &locationEnabled
			}

			if v, ok := optsMap["mfa_device_info_enabled"]; ok {
				mfaEnabled := v.(bool)
				opts.MFADeviceInfoEnabled = &mfaEnabled
			}

			hook.Options = opts
		}
	}

	// Handle conditions
	if v, ok := d.GetOk("conditions"); ok {
		condsList := v.([]interface{})
		conditions := make([]models.Condition, 0, len(condsList))

		for _, c := range condsList {
			condMap := c.(map[string]interface{})
			condition := models.Condition{
				Source:   condMap["source"].(string),
				Operator: condMap["operator"].(string),
				Value:    condMap["value"].(string),
			}
			conditions = append(conditions, condition)
		}

		hook.Conditions = conditions
	}

	// Update the hook
	_, err := client.UpdateSmartHook(d.Id(), hook)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem updating the smart hook!", map[string]interface{}{"error": err})
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "[UPDATED] Updated smart hook", map[string]interface{}{"id": d.Id()})
	return smartHookRead(ctx, d, m)
}

// smartHookDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete a SmartHook
func smartHookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	result, err := client.DeleteHook(d.Id())
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem deleting the smart hook!", map[string]interface{}{"error": err})
		return diag.FromErr(err)
	}

	// Check result for errors
	if resultMap, ok := result.(map[string]interface{}); ok {
		if status, exists := resultMap["status"]; exists && status.(string) != "ok" {
			tflog.Error(ctx, "[ERROR] There was a problem deleting the smart hook!", map[string]interface{}{"status": status})
			return diag.Errorf("Failed to delete hook: %v", status)
		}
	}

	tflog.Info(ctx, "[DELETED] Deleted smart hook", map[string]interface{}{"id": d.Id()})
	d.SetId("")

	return nil
}
