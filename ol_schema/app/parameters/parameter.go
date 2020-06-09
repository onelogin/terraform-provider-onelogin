package appparametersschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
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
// an AppParameter instance.
func Inflate(s map[string]interface{}) apps.AppParameters {
	out := apps.AppParameters{}
	var b, notNil bool
	var d int
	var st string

	if st, notNil = s["label"].(string); notNil {
		out.Label = oltypes.String(st)
	}

	if st, notNil = s["user_attribute_mappings"].(string); notNil {
		out.UserAttributeMappings = oltypes.String(st)
	}

	if st, notNil = s["user_attribute_macros"].(string); notNil {
		out.UserAttributeMacros = oltypes.String(st)
	}

	if st, notNil = s["attributes_transformations"].(string); notNil {
		out.AttributesTransformations = oltypes.String(st)
	}

	if st, notNil = s["values"].(string); notNil {
		out.Values = oltypes.String(st)
	}

	if st, notNil = s["default_values"].(string); notNil {
		out.DefaultValues = oltypes.String(st)
	}

	if b, notNil = s["skip_if_blank"].(bool); notNil {
		out.SkipIfBlank = oltypes.Bool(b)
	}

	if b, notNil = s["provisioned_entitlements"].(bool); notNil {
		out.ProvisionedEntitlements = oltypes.Bool(b)
	}

	if b, notNil = s["safe_entitlements_enabled"].(bool); notNil {
		out.SafeEntitlementsEnabled = oltypes.Bool(b)
	}

	if b, notNil = s["include_in_saml_assertion"].(bool); notNil {
		out.IncludeInSamlAssertion = oltypes.Bool(b)
	}

	if d, notNil = s["param_id"].(int); d != 0 && notNil {
		out.ID = oltypes.Int32(int32(d))
	}
	return out
}

// Flatten takes a map of AppParamters instances and returns an array of maps
func Flatten(params map[string]apps.AppParameters) []map[string]interface{} {
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
			"safe_entitlements_enabled":  v.SafeEntitlementsEnabled,
			"include_in_saml_assertion":  v.IncludeInSamlAssertion,
		}
		out = append(out, param)
	}
	return out
}
