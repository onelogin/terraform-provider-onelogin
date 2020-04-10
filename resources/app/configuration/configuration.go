package configuration

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type ConfigurationSchema func() map[string]*schema.Schema

func AddConfigurationSchema(parentSchema *map[string]*schema.Schema, configSchema ConfigurationSchema) {
	(*parentSchema)["configuration"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: configSchema(),
		},
	}
}
