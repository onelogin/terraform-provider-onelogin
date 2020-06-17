package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/actions"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/user_mapping/conditions"
)

// UserMappings attaches additional configuration and sso schemas and
// returns a resource with the CRUD methods and Terraform Schema defined
func UserMappings() *schema.Resource {
	mappingSchema := usermappingschema.Schema()

	return &schema.Resource{
		Create:   userMappingCreate,
		Read:     userMappingRead,
		Update:   userMappingUpdate,
		Delete:   userMappingDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   mappingSchema,
	}
}

// userMappingCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an samlApp with its sub-resources
func userMappingCreate(d *schema.ResourceData, m interface{}) error {
	mappings := usermappingschema.Inflate(map[string]interface{}{
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"enabled":    d.Get("enabled"),
		"position":   d.Get("position"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	client := m.(*client.APIClient)
	mappingsResp, err := client.Services.UserMappingsV2.Create(&mappings)
	if err != nil {
		log.Println("[ERROR] There was a problem creating the user mapping!", err)
		return err
	}
	log.Printf("[CREATED] Created user mapping with %d", *(mappingsResp.ID))

	d.SetId(fmt.Sprintf("%d", *(mappingsResp.ID)))
	return userMappingRead(d, m)
}

// UserMappingRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an samlApp with its sub-resources
func userMappingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	umid, _ := strconv.Atoi(d.Id())
	mapping, err := client.Services.UserMappingsV2.GetOne(int32(umid))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the mapping!")
		log.Println(err)
		return err
	}
	if mapping == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading app with %d", *(mapping.ID))

	d.Set("name", mapping.Name)
	d.Set("match", mapping.Match)
	d.Set("enabled", mapping.Enabled)
	d.Set("position", mapping.Position)
	d.Set("conditions", usermappingconditionsschema.Flatten(mapping.Conditions))
	d.Set("actions", usermappingactionsschema.Flatten(mapping.Actions))

	return nil
}

// UserMappingUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an samlApp and its sub-resources
func userMappingUpdate(d *schema.ResourceData, m interface{}) error {
	userMapping := usermappingschema.Inflate(map[string]interface{}{
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"enabled":    d.Get("enabled"),
		"position":   d.Get("position"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})

	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	userMappingResp, err := client.Services.UserMappingsV2.Update(int32(aid), &userMapping)
	if err != nil {
		log.Println("[ERROR] There was a problem Updating the user mapping!", err)
		return err
	}
	if userMappingResp == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[UPDATED] Updated user mapping with %d", *(userMappingResp.ID))
	d.SetId(fmt.Sprintf("%d", *(userMappingResp.ID)))
	return userMappingRead(d, m)
}

// userMappingDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete a user mapping
func userMappingDelete(d *schema.ResourceData, m interface{}) error {
	aid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	err := client.Services.UserMappingsV2.Destroy(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the user mapping!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted user mapping with %d", aid)
		d.SetId("")
	}

	return nil
}
