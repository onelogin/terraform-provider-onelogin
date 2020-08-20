package authserverschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/auth_server/configuration"
)

// Schema returns a key/value map of the various fields that make up an App at OneLogin.
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
func Inflate(s map[string]interface{}) (authservers.AuthServer, error) {
	var err error
	authServer := authservers.AuthServer{
		Name:        oltypes.String(s["name"].(string)),
		Description: oltypes.String(s["description"].(string)),
	}
	if s["configuration"] != nil {
		var conf authservers.AuthServerConfiguration
		conf = authserverconfigurationschema.Inflate(s["configuration"].([]interface{}))
		authServer.Configuration = &conf
	}
	return authServer, err
}
