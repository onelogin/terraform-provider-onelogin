package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	appschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app"
	appconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
	appparametersschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	appprovisioningschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
	appssoschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/sso"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// SAMLAppsSimplified attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func SAMLAppsSimplified() *schema.Resource {
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
	appSchema["certificate"] = &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	return &schema.Resource{
		CreateContext: samlAppCreateSimplified,
		ReadContext:   samlAppReadSimplified,
		UpdateContext: samlAppUpdateSimplified,
		DeleteContext: samlAppDeleteSimplified,
		Importer:      &schema.ResourceImporter{},
		Schema:        appSchema,
	}
}

// samlAppCreateSimplified takes a context, pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create a SAML app with its sub-resources
func samlAppCreateSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	samlApp, err := appschema.Inflate(map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryCreate, "SAML App", "")
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating SAML app", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	result, err := client.CreateApp(samlApp)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "SAML App", "")
	}

	// Extract app ID from the result
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse SAML app creation response")
	}

	id, ok := appMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract SAML app ID from response")
	}

	appID := int(id)
	tflog.Info(ctx, "[CREATED] Created SAML app", map[string]interface{}{
		"id": appID,
	})

	d.SetId(fmt.Sprintf("%d", appID))
	return samlAppReadSimplified(ctx, d, m)
}

// samlAppReadSimplified takes a context, pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read a SAML app with its sub-resources
func samlAppReadSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	aid, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "[READ] Reading SAML app", map[string]interface{}{
		"id": aid,
	})

	result, err := client.GetAppByID(aid, nil)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "SAML App", d.Id())
	}

	// Check if app exists
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] SAML app not found", map[string]interface{}{
			"id": aid,
		})
		d.SetId("")
		return nil
	}

	// Parse the app map from the result
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse SAML app response")
	}

	// Set basic fields
	basicFields := []string{
		"name", "visible", "description", "notes", "icon_url",
		"auth_method", "policy_id", "allow_assumed_signin", "tab_id",
		"brand_id", "connector_id", "created_at", "updated_at",
	}
	utils.SetResourceFields(d, appMap, basicFields)

	// Handle parameters if they exist
	if v, ok := appMap["parameters"]; ok {
		if params, ok := v.(map[string]interface{}); ok {
			d.Set("parameters", appparametersschema.FlattenV4(params))
		}
	}

	// Handle provisioning if it exists
	if v, ok := appMap["provisioning"]; ok {
		if provData, ok := v.(map[string]interface{}); ok {
			d.Set("provisioning", appprovisioningschema.FlattenMap(provData))
		}
	}

	// Handle configuration if it exists
	if v, ok := appMap["configuration"]; ok {
		if configData, ok := v.(map[string]interface{}); ok {
			d.Set("configuration", appconfigurationschema.Flatten(configData))
		}
	}

	// Handle sso if it exists
	if v, ok := appMap["sso"]; ok {
		if ssoData, ok := v.(map[string]interface{}); ok {
			d.Set("sso", appssoschema.Flatten(ssoData))
		}
	}

	return nil
}

// samlAppUpdateSimplified takes a context, pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update a SAML app and its sub-resources
func samlAppUpdateSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aid, _ := strconv.Atoi(d.Id())

	samlApp, err := appschema.Inflate(map[string]interface{}{
		"id":                   d.Id(),
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryUpdate, "SAML App", d.Id())
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[UPDATE] Updating SAML app", map[string]interface{}{
		"id": aid,
	})

	_, err = client.UpdateApp(aid, samlApp)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "SAML App", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated SAML app", map[string]interface{}{
		"id": aid,
	})

	return samlAppReadSimplified(ctx, d, m)
}

// samlAppDeleteSimplified takes a context, pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete a SAML app
func samlAppDeleteSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		aid, _ := strconv.Atoi(id)
		return client.DeleteApp(aid)
	}, "SAML App")
}
