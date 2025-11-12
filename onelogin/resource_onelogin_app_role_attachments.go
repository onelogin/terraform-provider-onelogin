package onelogin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// AppRoleAttachment attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func AppRoleAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: appRoleAttachmentCreate,
		ReadContext:   appRoleAttachmentRead,
		UpdateContext: appRoleAttachmentUpdate,
		DeleteContext: appRoleAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func appRoleAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	roleID := d.Get("role_id").(int)
	appID := d.Get("app_id").(int)

	if appErr := attachRoleToApp(ctx, client, appID, roleID); appErr != nil {
		return diag.Errorf("Unable to attach role to app: %s", appErr)
	}

	d.SetId(fmt.Sprintf("%d%d", roleID, appID))
	return appRoleAttachmentRead(ctx, d, m)
}

func appRoleAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	appID := d.Get("app_id").(int)
	roleID := d.Get("role_id").(int)

	result, err := client.GetAppByID(appID, nil)
	if err != nil {
		// Check if this is a 404 (resource not found)
		if utils.IsNotFoundError(err) {
			tflog.Info(ctx, "[NOT FOUND] App not found for role attachment", map[string]interface{}{
				"app_id": appID,
			})
			d.SetId("")
			return nil
		}
		// For other errors, return the error
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "App Role Attachment", d.Id())
	}

	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("Failed to parse app response")
	}

	roleIdsInterface, hasRoles := appMap["role_ids"].([]interface{})
	if !hasRoles {
		d.SetId("")
		return diag.Errorf("App %d does not have any roles", appID)
	}

	for _, rIDInterface := range roleIdsInterface {
		rID := int(rIDInterface.(float64))
		if rID == roleID {
			d.Set("role_id", rID)
			d.Set("app_id", appID)
			return nil
		}
	}

	d.SetId("")
	return diag.Errorf("App %d does not have role %d", appID, roleID)
}

func appRoleAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	oldApp, newApp := d.GetChange("app_id")
	oldRole, newRole := d.GetChange("role_id")

	var err error
	if err = removeRoleFromApp(ctx, client, oldApp.(int), oldRole.(int)); err != nil {
		return diag.Errorf("Unable to remove role from app: %s", err)
	}

	if err = attachRoleToApp(ctx, client, newApp.(int), newRole.(int)); err != nil {
		return diag.Errorf("Unable to attach role to app: %s", err)
	}

	d.SetId(fmt.Sprintf("%d%d", newRole, newApp))
	return appRoleAttachmentRead(ctx, d, m)
}

func appRoleAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	appID := d.Get("app_id").(int)
	roleID := d.Get("role_id").(int)

	var err error
	if err = removeRoleFromApp(ctx, client, appID, roleID); err != nil {
		return diag.Errorf("Unable to remove role from app: %s", err)
	}
	d.SetId("")
	return nil
}

func removeRoleFromApp(ctx context.Context, client *onelogin.OneloginSDK, appID int, roleID int) error {
	result, err := client.GetAppByID(appID, nil)
	if err != nil {
		return err
	}

	appMap, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Failed to parse app response")
	}

	roleIdsInterface, hasRoles := appMap["role_ids"].([]interface{})
	if !hasRoles {
		return fmt.Errorf("App %d does not have any roles", appID)
	}

	// Create a new slice with all roles except the one to remove
	newRoleIDs := []int{}
	for _, rIDInterface := range roleIdsInterface {
		rID := int(rIDInterface.(float64))
		if rID != roleID {
			newRoleIDs = append(newRoleIDs, rID)
		}
	}

	// Update the app with the new role IDs
	appToUpdate := models.App{
		RoleIDs: &newRoleIDs,
	}

	_, err = client.UpdateApp(appID, appToUpdate)
	if err != nil {
		return err
	}

	tflog.Info(ctx, "[UPDATED] Removed role %d from app %d", roleID, appID)
	return nil
}

func attachRoleToApp(ctx context.Context, client *onelogin.OneloginSDK, appID int, roleID int) error {
	result, err := client.GetAppByID(appID, nil)
	if err != nil {
		return err
	}

	appMap, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Failed to parse app response")
	}

	// Get existing role IDs or initialize empty slice
	roleIDs := []int{}
	roleIdsInterface, hasRoles := appMap["role_ids"].([]interface{})
	if hasRoles {
		for _, rIDInterface := range roleIdsInterface {
			rID := int(rIDInterface.(float64))
			roleIDs = append(roleIDs, rID)
		}
	}

	// Add the new role ID
	roleIDs = append(roleIDs, roleID)

	// Update the app with the new role IDs
	appToUpdate := models.App{
		RoleIDs: &roleIDs,
	}

	_, err = client.UpdateApp(appID, appToUpdate)
	if err != nil {
		return err
	}

	tflog.Info(ctx, "[UPDATED] Added role %d to app %d", roleID, appID)
	return nil
}
