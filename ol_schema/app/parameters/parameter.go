package appparametersschema

import (
	"context"
	"fmt"

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
		// A list, because the API can hold several values here and a string
		// cannot say so. See SchemaV0 for what this used to be and
		// UpgradeParameterValuesV0 for how old state is carried across.
		"default_values": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"skip_if_blank": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"values": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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

// HashByKeyName identifies a parameter in the set by its key name alone.
//
// The default set hash covers every attribute, and all but param_key_name are
// Computed. A parameter whose values or user_attribute_macros the practitioner
// has not written is unknown at plan time, so the element hashes to something
// Terraform cannot match against the one already in state, and every plan
// proposes replacing the whole set. Hashing the key name is knowable from the
// configuration alone, so the element matches and the Computed attributes are
// then compared one at a time, which is what Optional+Computed is for.
//
// param_key_name is the API's own identity for a parameter: the endpoint
// returns parameters as an object keyed by it, and Inflate builds that object
// back, so two parameters cannot share one.
func HashByKeyName(v interface{}) int {
	param, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}
	name, _ := param["param_key_name"].(string)
	return schema.HashString(name)
}

// RetainManaged narrows flattened parameters to the ones the practitioner is
// managing.
//
// A connector brings its own parameters -- the AWS connector alone adds
// saml_username, Role and RoleSessionName -- and the app endpoint returns them
// alongside the configured ones. Recording them puts elements in state that
// are in no configuration, which every plan then proposes removing and no
// apply can settle: the PUT merges rather than replaces, so they come straight
// back.
//
// The prior state's key names are the practitioner's. On create they are the
// planned set, on update the new one, and on refresh whatever the last apply
// wrote. An empty prior state is an import, where there is no configuration to
// respect yet and everything the app has is written.
func RetainManaged(prior interface{}, flattened []map[string]interface{}) []map[string]interface{} {
	set, ok := prior.(*schema.Set)
	if !ok || set.Len() == 0 {
		return flattened
	}

	managed := make(map[string]bool, set.Len())
	for _, raw := range set.List() {
		if param, ok := raw.(map[string]interface{}); ok {
			if name, ok := param["param_key_name"].(string); ok {
				managed[name] = true
			}
		}
	}

	out := make([]map[string]interface{}, 0, len(managed))
	for _, param := range flattened {
		if name, ok := param["param_key_name"].(string); ok && managed[name] {
			out = append(out, param)
		}
	}
	return out
}

// parameterValues turns the list the schema declares into the shape the API
// takes: nil for nothing, a bare string for one value, an array for several.
func parameterValues(raw interface{}) interface{} {
	values, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	out := make([]string, 0, len(values))
	for _, value := range values {
		if str, ok := value.(string); ok {
			out = append(out, str)
		}
	}

	switch len(out) {
	case 0:
		// Nothing at all, which is not the same as [""] -- the API keeps an
		// empty string as a value and returns it.
		return nil
	case 1:
		return out[0]
	default:
		return out
	}
}

