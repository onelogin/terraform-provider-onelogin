package resources

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func OneloginApps() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
}
