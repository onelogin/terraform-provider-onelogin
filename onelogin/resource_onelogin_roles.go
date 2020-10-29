package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/role"
)

// Roles returns a resource with the CRUD methods and Terraform Schema defined
func Roles() *schema.Resource {
	return &schema.Resource{
		Create:   rolesCreate,
		Read:     rolesRead,
		Update:   rolesUpdate,
		Delete:   rolesDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   roleschema.Schema(),
	}
}

func rolesCreate(d *schema.ResourceData, m interface{}) error {
	role := roleschema.Inflate(map[string]interface{}{
		"name":   d.Get("name"),
		"apps":   d.Get("apps"),
		"users":  d.Get("users"),
		"admins": d.Get("admins"),
	})
	client := m.(*client.APIClient)
	err := client.Services.RolesV1.Create(&role)
	if err != nil {
		log.Println("[ERROR] There was a problem Creating the role")
		return err
	}
	log.Printf("[CREATED] Created role with id %d", *(role.ID))
	d.SetId(fmt.Sprintf("%d", *(role.ID)))
	return rolesRead(d, m)
}

func rolesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	uid, _ := strconv.Atoi(d.Id())
	role, err := client.Services.RolesV1.GetOne(int32(uid))
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
	d.Set("name", role.Name)
	d.Set("apps", role.Apps)
	d.Set("users", role.Users)
	d.Set("admins", role.Admins)
	return nil
}

func rolesUpdate(d *schema.ResourceData, m interface{}) error {
	var id interface{}
	id, _ = strconv.Atoi(d.Id())
	role := roleschema.Inflate(map[string]interface{}{
		"id":     id,
		"name":   d.Get("name"),
		"apps":   d.Get("apps"),
		"users":  d.Get("users"),
		"admins": d.Get("admins"),
	})
	client := m.(*client.APIClient)
	err := client.Services.RolesV1.Update(&role)
	if err != nil {
		log.Println("[ERROR] There was a problem Updating the role")
		return err
	}
	log.Printf("[UPDATED] Created role with id %d", *(role.ID))
	d.SetId(fmt.Sprintf("%d", *(role.ID)))
	return rolesRead(d, m)
}

func rolesDelete(d *schema.ResourceData, m interface{}) error {
	uid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	err := client.Services.RolesV1.Destroy(int32(uid))
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the role!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted role with %d", uid)
		d.SetId("")
	}

	return nil
}
