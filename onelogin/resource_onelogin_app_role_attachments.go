package onelogin

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
)

// AppRoleAttachment attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func AppRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: appRoleAttachmentCreate,
		Read:   appRoleAttachmentRead,
		Update: appRoleAttachmentUpdate,
		Delete: appRoleAttachmentDelete,
		Schema: map[string]*schema.Schema{
      "role_id": {
  			Type:         schema.TypeInt,
  			Required:     true,
  		},
      "app_id": {
  			Type:     schema.TypeInt,
  			Required: true,
  		},
    },
	}
}

func appRoleAttachmentCreate(d *schema.ResourceData, m interface{}) error {
  client := m.(*client.APIClient)

  roleID := d.Get("role_id")
	appID := d.Get("app_id")

	if appErr := attachRoleToApp(client, appID, roleID); appErr != nil {
    return fmt.Errorf("Unable to attach role to app %s", appErr)
  }

  d.SetId(fmt.Sprintf("%d%d", roleID, appID))
	return appRoleAttachmentRead(d, m)
}

func appRoleAttachmentRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*client.APIClient)
  appID := int32(d.Get("app_id").(int))
  roleID := d.Get("role_id").(int)

	app, err := client.Services.AppsV2.GetOne(appID)
  if err != nil {
    d.SetId("")
    return fmt.Errorf("App does not exist %s", err)
  }
  for _, rID := range app.RoleIDs {
    if rID == roleID {
      d.Set("role_id", rID)
      d.Set("app_id", *app.ID)
      return nil
    }
  }
  d.SetId("")
  return fmt.Errorf("App %d does not have role %d", appID, roleID)
}

func appRoleAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
  client := m.(*client.APIClient)

  oldApp, newApp := d.GetChange("app_id")
  oldRole, newRole := d.GetChange("role_id")

  var err error
  if err = removeRoleFromApp(client, oldApp, oldRole); err != nil {
    return fmt.Errorf("Unable to remove role from app %s", err)
  }

  if err = attachRoleToApp(client, newApp, newRole); err != nil {
    return fmt.Errorf("Unable to attach role to app %s", err)
  }

  d.SetId(fmt.Sprintf("%d%d", newRole, newApp))
  return appRoleAttachmentRead(d, m)
}

func appRoleAttachmentDelete(d *schema.ResourceData, m interface{}) error {
  client := m.(*client.APIClient)

  appID := d.Get("app_id")
  roleID := d.Get("role_id")

  var err error
  if err = removeRoleFromApp(client, appID, roleID); err != nil {
    return fmt.Errorf("Unable to remove role from app %s", err)
  }
  d.SetId("")
  return nil
}

func removeRoleFromApp(client *client.APIClient, appID interface{}, roleID interface{}) error {
  app, err := client.Services.AppsV2.GetOne(int32(appID.(int)))
  if err != nil {
    return err
  }
  newRoleIDs := make([]int, 0)
  for _, rID := range app.RoleIDs {
    if rID != roleID {
      newRoleIDs = append(newRoleIDs, rID)
    }
  }
  app.RoleIDs = newRoleIDs
  app, err = client.Services.AppsV2.Update(app)
  if err != nil {
    return err
  }
  return nil
}

func attachRoleToApp(client *client.APIClient, appID interface{}, roleID interface{}) error {
  app, err := client.Services.AppsV2.GetOne(int32(appID.(int)))
  if err != nil {
    return err
  }
  app.RoleIDs = append(app.RoleIDs, roleID.(int))
  app, err = client.Services.AppsV2.Update(app)
  if err != nil {
    return err
  }
  return nil
}
