package onelogin

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
		d.SetId("")
		return diag.Errorf("App does not exist: %s", err)
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

// appRoleLocks holds one mutex per app ID, serialising the read-modify-write
// that every attachment does on its app's role_ids.
//
// The API offers nothing finer than the whole array: an update replaces it with
// whatever is sent. Two attachments on the same app that both read before
// either writes each send the array they read plus their own role, and the
// later write drops the earlier one's. With for_each that is the normal case,
// not an edge one -- #272.
//
// This only coordinates attachments inside one provider process, which is one
// provider configuration in one run; an aliased configuration gets a process,
// and a set of locks, of its own. Anything else that changes an app's roles --
// the app resources' role_ids argument, onelogin_roles.apps, the OneLogin UI --
// can still overwrite them.
//
// Entries are never removed. A provider process lives for one operation and
// touches a bounded number of apps, so one mutex per app is not worth the
// bookkeeping of reclaiming.
var appRoleLocks sync.Map

// lockAppRoles takes the lock for appID and returns the function that releases
// it.
func lockAppRoles(appID int) func() {
	v, _ := appRoleLocks.LoadOrStore(appID, new(sync.Mutex))
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

func removeRoleFromApp(ctx context.Context, client *onelogin.OneloginSDK, appID int, roleID int) error {
	defer lockAppRoles(appID)()

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

	tflog.Info(ctx, "[UPDATED] Removed role from app", map[string]interface{}{
		"role_id": roleID,
		"app_id":  appID,
	})
	return nil
}

func attachRoleToApp(ctx context.Context, client *onelogin.OneloginSDK, appID int, roleID int) error {
	defer lockAppRoles(appID)()

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

	tflog.Info(ctx, "[UPDATED] Added role to app", map[string]interface{}{
		"role_id": roleID,
		"app_id":  appID,
	})
	return nil
}
