// Package onelogin provides resources for interacting with the OneLogin API
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
	client := m.(*onelogin.OneloginSDK)

	// Create a role object with name
	name := d.Get("name").(string)
	roleData := map[string]interface{}{
		"name": name,
	}

	// Add optional fields if present
	if users, ok := d.GetOk("users"); ok && users.(*schema.Set).Len() > 0 {
		userIDs := make([]int, 0, users.(*schema.Set).Len())
		for _, user := range users.(*schema.Set).List() {
			userIDs = append(userIDs, user.(int))
		}
		roleData["users"] = userIDs
	}

	if apps, ok := d.GetOk("apps"); ok && apps.(*schema.Set).Len() > 0 {
		appIDs := make([]int, 0, apps.(*schema.Set).Len())
		for _, app := range apps.(*schema.Set).List() {
			appIDs = append(appIDs, app.(int))
		}
		roleData["apps"] = appIDs
	}

	if admins, ok := d.GetOk("admins"); ok && admins.(*schema.Set).Len() > 0 {
		adminIDs := make([]int, 0, admins.(*schema.Set).Len())
		for _, admin := range admins.(*schema.Set).List() {
			adminIDs = append(adminIDs, admin.(int))
		}
		roleData["admins"] = adminIDs
	}

	tflog.Info(ctx, "[CREATE] Creating role with complete properties", map[string]interface{}{
		"name": name,
	})

	// Create the role object
	role := roleschema.Inflate(roleData)

	// Log the role object for debugging
	roleJSON, _ := json.MarshalIndent(role, "", "  ")
	tflog.Info(ctx, "[DEBUG] Role object being sent", map[string]interface{}{
		"role_json": string(roleJSON),
	})

	// Create the role with complete properties
	result, err := client.CreateRoleWithContext(ctx, role)
	if err != nil {
		tflog.Error(ctx, "[ERROR] Failed to create role", map[string]interface{}{
			"name":  name,
			"error": err.Error(),
		})

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
		"name": name,
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

	// Use the GetRoles endpoint which returns complete role objects
	// We'll need to handle pagination and filter to find our specific role
	foundRole := false
	var roleObj map[string]interface{}

	// Initialize query with pagination parameters
	query := &roleschema.RoleQuery{
		Limit: "100", // Get reasonable batch size
	}

	// Enable debugging
	utils.AddRequestResponseLogging(ctx, client)

	for !foundRole {
		// Get a batch of roles
		result, err := client.GetRolesWithContext(ctx, query)
		if err != nil {
			return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Role", d.Id())
		}

		// Parse the roles from the result
		roles, ok := result.([]interface{})
		if !ok {
			// Try to handle if the API returns an object with array under a key
			resultMap, mapOk := result.(map[string]interface{})
			if mapOk {
				if data, dataOk := resultMap["data"].([]interface{}); dataOk {
					roles = data
				} else {
					return diag.Errorf("failed to parse roles response: unexpected structure")
				}
			} else {
				return diag.Errorf("failed to parse roles response: not an array or map")
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("[READ] Found %d roles to search through", len(roles)))

		// Check if our role is in this batch
		for _, r := range roles {
			role, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			// Check if this is the role we're looking for
			if roleID, ok := role["id"].(float64); ok && int(roleID) == rid {
				roleObj = role
				foundRole = true
				break
			}
		}

		// If we found the role, break out of the pagination loop
		if foundRole {
			break
		}

		// Handle pagination
		// Check if there's more data to fetch
		if len(roles) < 100 {
			// No more roles to check, the role doesn't exist
			tflog.Info(ctx, "[NOT FOUND] Role not found in any page", map[string]interface{}{
				"id": rid,
			})
			d.SetId("")
			return nil
		}

		// Use the last role's ID as the cursor for the next page
		if lastRole, ok := roles[len(roles)-1].(map[string]interface{}); ok {
			if lastID, ok := lastRole["id"].(float64); ok {
				query.Cursor = fmt.Sprintf("%d", int(lastID))
				// Clear limit and page when using cursor - API requires cursor XOR pagination
				query.Limit = ""
				query.Page = ""
			}
		}
	}

	// If we get here and haven't found the role, it doesn't exist
	if !foundRole {
		tflog.Info(ctx, "[NOT FOUND] Role not found", map[string]interface{}{
			"id": rid,
		})
		d.SetId("")
		return nil
	}

	// Role found - set its properties in the state
	// No need to use Inflate here since we're directly reading from the API response

	// Set basic fields
	d.Set("name", roleObj["name"])

	// Handle apps
	if v, ok := roleObj["apps"].([]interface{}); ok {
		var appIDs []int
		for _, app := range v {
			if id, ok := app.(float64); ok {
				appIDs = append(appIDs, int(id))
			}
		}
		d.Set("apps", appIDs)
	} else {
		// Always ensure we have an empty array rather than nil
		d.Set("apps", []int{})
	}

	// Handle users
	if v, ok := roleObj["users"].([]interface{}); ok {
		var userIDs []int
		for _, user := range v {
			if id, ok := user.(float64); ok {
				userIDs = append(userIDs, int(id))
			}
		}
		d.Set("users", userIDs)
	} else {
		// Always ensure we have an empty array rather than nil
		d.Set("users", []int{})
	}

	// Handle admins
	if v, ok := roleObj["admins"].([]interface{}); ok {
		var adminIDs []int
		for _, admin := range v {
			if id, ok := admin.(float64); ok {
				adminIDs = append(adminIDs, int(id))
			}
		}
		d.Set("admins", adminIDs)
	} else {
		// Always ensure we have an empty array rather than nil
		d.Set("admins", []int{})
	}

	return nil
}

// roleUpdate updates a role by ID in OneLogin
func roleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rid, _ := strconv.Atoi(d.Id())
	client := m.(*onelogin.OneloginSDK)

	// Add debug logging to client
	utils.AddRequestResponseLogging(ctx, client)

	tflog.Info(ctx, "[UPDATE] Updating role", map[string]interface{}{
		"id": rid,
	})

	// Create a role object with the required fields
	// Don't include ID in the payload as we're already specifying it in the URL
	roleData := map[string]interface{}{
		"name": d.Get("name"),
	}

	// Handle users array - only add to roleData if present
	if users, ok := d.GetOk("users"); ok {
		userIDs := make([]int, 0, users.(*schema.Set).Len())
		for _, user := range users.(*schema.Set).List() {
			userIDs = append(userIDs, user.(int))
		}
		roleData["users"] = userIDs
		tflog.Info(ctx, "[UPDATE] Setting users for role", map[string]interface{}{
			"role_id": rid,
			"users":   userIDs,
		})
	}

	// Handle apps array - only add to roleData if present
	if apps, ok := d.GetOk("apps"); ok {
		appIDs := make([]int, 0, apps.(*schema.Set).Len())
		for _, app := range apps.(*schema.Set).List() {
			appIDs = append(appIDs, app.(int))
		}
		roleData["apps"] = appIDs
		tflog.Info(ctx, "[UPDATE] Setting apps for role", map[string]interface{}{
			"role_id": rid,
			"apps":    appIDs,
		})
	}

	// Handle admins array - only add to roleData if present
	if admins, ok := d.GetOk("admins"); ok {
		adminIDs := make([]int, 0, admins.(*schema.Set).Len())
		for _, admin := range admins.(*schema.Set).List() {
			adminIDs = append(adminIDs, admin.(int))
		}
		roleData["admins"] = adminIDs
		tflog.Info(ctx, "[UPDATE] Setting admins for role", map[string]interface{}{
			"role_id": rid,
			"admins":  adminIDs,
		})
	}

	// Create a role object from the data
	role := roleschema.Inflate(roleData)

	// Print the exact JSON payload that will be sent to the API
	// Debug logging of the role model
	tflog.Debug(ctx, "[DEBUG] Role model details", map[string]interface{}{
		"role_id":          rid,            // ID is in the URL path already, not in the payload
		"has_id_field":     role.ID != nil, // Should be nil for update operations
		"has_users_field":  role.Users != nil,
		"has_apps_field":   role.Apps != nil,
		"has_admins_field": role.Admins != nil,
		"users_len":        fmt.Sprintf("%v", (role.Users != nil && len(role.Users) > 0)),
		"apps_len":         fmt.Sprintf("%v", (role.Apps != nil && len(role.Apps) > 0)),
		"admins_len":       fmt.Sprintf("%v", (role.Admins != nil && len(role.Admins) > 0)),
	})

	// Print the exact JSON payload that will be sent to the API
	requestJSON, _ := json.MarshalIndent(role, "", "  ")
	tflog.Debug(ctx, "[DEBUG] Exact request payload", map[string]interface{}{
		"payload": string(requestJSON),
	})

	// Simple logging
	tflog.Info(ctx, "[UPDATE] Sending role update to API", map[string]interface{}{
		"role_id":      rid,
		"endpoint":     fmt.Sprintf("api/2/roles/%d", rid),
		"method":       "PUT",
		"has_id_field": role.ID != nil, // Should be false for updates
	})

	// Update the role with a single call
	resp, err := client.UpdateRoleWithContext(ctx, rid, role)
	if err != nil {
		// Print error details directly to console for visibility
		errorJSON, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Printf("\n\nAPI ERROR: %s\nRESPONSE: %s\n\n", err.Error(), string(errorJSON))

		// Try to extract more error details
		respJSON, _ := json.MarshalIndent(resp, "", "  ")
		tflog.Error(ctx, "[ERROR] Failed to update role", map[string]interface{}{
			"id":       rid,
			"error":    err.Error(),
			"response": string(respJSON),
		})

		// Print additional debug info with detailed type information
		roleMarshalled, _ := json.Marshal(role)
		tflog.Error(ctx, "[DEBUG] API Error Details", map[string]interface{}{
			"role_type":       fmt.Sprintf("%T", role),
			"users":           fmt.Sprintf("%v (type: %T)", role.Users, role.Users),
			"apps":            fmt.Sprintf("%v (type: %T)", role.Apps, role.Apps),
			"admins":          fmt.Sprintf("%v (type: %T)", role.Admins, role.Admins),
			"marshalled_json": string(roleMarshalled),
		})

		// Check for specific user ID format issues
		if len(role.Users) > 0 {
			tflog.Error(ctx, "[DEBUG] User IDs", map[string]interface{}{
				"user_ids":           fmt.Sprintf("%v", role.Users),
				"first_id_type":      fmt.Sprintf("%T", role.Users[0]),
				"first_id_value":     role.Users[0],
				"first_id_max_int32": role.Users[0] == 2147483647,
			})
		}

		// Try to decode error response more thoroughly
		if errResp, ok := resp.(map[string]interface{}); ok {
			if status, ok := errResp["status"].(map[string]interface{}); ok {
				if message, ok := status["message"].(string); ok {
					tflog.Error(ctx, "[DEBUG] API Error Message", map[string]interface{}{
						"message": message,
					})
				}
			}

			if errorObj, ok := errResp["error"]; ok {
				tflog.Error(ctx, "[DEBUG] API Error Object", map[string]interface{}{
					"error": fmt.Sprintf("%v", errorObj),
				})
			}

			// Log the full structure
			tflog.Error(ctx, "[DEBUG] Full Error Response Structure", map[string]interface{}{
				"response_structure": fmt.Sprintf("%#v", errResp),
			})
		}

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
