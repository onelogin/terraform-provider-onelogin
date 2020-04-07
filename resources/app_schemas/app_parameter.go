package app_schemas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
)

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

func InflateAppParameter(s map[string]interface{}) models.AppParameters {

	pid := int32(s["param_id"].(int))
	lbl := s["label"].(string)
	uam := s["user_attribute_mappings"].(string)
	uac := s["user_attribute_macros"].(string)
	atr := s["attributes_transformations"].(string)
	sib := s["skip_if_blank"].(bool)
	val := s["values"].(string)
	dfv := s["default_values"].(string)
	pet := s["provisioned_entitlements"].(bool)
	see := s["safe_entitlements_enabled"].(bool)

	return models.AppParameters{
		ID:                        &pid,
		Label:                     &lbl,
		UserAttributeMappings:     &uam,
		UserAttributeMacros:       &uac,
		AttributesTransformations: &atr,
		SkipIfBlank:               &sib,
		Values:                    &val,
		DefaultValues:             &dfv,
		ProvisionedEntitlements:   &pet,
		SafeEntitlementsEnabled:   &see,
	}

}
