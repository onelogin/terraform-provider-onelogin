package onelogin

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	privilegeschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/privilege"
)

// Privileges returns a resource with the CRUD methods and Terraform Schema defined
func Privileges() *schema.Resource {
	privilegeSchema := privilegeschema.Schema()
	return &schema.Resource{
		CreateContext: privilegeCreate,
		ReadContext:   privilegeRead,
		UpdateContext: privilegeUpdate,
		DeleteContext: privilegeDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        privilegeSchema,
	}
}

// privilegeCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create a privilege with its sub-resources
func privilegeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	privilege, err := privilegeschema.Inflate(map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"user_ids":    d.Get("user_ids"),
		"role_ids":    d.Get("role_ids"),
		"privilege":   d.Get("privilege"),
	})
	if err != nil {
		return diag.Errorf("unable to inflate privilege: %v", err)
	}

	result, err := client.CreatePrivilege(privilege)
	if err != nil {
		return diag.Errorf("error creating privilege: %v", err)
	}

	// Extract the privilege ID from the response
	privilegeMap, ok := result.(map[string]interface{})
	if !ok || privilegeMap["id"] == nil {
		return diag.Errorf("failed to parse privilege creation response or privilege ID not found in response")
	}

	privilegeID := privilegeMap["id"].(string)
	d.SetId(privilegeID)
	log.Printf("[CREATED] Created privilege with id %s", privilegeID)

	return privilegeRead(ctx, d, m)
}

// privilegeRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read a privilege with its sub-resources
func privilegeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	result, err := client.GetPrivilege(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the privilege: %v", err)
		return diag.FromErr(err)
	}

	// Check if the resource was not found
	if result == nil {
		d.SetId("")
		return nil
	}

	// Parse the response
	privilegeMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse privilege response")
	}

	log.Printf("[READ] Reading privilege with id %s", d.Id())

	if privilegeMap["name"] != nil {
		d.Set("name", privilegeMap["name"])
	}

	if privilegeMap["description"] != nil {
		d.Set("description", privilegeMap["description"])
	}

	// Handle user_ids
	if privilegeMap["user_ids"] != nil {
		d.Set("user_ids", privilegeMap["user_ids"])
	}

	// Handle role_ids
	if privilegeMap["role_ids"] != nil {
		d.Set("role_ids", privilegeMap["role_ids"])
	}

	// Handle privilege data
	if privilegeMap["privilege"] != nil {
		privilegeData, ok := privilegeMap["privilege"].(map[string]interface{})
		if ok {
			// Process statements
			statements := []map[string]interface{}{}
			if stmts, ok := privilegeData["Statement"].([]interface{}); ok {
				for _, s := range stmts {
					stmt, ok := s.(map[string]interface{})
					if ok {
						statements = append(statements, map[string]interface{}{
							"effect": stmt["Effect"],
							"action": stmt["Action"],
							"scope":  stmt["Scope"],
						})
					}
				}
			}

			d.Set("privilege", []map[string]interface{}{
				{
					"version":   privilegeData["version"],
					"statement": statements,
				},
			})
		}
	}

	return nil
}

// privilegeUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update a privilege and its sub-resources
func privilegeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	privilege, err := privilegeschema.Inflate(map[string]interface{}{
		"id":          d.Id(),
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"user_ids":    d.Get("user_ids"),
		"role_ids":    d.Get("role_ids"),
		"privilege":   d.Get("privilege"),
	})
	if err != nil {
		return diag.Errorf("unable to inflate privilege: %v", err)
	}

	_, err = client.UpdatePrivilege(d.Id(), privilege)
	if err != nil {
		return diag.Errorf("error updating privilege: %v", err)
	}

	log.Printf("[UPDATED] Updated privilege with id %s", d.Id())
	return privilegeRead(ctx, d, m)
}

// privilegeDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete a privilege and its sub-resources
func privilegeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	_, err := client.DeletePrivilege(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the privilege: %v", err)
		return diag.FromErr(err)
	}

	log.Printf("[DELETED] Deleted privilege with id %s", d.Id())
	d.SetId("")

	return nil
}
