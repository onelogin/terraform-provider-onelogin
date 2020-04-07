package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
)

func AppProvisioning() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}

func InflateAppProvisioning(s map[string]interface{}) *models.AppProvisioning {
	e := s["enabled"].(bool)

	return &models.AppProvisioning{
		Enabled: &e,
	}
}
