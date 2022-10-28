package onelogin

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
)

// UserRoleAttachment attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func UserRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: userRoleAttachmentCreate,
		Read:   userRoleAttachmentRead,
		Update: userRoleAttachmentUpdate,
		Delete: userRoleAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
  		"users": {
  			Type:     schema.TypeSet,
  			Required: true,
  			Elem:     &schema.Schema{Type: schema.TypeInt},
  		},
		},
	}
}

func userRoleAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	roleID := d.Get("role_id")
	users := d.Get("users")

	if appErr := attachRoleToUser(client, users, roleID); appErr != nil {
		return fmt.Errorf("Unable to attach role to app %s", appErr)
	}

	d.SetId(fmt.Sprintf("%d%d", roleID, users))
	return userRoleAttachmentRead(d, m)
}

func userRoleAttachmentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	users := d.Get("users").(int)
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

func userRoleAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	oldApp, newApp := d.GetChange("app_id")
	oldRole, newRole := d.GetChange("role_id")

	var err error
	if err = removeRoleFromUser(client, oldApp, oldRole); err != nil {
		return fmt.Errorf("Unable to remove role from app %s", err)
	}

	if err = attachRoleToUser(client, newApp, newRole); err != nil {
		return fmt.Errorf("Unable to attach role to app %s", err)
	}

	d.SetId(fmt.Sprintf("%d%d", newRole, newApp))
	return userRoleAttachmentRead(d, m)
}

func userRoleAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	appID := d.Get("app_id")
	roleID := d.Get("role_id")

	var err error
	if err = removeRoleFromUser(client, appID, roleID); err != nil {
		return fmt.Errorf("Unable to remove role from app %s", err)
	}
	d.SetId("")
	return nil
}

func removeRoleFromUser(client *client.APIClient, appID interface{}, roleID interface{}) error {
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

func attachRoleToUser(client *client.APIClient, appID interface{}, roleID interface{}) error {
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
