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
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	usermappingschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping"
	usermappingactionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/actions"
	usermappingconditionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/conditions"
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
	userMapping := usermappingschema.Inflate(map[string]interface{}{
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"enabled":    d.Get("enabled"),
		"position":   d.Get("position"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})

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
	return userMappingRead(ctx, d, m)
}

// userMappingRead reads a user mapping by ID from OneLogin
func userMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		if err := d.Set("conditions", usermappingconditionsschema.Flatten(convertToUserMappingConditions(conditions))); err != nil {
			return diag.FromErr(fmt.Errorf("error setting conditions: %s", err))
		}
	}

	// Handle actions
	if actions, ok := mappingMap["actions"].([]interface{}); ok {
		if err := d.Set("actions", usermappingactionsschema.Flatten(convertToUserMappingActions(actions))); err != nil {
			return diag.FromErr(fmt.Errorf("error setting actions: %s", err))
		}
	}

	return nil
}

// userMappingUpdate updates a user mapping by ID in OneLogin
func userMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mid, _ := strconv.Atoi(d.Id())

	userMapping := usermappingschema.Inflate(map[string]interface{}{
		"id":         d.Id(),
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"enabled":    d.Get("enabled"),
		"position":   d.Get("position"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})

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

	return userMappingRead(ctx, d, m)
}

// userMappingDelete deletes a user mapping by ID from OneLogin
func userMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		mid, _ := strconv.Atoi(id)
		return client.DeleteUserMapping(mid)
	}, "User Mapping")
}

// convertToUserMappingConditions converts an array of interface{} to an array of models.UserMappingConditions
func convertToUserMappingConditions(conditions []interface{}) []models.UserMappingConditions {
	result := make([]models.UserMappingConditions, len(conditions))
	for i, condition := range conditions {
		if condMap, ok := condition.(map[string]interface{}); ok {
			var src, op, val string
			if s, ok := condMap["source"].(string); ok {
				src = s
			}
			if o, ok := condMap["operator"].(string); ok {
				op = o
			}
			if v, ok := condMap["value"].(string); ok {
				val = v
			}
			result[i] = models.UserMappingConditions{
				Source:   &src,
				Operator: &op,
				Value:    &val,
			}
		}
	}
	return result
}

// convertToUserMappingActions converts an array of interface{} to an array of models.UserMappingActions
func convertToUserMappingActions(actions []interface{}) []models.UserMappingActions {
	result := make([]models.UserMappingActions, len(actions))
	for i, action := range actions {
		if actMap, ok := action.(map[string]interface{}); ok {
			var act string
			var vals []string
			if a, ok := actMap["action"].(string); ok {
				act = a
			}
			if v, ok := actMap["value"].([]interface{}); ok {
				vals = make([]string, len(v))
				for j, val := range v {
					if s, ok := val.(string); ok {
						vals[j] = s
					}
				}
			}
			result[i] = models.UserMappingActions{
				Action: &act,
				Value:  vals,
			}
		}
	}
	return result
}
