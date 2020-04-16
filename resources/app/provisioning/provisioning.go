package provisioning

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// AppProvisioning returns a key/value map of the various fields that make up
// the AppProvisioning field for a OneLogin App.
func ProvisioningSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}

// InflateAppProvisioning takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func InflateProvisioning(s *map[string]interface{}) *models.AppProvisioning {
	out := models.AppProvisioning{}
	if enb, notNil := (*s)["enabled"].(bool); notNil {
		out.Enabled = oltypes.Bool(enb)
	}
	return &out
}

func Flatten(prov *models.AppProvisioning) []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"enabled": *prov.Enabled,
		},
	}
}
