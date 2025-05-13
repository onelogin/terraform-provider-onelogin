//go:build exclude
// +build exclude

package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	usermappingschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// UserMappingsSimplified returns a resource with the CRUD methods and Terraform Schema defined
func UserMappingsSimplified() *schema.Resource {
	return &schema.Resource{
		CreateContext: userMappingCreateSimplified,
		ReadContext:   userMappingReadSimplified,
		UpdateContext: userMappingUpdateSimplified,
		DeleteContext: userMappingDeleteSimplified,
		Importer:      &schema.ResourceImporter{},
		Schema:        usermappingschema.Schema(),
	}
}

// userMappingCreateSimplified creates a new user mapping in OneLogin
func userMappingCreateSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userMapping, err := usermappingschema.Inflate(map[string]interface{}{
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"enabled":    d.Get("enabled"),
		"position":   d.Get("position"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryCreate, "User Mapping", "")
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating user mapping", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	result, err := client.CreateUserMapping(userMapping)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "User Mapping", "")
	}

	// Extract user mapping ID from the result
	mappingMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse user mapping creation response")
	}

	id, ok := mappingMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract user mapping ID from response")
	}

	mappingID := int(id)
	tflog.Info(ctx, "[CREATED] Created user mapping", map[string]interface{}{
		"id":   mappingID,
		"name": d.Get("name").(string),
	})

	d.SetId(fmt.Sprintf("%d", mappingID))
	return userMappingReadSimplified(ctx, d, m)
}

// userMappingReadSimplified reads a user mapping by ID from OneLogin
func userMappingReadSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	mid, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "[READ] Reading user mapping", map[string]interface{}{
		"id": mid,
	})

	result, err := client.GetUserMappingByID(mid)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "User Mapping", d.Id())
	}

	// Check if mapping exists
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] User mapping not found", map[string]interface{}{
			"id": mid,
		})
		d.SetId("")
		return nil
	}

	// Parse the mapping from the result
	mappingMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse user mapping response")
	}

	// Set basic fields
	basicFields := []string{
		"name", "match", "enabled", "position", "created_at", "updated_at",
	}
	utils.SetResourceFields(d, mappingMap, basicFields)

	// Handle conditions
	if conditions, ok := mappingMap["conditions"].([]interface{}); ok {
		if err := d.Set("conditions", usermappingschema.FlattenConditions(conditions)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting conditions: %s", err))
		}
	}

	// Handle actions
	if actions, ok := mappingMap["actions"].([]interface{}); ok {
		if err := d.Set("actions", usermappingschema.FlattenActions(actions)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting actions: %s", err))
		}
	}

	return nil
}

// userMappingUpdateSimplified updates a user mapping by ID in OneLogin
func userMappingUpdateSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mid, _ := strconv.Atoi(d.Id())

	userMapping, err := usermappingschema.Inflate(map[string]interface{}{
		"id":         d.Id(),
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"enabled":    d.Get("enabled"),
		"position":   d.Get("position"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryUpdate, "User Mapping", d.Id())
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[UPDATE] Updating user mapping", map[string]interface{}{
		"id": mid,
	})

	_, err = client.UpdateUserMapping(mid, userMapping)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "User Mapping", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated user mapping", map[string]interface{}{
		"id": mid,
	})

	return userMappingReadSimplified(ctx, d, m)
}

// userMappingDeleteSimplified deletes a user mapping by ID from OneLogin
func userMappingDeleteSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		mid, _ := strconv.Atoi(id)
		return client.DeleteUserMapping(mid)
	}, "User Mapping")
}
