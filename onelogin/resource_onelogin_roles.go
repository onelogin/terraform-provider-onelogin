package onelogin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	roleschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/role"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Roles returns a roles resource with CRUD methods and the appropriate schemas
func Roles() *schema.Resource {
	return &schema.Resource{
		CreateContext: roleCreate,
		ReadContext:   roleRead,
		UpdateContext: roleUpdate,
		DeleteContext: roleDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        roleschema.Schema(),
	}
}

// roleCreate creates a new role in OneLogin
func roleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	role := roleschema.Inflate(map[string]interface{}{
		"name":   d.Get("name"),
		"apps":   d.Get("apps"),
		"admins": d.Get("admins"),
		"users":  d.Get("users"), // Include users in initial creation
	})

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating role", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	result, err := client.CreateRoleWithContext(ctx, role)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Role", "")
	}

	// Extract role ID from the result
	roleMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse role creation response")
	}

	id, ok := roleMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract role ID from response")
	}

	roleID := int(id)
	tflog.Info(ctx, "[CREATED] Created role", map[string]interface{}{
		"id":   roleID,
		"name": d.Get("name").(string),
	})

	d.SetId(fmt.Sprintf("%d", roleID))
	return roleRead(ctx, d, m)
}

// roleRead reads a role by ID from OneLogin
func roleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	rid, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "[READ] Reading role", map[string]interface{}{
		"id": rid,
	})

	result, err := client.GetRoleByIDWithContext(ctx, rid, &roleschema.RoleQuery{})
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Role", d.Id())
	}

	// Check if role exists
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] Role not found", map[string]interface{}{
			"id": rid,
		})
		d.SetId("")
		return nil
	}

	// Parse the role from the result
	roleMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse role response")
	}

	// Set basic fields
	d.Set("name", roleMap["name"])

	// Handle apps
	if v, ok := roleMap["apps"].([]interface{}); ok {
		var appIDs []int
		for _, app := range v {
			if id, ok := app.(float64); ok {
				appIDs = append(appIDs, int(id))
			}
		}
		d.Set("apps", appIDs)
	}

	// Handle users
	if v, ok := roleMap["users"].([]interface{}); ok {
		var userIDs []int
		for _, user := range v {
			if id, ok := user.(float64); ok {
				userIDs = append(userIDs, int(id))
			}
		}
		d.Set("users", userIDs)
	}

	// Handle admins
	if v, ok := roleMap["admins"].([]interface{}); ok {
		var adminIDs []int
		for _, admin := range v {
			if id, ok := admin.(float64); ok {
				adminIDs = append(adminIDs, int(id))
			}
		}
		d.Set("admins", adminIDs)
	}

	return nil
}

// roleUpdate updates a role by ID in OneLogin
func roleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rid, _ := strconv.Atoi(d.Id())
	client := m.(*onelogin.OneloginSDK)
	
	tflog.Info(ctx, "[UPDATE] Updating role", map[string]interface{}{
		"id": rid,
	})
	
	// Handle user changes separately using the specialized user management APIs
	if d.HasChange("users") {
		old, new := d.GetChange("users")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		
		// Users to add (in new set but not in old set)
		usersToAdd := newSet.Difference(oldSet)
		if usersToAdd.Len() > 0 {
			userIDs := make([]int, 0, usersToAdd.Len())
			for _, user := range usersToAdd.List() {
				userIDs = append(userIDs, user.(int))
			}
			
			tflog.Info(ctx, "[UPDATE] Adding users to role", map[string]interface{}{
				"role_id": rid,
				"users":   userIDs,
			})
			
			_, err := client.AddRoleUsers(rid, userIDs)
			if err != nil {
				tflog.Error(ctx, "[ERROR] Failed to add users to role", map[string]interface{}{
					"role_id": rid,
					"error":   err.Error(),
				})
				return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Role Users", d.Id())
			}
		}
		
		// Users to remove (in old set but not in new set)
		usersToRemove := oldSet.Difference(newSet)
		if usersToRemove.Len() > 0 {
			userIDs := make([]int, 0, usersToRemove.Len())
			for _, user := range usersToRemove.List() {
				userIDs = append(userIDs, user.(int))
			}
			
			tflog.Info(ctx, "[UPDATE] Removing users from role", map[string]interface{}{
				"role_id": rid,
				"users":   userIDs,
			})
			
			_, err := client.DeleteRoleUsers(rid, userIDs)
			if err != nil {
				tflog.Error(ctx, "[ERROR] Failed to remove users from role", map[string]interface{}{
					"role_id": rid,
					"error":   err.Error(),
				})
				return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Role Users", d.Id())
			}
		}
	}
	
	// Create a role object with all current values
	role := roleschema.Inflate(map[string]interface{}{
		"id":     d.Id(),
		"name":   d.Get("name"),
		"apps":   d.Get("apps"),
		"admins": d.Get("admins"),
		"users":  d.Get("users"), // Include users in the update
	})
	
	// Log the role object for debugging
	roleJSON, _ := json.Marshal(role)
	tflog.Info(ctx, "[DEBUG] Role object being sent", map[string]interface{}{
		"role_json": string(roleJSON),
	})
	
	// Update the role
	updateResponse, err := client.UpdateRoleWithContext(ctx, rid, role)
	if err != nil {
		// Log more details about the error
		tflog.Error(ctx, "[ERROR] Failed to update role", map[string]interface{}{
			"id":    rid,
			"error": err.Error(),
		})

		// Try to get more error details
		respJSON, _ := json.Marshal(updateResponse)
		tflog.Error(ctx, "[ERROR] Update response", map[string]interface{}{
			"response": string(respJSON),
		})

		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Role", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated role", map[string]interface{}{
		"id": rid,
	})

	return roleRead(ctx, d, m)
}

// roleDelete deletes a role by ID from OneLogin
func roleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		rid, _ := strconv.Atoi(id)
		return client.DeleteRoleWithContext(ctx, rid)
	}, "Role")
}