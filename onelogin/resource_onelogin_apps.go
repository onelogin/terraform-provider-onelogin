package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	appschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app"
	appparametersschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	appprovisioningschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Apps returns a resource with enhanced CRUD methods and the appropriate schemas
// This implementation uses utility functions for better error handling and logging
func Apps() *schema.Resource {
	return &schema.Resource{
		CreateContext: appCreate,
		ReadContext:   appRead,
		UpdateContext: appUpdate,
		DeleteContext: appDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        appschema.Schema(),
	}
}

// appCreate creates an app
func appCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	basicApp, _ := appschema.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
	})

	client := m.(*onelogin.OneloginSDK)
	result, err := client.CreateApp(basicApp)
	if err != nil {
		tflog.Error(ctx, "[ERROR] Error creating app", map[string]interface{}{"error": err})
		return diag.FromErr(err)
	}

	// Extract app ID from the result
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse app creation response")
	}

	id, ok := appMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract app ID from response")
	}

	appID := int(id)
	tflog.Info(ctx, "[CREATED] Created app", map[string]interface{}{"id": appID})

	d.SetId(fmt.Sprintf("%d", appID))
	return appRead(ctx, d, m)
}

// appRead reads an app
func appRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	aid, _ := strconv.Atoi(d.Id())

	result, err := client.GetAppByID(aid, nil)
	if err != nil {
		// Check if this is a 404 (resource not found)
		if utils.IsNotFoundError(err) {
			tflog.Info(ctx, "[NOT FOUND] App not found", map[string]interface{}{"id": aid})
			d.SetId("")
			return nil
		}
		// For other errors, log and return the error
		tflog.Error(ctx, "[ERROR] Error reading app", map[string]interface{}{"id": aid, "error": err})
		return diag.FromErr(err)
	}

	// Additional nil check for safety
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] App not found (nil result)", map[string]interface{}{"id": aid})
		d.SetId("")
		return nil
	}

	// Parse the app map from the result
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse app response")
	}

	tflog.Info(ctx, "[READ] Reading app", map[string]interface{}{"id": aid})

	// Use utility function to set basic fields
	basicFields := []string{
		"name", "visible", "description", "notes", "icon_url",
		"auth_method", "policy_id", "allow_assumed_signin", "tab_id",
		"brand_id", "connector_id", "created_at", "updated_at",
	}
	utils.SetResourceFields(d, appMap, basicFields)

	// Handle parameters if they exist
	if v, ok := appMap["parameters"]; ok {
		if params, ok := v.(map[string]interface{}); ok {
			paramMap := make(map[string]models.Parameter)
			for key, val := range params {
				if paramData, ok := val.(map[string]interface{}); ok {
					// Convert each param to a Parameter model
					param := models.Parameter{}

					// Set string fields
					if label, ok := paramData["label"].(string); ok {
						param.Label = label
					}

					// Set optional fields
					param.UserAttributeMappings = paramData["user_attribute_mappings"]
					param.UserAttributeMacros = paramData["user_attribute_macros"]
					param.AttributesTransformations = paramData["attributes_transformations"]
					param.Values = paramData["values"]
					param.DefaultValues = paramData["default_values"]

					// Set boolean fields
					if skipIfBlank, ok := paramData["skip_if_blank"].(bool); ok {
						param.SkipIfBlank = skipIfBlank
					}
					if provisioned, ok := paramData["provisioned_entitlements"].(bool); ok {
						param.ProvisionedEntitlements = provisioned
					}
					if includeInAssertion, ok := paramData["include_in_saml_assertion"].(bool); ok {
						param.IncludeInSamlAssertion = includeInAssertion
					}

					// Set ID if it exists
					if id, ok := paramData["id"].(float64); ok {
						param.ID = int(id)
					}

					paramMap[key] = param
				}
			}
			d.Set("parameters", appparametersschema.Flatten(paramMap))
		}
	}

	// Handle provisioning if it exists
	if v, ok := appMap["provisioning"]; ok {
		if provData, ok := v.(map[string]interface{}); ok {
			if enabled, ok := provData["enabled"].(bool); ok {
				prov := models.Provisioning{
					Enabled: enabled,
				}
				d.Set("provisioning", appprovisioningschema.Flatten(prov))
			}
		}
	}

	return nil
}

// appUpdate updates an app
func appUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aid, _ := strconv.Atoi(d.Id())

	basicApp, _ := appschema.Inflate(map[string]interface{}{
		"id":                   d.Id(),
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"brand_id":             d.Get("brand_id"),
	})

	client := m.(*onelogin.OneloginSDK)
	_, err := client.UpdateApp(aid, basicApp)
	if err != nil {
		tflog.Error(ctx, "[ERROR] Error updating app", map[string]interface{}{"id": aid, "error": err})
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "[UPDATED] Updated app", map[string]interface{}{"id": aid})
	return appRead(ctx, d, m)
}

// appDelete deletes an app using the standard delete pattern
func appDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		aid, _ := strconv.Atoi(id)
		return client.DeleteApp(aid)
	}, "app")
}
