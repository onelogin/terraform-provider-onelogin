package appparametersschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up
// the Parameters field for a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"param_key_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"param_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"label": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_attribute_mappings": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_attribute_macros": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"attributes_transformations": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"default_values": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"skip_if_blank": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"values": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"provisioned_entitlements": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"safe_entitlements_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"include_in_saml_assertion": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
}

// Inflate takes a map of interfaces and uses the fields to construct
// a Parameter instance.
func Inflate(s map[string]interface{}) models.Parameter {
	out := models.Parameter{}
	var b, notNil bool
	var d int
	var st string

	if st, notNil = s["label"].(string); notNil {
		out.Label = st
	}

	if st, notNil = s["user_attribute_mappings"].(string); notNil {
		out.UserAttributeMappings = st
	}

	if st, notNil = s["user_attribute_macros"].(string); notNil {
		out.UserAttributeMacros = st
	}

	if st, notNil = s["attributes_transformations"].(string); notNil {
		out.AttributesTransformations = st
	}

	if st, notNil = s["values"].(string); notNil {
		out.Values = st
	}

	if st, notNil = s["default_values"].(string); notNil {
		out.DefaultValues = st
	}

	if b, notNil = s["skip_if_blank"].(bool); notNil {
		out.SkipIfBlank = b
	}

	if b, notNil = s["provisioned_entitlements"].(bool); notNil {
		out.ProvisionedEntitlements = b
	}

	if b, notNil = s["include_in_saml_assertion"].(bool); notNil {
		out.IncludeInSamlAssertion = b
	}

	if d, notNil = s["param_id"].(int); d != 0 && notNil {
		out.ID = d
	}
	return out
}

// Flatten takes a map of Parameter instances and returns an array of maps
func Flatten(params map[string]models.Parameter) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	for k, v := range params {
		param := map[string]interface{}{
			"param_key_name":             k,
			"param_id":                   v.ID,
			"label":                      v.Label,
			"user_attribute_mappings":    v.UserAttributeMappings,
			"user_attribute_macros":      v.UserAttributeMacros,
			"attributes_transformations": v.AttributesTransformations,
			"skip_if_blank":              v.SkipIfBlank,
			"values":                     v.Values,
			"default_values":             v.DefaultValues,
			"provisioned_entitlements":   v.ProvisionedEntitlements,
			"include_in_saml_assertion":  v.IncludeInSamlAssertion,
		}
		out = append(out, param)
	}
	return out
}

// FlattenV4 takes a map of interface{} and returns an array of maps for V4 SDK compatibility
func FlattenV4(params map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	for k, v := range params {
		if paramMap, ok := v.(map[string]interface{}); ok {
			param := map[string]interface{}{
				"param_key_name": k,
			}

			if id, ok := paramMap["id"].(float64); ok {
				param["param_id"] = int(id)
			}

			if val, ok := paramMap["label"].(string); ok {
				param["label"] = val
			}

			if val, ok := paramMap["user_attribute_mappings"].(string); ok {
				param["user_attribute_mappings"] = val
			}

			if val, ok := paramMap["user_attribute_macros"].(string); ok {
				param["user_attribute_macros"] = val
			}

			if val, ok := paramMap["attributes_transformations"].(string); ok {
				param["attributes_transformations"] = val
			}

			if val, ok := paramMap["skip_if_blank"].(bool); ok {
				param["skip_if_blank"] = val
			}

			if val, ok := paramMap["values"].(string); ok {
				param["values"] = val
			}

			if val, ok := paramMap["default_values"].(string); ok {
				param["default_values"] = val
			}

			if val, ok := paramMap["provisioned_entitlements"].(bool); ok {
				param["provisioned_entitlements"] = val
			}

			if val, ok := paramMap["include_in_saml_assertion"].(bool); ok {
				param["include_in_saml_assertion"] = val
			}

			out = append(out, param)
		}
	}
	return out
}
