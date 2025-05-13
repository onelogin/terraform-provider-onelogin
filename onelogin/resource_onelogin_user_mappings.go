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

// UserMappings returns a resource with the CRUD methods and Terraform Schema defined
func UserMappings() *schema.Resource {
	return &schema.Resource{
		CreateContext: userMappingCreate,
		ReadContext:   userMappingRead,
		UpdateContext: userMappingUpdate,
		DeleteContext: userMappingDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        usermappingschema.Schema(),
	}
}

// userMappingCreate creates a new user mapping in OneLogin
func userMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if result != nil && result.ID != nil {
		mappingID := int(*result.ID)
		tflog.Info(ctx, "[CREATED] Created user mapping", map[string]interface{}{
			"id":   mappingID,
			"name": d.Get("name").(string),
		})

		d.SetId(fmt.Sprintf("%d", mappingID))
	} else {
		return diag.Errorf("failed to extract user mapping ID from response")
	}

	return userMappingRead(ctx, d, m)
}

// userMappingRead reads a user mapping by ID from OneLogin
func userMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	mid, _ := strconv.Atoi(d.Id())
	mid32 := int32(mid)

	tflog.Info(ctx, "[READ] Reading user mapping", map[string]interface{}{
		"id": mid,
	})

	result, err := client.GetUserMapping(mid32)
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

	// Set basic fields
	if result.Name != nil {
		d.Set("name", *result.Name)
	}
	if result.Match != nil {
		d.Set("match", *result.Match)
	}
	if result.Enabled != nil {
		d.Set("enabled", *result.Enabled)
	}
	if result.Position != nil {
		d.Set("position", *result.Position)
	}

	// Handle conditions
	if len(result.Conditions) > 0 {
		conditions := make([]map[string]interface{}, len(result.Conditions))
		for i, condition := range result.Conditions {
			condMap := map[string]interface{}{}
			if condition.Source != nil {
				condMap["source"] = *condition.Source
			}
			if condition.Operator != nil {
				condMap["operator"] = *condition.Operator
			}
			if condition.Value != nil {
				condMap["value"] = *condition.Value
			}
			conditions[i] = condMap
		}
		d.Set("conditions", conditions)
	}

	// Handle actions
	if len(result.Actions) > 0 {
		actions := make([]map[string]interface{}, len(result.Actions))
		for i, action := range result.Actions {
			actMap := map[string]interface{}{}
			if action.Action != nil {
				actMap["action"] = *action.Action
			}
			actMap["value"] = action.Value
			actions[i] = actMap
		}
		d.Set("actions", actions)
	}

	return nil
}

// userMappingUpdate updates a user mapping by ID in OneLogin
func userMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mid, _ := strconv.Atoi(d.Id())
	mid32 := int32(mid)

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

	_, err = client.UpdateUserMapping(mid32, userMapping)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "User Mapping", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated user mapping", map[string]interface{}{
		"id": mid,
	})

	return userMappingRead(ctx, d, m)
}

// userMappingDelete deletes a user mapping by ID from OneLogin
func userMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	mid, _ := strconv.Atoi(d.Id())
	mid32 := int32(mid)

	tflog.Info(ctx, "[DELETE] Deleting user mapping", map[string]interface{}{
		"id": mid,
	})

	err := client.DeleteUserMapping(mid32)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryDelete, "User Mapping", d.Id())
	}

	tflog.Info(ctx, "[DELETED] Deleted user mapping", map[string]interface{}{
		"id": mid,
	})

	d.SetId("")
	return nil
}
