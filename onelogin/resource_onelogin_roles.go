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
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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

// roleRead reads a role by ID from OneLogin using the direct GET /api/2/roles/{id}
// endpoint for the base object, then fetches membership via the sub-endpoints
// GET /api/2/roles/{id}/apps, /users, /admins.  This is O(roles) total requests
// rather than O(roles × pages) from the previous full-list pagination approach.
func roleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	rid, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "[READ] Reading role", map[string]interface{}{
		"id": rid,
	})

	// Fetch base role object (returns {"id": ..., "name": ...}).
	result, err := client.GetRoleByIDWithContext(ctx, rid, nil)
	if err != nil {
		if utils.IsNotFoundError(err) {
			tflog.Info(ctx, "[NOT FOUND] Role not found, removing from state", map[string]interface{}{
				"id": rid,
			})
			d.SetId("")
			return nil
		}
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Role", d.Id())
	}

	roleObj, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse role response")
	}
	d.Set("name", roleObj["name"])

	// Membership is not part of the base role object — it lives on three separate
	// sub-endpoints, each of which is paginated.
	memberAttrs := []struct {
		attr  string
		fetch memberFetcher
	}{
		{"apps", func(ctx context.Context, q *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error) {
			return client.GetRoleAppsWithPaginationAndContext(ctx, rid, q)
		}},
		{"users", func(ctx context.Context, q *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error) {
			return client.GetRoleUsersWithPaginationAndContext(ctx, rid, q)
		}},
		{"admins", func(ctx context.Context, q *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error) {
			return client.GetRoleAdminsWithPaginationAndContext(ctx, rid, q)
		}},
	}

	for _, m := range memberAttrs {
		ids, err := fetchAllMemberIDs(ctx, m.fetch)
		if err != nil {
			// A 404 here means the role was deleted between the base fetch and now.
			// Treat it the same as a missing base object rather than recording the
			// role as having no members.
			if utils.IsNotFoundError(err) {
				tflog.Info(ctx, "[NOT FOUND] Role disappeared while reading membership, removing from state", map[string]interface{}{
					"id":        rid,
					"attribute": m.attr,
				})
				d.SetId("")
				return nil
			}
			return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, fmt.Sprintf("Role %s", m.attr), d.Id())
		}
		if err := d.Set(m.attr, ids); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// memberFetcher retrieves one page of a role sub-endpoint (apps, users or admins).
type memberFetcher func(context.Context, *roleschema.RoleQuery) (interface{}, *models.PaginationInfo, error)

// rolePageLimit is the page size requested when walking role sub-endpoints.
const rolePageLimit = "100"

// fetchAllMemberIDs walks every page of a role sub-endpoint and returns the flat
// list of member IDs. These endpoints are paginated, so reading only the first
// page would write a truncated list into state and produce a permanent diff for
// any role with more members than fit on one page.
func fetchAllMemberIDs(ctx context.Context, fetch memberFetcher) ([]int, error) {
	ids := []int{}
	query := &roleschema.RoleQuery{Limit: rolePageLimit}

	for {
		result, pagination, err := fetch(ctx, query)
		if err != nil {
			return nil, err
		}
		ids = append(ids, extractMemberIDs(result)...)

		if pagination == nil || pagination.AfterCursor == "" {
			return ids, nil
		}
		// Guard against a server that keeps handing back the same cursor, which
		// would otherwise hang the plan rather than fail it.
		if pagination.AfterCursor == query.Cursor {
			return ids, fmt.Errorf("pagination stalled: role sub-endpoint repeated cursor %q", query.Cursor)
		}
		// Cursor and limit/page are mutually exclusive in the V2 API.
		query.Cursor = pagination.AfterCursor
		query.Limit, query.Page = "", ""
	}
}

// extractMemberIDs converts the slice of member objects returned by role sub-endpoints
// (apps, users, admins) into a flat []int of IDs.  Each element is an object such as
// {"id": 123, "name": "..."} — only the "id" field is required; others are ignored.
func extractMemberIDs(result interface{}) []int {
	if result == nil {
		return []int{}
	}
	items, ok := result.([]interface{})
	if !ok {
		return []int{}
	}
	ids := make([]int, 0, len(items))
	for _, item := range items {
		if obj, ok := item.(map[string]interface{}); ok {
			if id, ok := obj["id"].(float64); ok {
				ids = append(ids, int(id))
			}
		}
	}
	return ids
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
