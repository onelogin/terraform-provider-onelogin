package authserverschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	authserverconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/auth_server/configuration"
)

// Schema returns a key/value map of the various fields that make up an AuthServer at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"configuration": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: authserverconfigurationschema.Schema(),
			},
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin AuthServer.
func Inflate(s map[string]interface{}) (models.AuthServer, error) {
	var err error
	authServer := models.AuthServer{}

	// Handle basic fields
	if name, notNil := s["name"].(string); notNil {
		authServer.Name = &name
	}

	if desc, notNil := s["description"].(string); notNil {
		authServer.Description = &desc
	}

	// Handle ID if present
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			id32 := int32(id)
			authServer.ID = &id32
		}
	}

	// Handle configuration if present
	if s["configuration"] != nil {
		conf := authserverconfigurationschema.Inflate(s["configuration"].([]interface{}))
		authServer.Configuration = &conf
	}

	return authServer, err
}
