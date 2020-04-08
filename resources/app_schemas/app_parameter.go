package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
)

// AppParameter returns a key/value map of the various fields that make up
// the AppParameter field for a OneLogin App.
func AppParameter() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"param_key_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"param_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"label": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_attribute_mappings": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_attribute_macros": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"attributes_transformations": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"default_values": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"skip_if_blank": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"values": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"provisioned_entitlements": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"safe_entitlements_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

// InflateAppParameter takes a key/value map of interfaces and uses the fields to construct
// an AppParameter struct, a sub-field of a OneLogin App.
func InflateAppParameter(s map[string]interface{}) models.AppParameters {
	return models.AppParameters{
		ID:                        oltypes.Int32(int32(s["param_id"].(int))),
		Label:                     oltypes.String(s["label"].(string)),
		UserAttributeMappings:     oltypes.String(s["user_attribute_mappings"].(string)),
		UserAttributeMacros:       oltypes.String(s["user_attribute_macros"].(string)),
		AttributesTransformations: oltypes.String(s["attributes_transformations"].(string)),
		SkipIfBlank:               oltypes.Bool(s["skip_if_blank"].(bool)),
		Values:                    oltypes.String(s["values"].(string)),
		DefaultValues:             oltypes.String(s["default_values"].(string)),
		ProvisionedEntitlements:   oltypes.Bool(s["provisioned_entitlements"].(bool)),
		SafeEntitlementsEnabled:   oltypes.Bool(s["safe_entitlements_enabled"].(bool)),
	}

}
