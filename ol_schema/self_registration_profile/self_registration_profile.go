package selfregistrationprofile

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema returns a map of the schema for the self registration profile resource
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the self-registration profile",
		},
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The URL path for the self-registration profile",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether the self-registration profile is enabled",
		},
		"moderated": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the self-registration profile requires moderation",
		},
		"default_role_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The default role ID to assign to users who register through this profile",
		},
		"default_group_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The default group ID to assign to users who register through this profile",
		},
		"helptext": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Help text displayed on the registration page",
		},
		"thankyou_message": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Thank you message displayed after registration",
		},
		"domain_blacklist": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of domains to blacklist",
		},
		"domain_whitelist": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of domains to whitelist",
		},
		"domain_list_strategy": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Domain list strategy: 0 for blacklist, 1 for whitelist",
		},
		"email_verification_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Email MagicLink",
			Description: "Email verification type: 'Email MagicLink' or 'Email OTP'",
		},
		"fields": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Custom fields for the self-registration profile",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "ID of the field",
					},
					"custom_attribute_id": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "ID of the custom attribute",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the field",
					},
				},
			},
		},
	}
}

// Inflate takes a map of interfaces and returns a SelfRegistrationProfile struct
func Inflate(s map[string]interface{}) (map[string]interface{}, error) {
	return s, nil
}
