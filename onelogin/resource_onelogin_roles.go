package onelogin

import (
	"context"
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
		"users":  d.Get("users"),
		"admins": d.Get("admins"),
	})

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating role", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	result, err := client.CreateRole(role)
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

	result, err := client.GetRoleByID(rid, &roleschema.RoleQuery{})
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

	role := roleschema.Inflate(map[string]interface{}{
		"id":     d.Id(),
		"name":   d.Get("name"),
		"apps":   d.Get("apps"),
		"users":  d.Get("users"),
		"admins": d.Get("admins"),
	})

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[UPDATE] Updating role", map[string]interface{}{
		"id": rid,
	})

	_, err := client.UpdateRole(rid, *role, map[string]string{})
	if err != nil {
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
		return client.DeleteRole(rid, map[string]string{})
	}, "Role")
}
