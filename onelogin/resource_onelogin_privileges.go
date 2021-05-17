package onelogin

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	privilegeschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/privilege"
)

// privileges attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func Privileges() *schema.Resource {
	privilegeSchema := privilegeschema.Schema()
	return &schema.Resource{
		Create:   privilegeCreate,
		Read:     privilegeRead,
		Update:   privilegeUpdate,
		Delete:   privilegeDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   privilegeSchema,
	}
}

// privilegeCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an privilege with its sub-resources
func privilegeCreate(d *schema.ResourceData, m interface{}) error {
	privilege, err := privilegeschema.Inflate(map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"user_ids":    d.Get("user_ids"),
		"role_ids":    d.Get("role_ids"),
		"privilege":   d.Get("privilege"),
	})
	if err != nil {
		log.Println("Unable to inflate privilege", err)
		return err
	}
	client := m.(*client.APIClient)
	err = client.Services.PrivilegesV1.Create(&privilege)
	if err != nil {
		log.Println("[ERROR] There was a problem creating the privilege!", err)
		return err
	}
	log.Printf("[CREATED] Created privilege with %s", *(privilege.ID))

	d.SetId(*(privilege.ID))
	return privilegeRead(d, m)
}

// privilegeRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an privilege with its sub-resources
func privilegeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	privilege, err := client.Services.PrivilegesV1.GetOne(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the privilege!")
		log.Println(err)
		return err
	}
	if privilege == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading privilege with %s", *(privilege.ID))

	d.Set("name", privilege.Name)
	d.Set("description", privilege.Description)
	d.Set("user_ids", privilege.UserIDs)
	d.Set("role_ids", privilege.RoleIDs)
	d.Set("privilege", privilegeschema.FlattenPrivilegeData(*privilege.Privilege))
	return nil
}

// privilegeUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an privilege and its sub-resources
func privilegeUpdate(d *schema.ResourceData, m interface{}) error {
	privilege, err := privilegeschema.Inflate(map[string]interface{}{
		"id":          d.Id(),
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"user_ids":    d.Get("user_ids"),
		"role_ids":    d.Get("role_ids"),
		"privilege":   d.Get("privilege"),
	})
	if err != nil {
		log.Println("Unable to inflate privilege", err)
		return err
	}
	client := m.(*client.APIClient)

	err = client.Services.PrivilegesV1.Update(&privilege)
	if err != nil {
		log.Println("[ERROR] There was a problem updating the privilege!", err)
		return err
	}

	log.Printf("[UPDATED] Updated privilege with %s", *(privilege.ID))
	d.SetId(*(privilege.ID))
	return privilegeRead(d, m)
}

// privilegeDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an privilege and its sub-resources
func privilegeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	err := client.Services.PrivilegesV1.Destroy(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the privilege!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted privilege with %s", d.Id())
		d.SetId("")
	}

	return nil
}