// parameterValueList turns whatever the API reported into the list the schema
// declares. The same field is a string on one parameter and an array on
// another, decided by nothing but what was last written to it.
func parameterValueList(raw interface{}) []interface{} {
	switch typed := raw.(type) {
	case string:
		return []interface{}{typed}
	case []interface{}:
		out := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			out = append(out, fmt.Sprint(item))
		}
		return out
	case []string:
		out := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			out = append(out, item)
		}
		return out
	}
	return nil
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

	// These are interface{} on the model, and an interface holding "" is not
	// empty for encoding/json -- only a nil interface is. Assigning whatever
	// Terraform reports therefore sent `"values": ""` for a parameter the API
	// had returned as null, and the API stores what it is sent. An update
	// touching only the app description rewrote all four.
	//
	// Leaving them nil omits them, and the app endpoint merges its PUT rather
	// than replacing, so an untouched field keeps whatever it already had.
	// default_values carries no omitempty and so goes out as an explicit null,
	// which the API also treats as "no value" -- the state it was already in.
	if st, notNil = s["user_attribute_mappings"].(string); notNil && st != "" {
		out.UserAttributeMappings = st
	}

	if st, notNil = s["user_attribute_macros"].(string); notNil && st != "" {
		out.UserAttributeMacros = st
	}

	if st, notNil = s["attributes_transformations"].(string); notNil && st != "" {
		out.AttributesTransformations = st
	}

	// One value goes out as a bare string, several as an array.
	//
	// The API stores exactly what it is sent and hands it back unchanged --
	// "x" stays "x" and ["x"] stays ["x"], neither is normalised into the
	// other. Sending a single value as a one-element array would therefore
	// rewrite the stored shape of every parameter that has been a string until
	// now, which is all of them, for no gain. An array is used only where a
	// string cannot express the value.
	out.Values = parameterValues(s["values"])
	out.DefaultValues = parameterValues(s["default_values"])

	if b, notNil = s["skip_if_blank"].(bool); notNil {
		out.SkipIfBlank = b
	}

	if b, notNil = s["provisioned_entitlements"].(bool); notNil {
		out.ProvisionedEntitlements = b
	}

	if b, notNil = s["include_in_saml_assertion"].(bool); notNil {
		out.IncludeInSamlAssertion = b
	}

	// A pointer on the model, unlike the flags above it, so a false is sent
	// rather than dropped by omitempty. The schema has declared this field for
	// some time, but nothing carried it: the value was discarded on the way
	// out and never came back on read.
	if b, notNil = s["safe_entitlements_enabled"].(bool); notNil {
		enabled := b
		out.SafeEntitlementsEnabled = &enabled
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
			"values":                     parameterValueList(v.Values),
			"default_values":             parameterValueList(v.DefaultValues),
			"provisioned_entitlements":   v.ProvisionedEntitlements,
			"include_in_saml_assertion":  v.IncludeInSamlAssertion,
		}

		// Absent stays absent: nil is a parameter the API said nothing about,
		// which is not the same as one with the setting switched off.
		if v.SafeEntitlementsEnabled != nil {
			param["safe_entitlements_enabled"] = *v.SafeEntitlementsEnabled
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

			// Both shapes, because the API returns whichever was last
			// written. A null is left out entirely, keeping "never set"
			// distinct from "set to nothing".
			if val := parameterValueList(paramMap["values"]); val != nil {
				param["values"] = val
			}

			if val := parameterValueList(paramMap["default_values"]); val != nil {
				param["default_values"] = val
			}

			if val, ok := paramMap["provisioned_entitlements"].(bool); ok {
				param["provisioned_entitlements"] = val
			}

			if val, ok := paramMap["include_in_saml_assertion"].(bool); ok {
				param["include_in_saml_assertion"] = val
			}

			// A null fails the assertion and is left out, keeping "unset"
			// distinct from "switched off".
			if val, ok := paramMap["safe_entitlements_enabled"].(bool); ok {
				param["safe_entitlements_enabled"] = val
			}

			out = append(out, param)
		}
	}
	return out
}

// SchemaV0 is the parameter schema as it was before default_values and values
// became lists. The state upgrader needs it to decode state written by an
// earlier provider version -- without it Terraform tries to read a string as a
// list and the resource cannot be refreshed at all.
//
// Derived from the current schema rather than restated, so a field added later
// appears here too and only the two that actually changed are described twice.
func SchemaV0() map[string]*schema.Schema {
	v0 := Schema()
	for _, field := range []string{"values", "default_values"} {
		v0[field] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		}
	}
	return v0
}

// UpgradeParameterValuesV0 rewrites a parameter's values and default_values
// from the string they used to be into the list they are now.
//
// The rules come from what the API does rather than what the string looks like:
// it stores exactly what it is given and returns it unchanged, so a stored
// "a,b" is one value containing a comma, not two values.
//
//	""     -> []        nothing. Inflate has always skipped an empty string, so
//	                    turning it into [""] would start sending a value where
//	                    none was sent before -- and [""] is a value the API
//	                    keeps and returns.
//	"a,b"  -> ["a,b"]   one value. Splitting on the comma would change what the
//	                    app sends to its service provider, silently, during an
//	                    upgrade nobody asked to be creative.
//	absent -> absent    nothing to say.
func UpgradeParameterValuesV0(ctx context.Context, state map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	parameters, ok := state["parameters"].([]interface{})
	if !ok {
		return state, nil
	}

	for _, raw := range parameters {
		parameter, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		for _, field := range []string{"values", "default_values"} {
			value, present := parameter[field]
			if !present || value == nil {
				continue
			}

			str, ok := value.(string)
			if !ok {
				// Already a list, so this state has been through the upgrade
				// or was written by a newer version. Leave it be.
				continue
			}
			if str == "" {
				parameter[field] = []interface{}{}
				continue
			}
			parameter[field] = []interface{}{str}
		}
	}

	return state, nil
}
