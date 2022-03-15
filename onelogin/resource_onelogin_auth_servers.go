package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	authserverschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/auth_server"
	authserverconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/auth_server/configuration"
)

// AuthServers returns a resource with the CRUD methods and Terraform Schema defined
func AuthServers() *schema.Resource {
	return &schema.Resource{
		Create:   authServersCreate,
		Read:     authServersRead,
		Update:   authServersUpdate,
		Delete:   authServersDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   authserverschema.Schema(),
	}
}

func authServersCreate(d *schema.ResourceData, m interface{}) error {
	AuthServer, _ := authserverschema.Inflate(map[string]interface{}{
		"name":          d.Get("name"),
		"description":   d.Get("description"),
		"configuration": d.Get("configuration"),
	})
	client := m.(*client.APIClient)
	err := client.Services.AuthServersV2.Create(&AuthServer)
	if err != nil {
		log.Println("[ERROR] There was a problem creating the AuthServer!", err)
		return err
	}
	log.Printf("[CREATED] Created AuthServer with %d", *(AuthServer.ID))

	d.SetId(fmt.Sprintf("%d", *(AuthServer.ID)))
	return authServersRead(d, m)
}

func authServersUpdate(d *schema.ResourceData, m interface{}) error {
	AuthServer, _ := authserverschema.Inflate(map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Get("name"),
		"description":   d.Get("description"),
		"configuration": d.Get("configuration"),
	})
	client := m.(*client.APIClient)
	err := client.Services.AuthServersV2.Update(&AuthServer)
	if err != nil {
		log.Println("[ERROR] There was a problem updating the AuthServer!", err)
		return err
	}
	log.Printf("[CREATED] Updated AuthServer with %d", *(AuthServer.ID))

	d.SetId(fmt.Sprintf("%d", *(AuthServer.ID)))
	return authServersRead(d, m)
}

func authServersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	uid, _ := strconv.Atoi(d.Id())
	authServer, err := client.Services.AuthServersV2.GetOne(int32(uid))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the AuthServer!")
		log.Println(err)
		return err
	}
	if authServer == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading AuthServer with %d", *(authServer.ID))

	d.Set("name", authServer.Name)
	d.Set("description", authServer.Description)
	d.Set("configuration", authserverconfigurationschema.Flatten(*authServer.Configuration))

	return nil
}

func authServersDelete(d *schema.ResourceData, m interface{}) error {
	uid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	err := client.Services.AuthServersV2.Destroy(int32(uid))
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the AuthServer!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted AuthServer with %d", uid)
		d.SetId("")
	}

	return nil
}
