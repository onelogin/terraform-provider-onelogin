package onelogin

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
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

	if appErr := updateUserRoleAttachment(client, users, roleID); appErr != nil {
		return fmt.Errorf("Unable to attach role to app %s", appErr)
	}

	d.SetId(fmt.Sprintf("%d", roleID))
	return userRoleAttachmentRead(d, m)
}

func userRoleAttachmentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	roleID := d.Get("role_id").(int)

	role, err := client.Services.RolesV1.GetOne(int32(roleID))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the role!")
		log.Println(err)
		return err
	}
	if role == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading role with %d", *(role.ID))
	d.Set("role_id", roleID)
	d.Set("users", role.Users)
	return nil
}

func userRoleAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	oldRole, newRole:= d.GetChange("role_id")
	oldUsers, newUsers:= d.GetChange("users")

	var err error
  if err = removeUserRoleAttachment(client, oldUsers, oldRole); err != nil {
		  return fmt.Errorf("Unable to delete mapping %s", err)
  }

	if err = updateUserRoleAttachment(client, newUsers, newRole); err != nil {
		return fmt.Errorf("Unable to update mapping %s", err)
	}

	d.SetId(fmt.Sprintf("%d", newRole))
	return userRoleAttachmentRead(d, m)
}

func userRoleAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	roleID := d.Get("role_id")
	users := d.Get("users")

	var err error
	if err = removeUserRoleAttachment(client, users, roleID); err != nil {
		return fmt.Errorf("Unable to remove role from users %s", err)
	}
	d.SetId("")
	return nil
}

func updateUserRoleAttachment(client *client.APIClient, userIDs interface{}, roleID interface{}) error {
  /*
    should be replaced with method provided by onelogin-go-sdk after sdk v3 is released
  */
  payload := make([]int32, userIDs.(*schema.Set).Len())
  for i, userID := range userIDs.(*schema.Set).List() {
    payload[i] = int32(userID.(int))
  }

  svc := client.Services.RolesV1
	_, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/users", svc.Endpoint, int32(roleID.(int))),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    &payload,
	})
	return err
}

func removeUserRoleAttachment(client *client.APIClient, userIDs interface{}, roleID interface{}) error {
  /*
    should be replaced with method provided by onelogin-go-sdk after sdk v3 is released
  */
  payload := make([]int32, userIDs.(*schema.Set).Len())
  for i, userID := range userIDs.(*schema.Set).List() {
    payload[i] = int32(userID.(int))
  }

  svc := client.Services.RolesV1
	_, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d/users", svc.Endpoint, int32(roleID.(int))),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    &payload,
	})
	return err
}
