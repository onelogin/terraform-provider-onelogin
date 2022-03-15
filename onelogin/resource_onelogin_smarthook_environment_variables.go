package onelogin

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	smarthookenvironmentvariablesschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/environment_variable"
)

// SmarthookEnvironmentVariable returns a resource with the CRUD methods and Terraform Schema defined
func SmarthookEnvironmentVariables() *schema.Resource {
	return &schema.Resource{
		Create:   environmentVariablesCreate,
		Read:     environmentVariablesRead,
		Update:   environmentVariablesUpdate,
		Delete:   environmentVariablesDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   smarthookenvironmentvariablesschema.Schema(),
	}
}

func environmentVariablesCreate(d *schema.ResourceData, m interface{}) error {
	envVar := smarthookenvironmentvariablesschema.Inflate(map[string]interface{}{
		"name":  d.Get("name"),
		"value": d.Get("value"),
	})
	client := m.(*client.APIClient)
	fullEnvVar, err := client.Services.SmartHooksEnvVarsV1.Create(&envVar)
	if err != nil {
		log.Println("[ERROR] There was a problem Creating the envVar")
		return err
	}
	log.Printf("[CREATED] Created envVar with id %s", *(fullEnvVar.ID))
	d.SetId(fmt.Sprintf("%s", *(fullEnvVar.ID)))
	return environmentVariablesRead(d, m)
}

func environmentVariablesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	envVar, err := client.Services.SmartHooksEnvVarsV1.GetOne(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the envVar!")
		log.Println(err)
		return err
	}
	if envVar == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading envVar with %s", *(envVar.ID))
	d.Set("name", envVar.Name)
	return nil
}

func environmentVariablesUpdate(d *schema.ResourceData, m interface{}) error {
	envVar := smarthookenvironmentvariablesschema.Inflate(map[string]interface{}{
		"id":    d.Id(),
		"value": d.Get("value"),
	})
	client := m.(*client.APIClient)
	fullEnvVar, err := client.Services.SmartHooksEnvVarsV1.Update(&envVar)
	if err != nil {
		log.Println("[ERROR] There was a problem Updating the envVar")
		return err
	}
	log.Printf("[UPDATED] Created envVar with id %s", *(fullEnvVar.ID))
	d.SetId(fmt.Sprintf("%s", *(fullEnvVar.ID)))
	return environmentVariablesRead(d, m)
}

func environmentVariablesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)

	err := client.Services.SmartHooksEnvVarsV1.Destroy(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the envVar!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted envVar with %s", d.Id())
		d.SetId("")
	}

	return nil
}
