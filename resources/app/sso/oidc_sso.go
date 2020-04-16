package sso

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// AppConfiguration returns a key/value map of the various fields that make up
// the AppConfiguration field for a OneLogin App.
func OIDCSSOSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"client_secret": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
