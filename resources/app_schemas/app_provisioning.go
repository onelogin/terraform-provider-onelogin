package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// AppProvisioning returns a key/value map of the various fields that make up
// the AppProvisioning field for a OneLogin App.
func AppProvisioning() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}

// InflateAppProvisioning takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func InflateAppProvisioning(s *map[string]interface{}) *models.AppProvisioning {
	return &models.AppProvisioning{
		Enabled: oltypes.Bool((*s)["enabled"].(bool)),
	}
}
