package appparametersschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func paramSet(names ...string) *schema.Set {
	list := make([]interface{}, 0, len(names))
	for _, name := range names {
		list = append(list, stateParam(map[string]interface{}{"param_key_name": name}))
	}
	return schema.NewSet(Hash, list)
}

// stateParam is a parameter the way the SDK hands one back from state or
// configuration: every attribute present, the ones nobody set at their zero
// value. The set hash runs over exactly this shape, so a test comparing hashes
// has to use it -- a missing bool or list does not hash like false or an empty
// list.
func stateParam(set map[string]interface{}) map[string]interface{} {
	param := map[string]interface{}{
		"param_key_name":             "",
		"param_id":                   0,
		"label":                      "",
		"user_attribute_mappings":    "",
		"user_attribute_macros":      "",
		"attributes_transformations": "",
		"default_values":             []interface{}{},
		"values":                     []interface{}{},
		"skip_if_blank":              false,
		"provisioned_entitlements":   false,
		"safe_entitlements_enabled":  false,
		"include_in_saml_assertion":  false,
	}
	for k, v := range set {
		param[k] = v
	}
	return param
}

// TestHash covers the set identity, and gh issue 274 with it. The SDK compares
// a set by its hash codes before it looks at any attribute, so an attribute the
// hash leaves out can never produce a diff: hashing the key name alone made a
// changed mapping plan nothing at all.
func TestHash(t *testing.T) {
	mapped := func(mapping string) map[string]interface{} {
		return stateParam(map[string]interface{}{
			"param_key_name":          "email",
			"label":                   "Email",
			"user_attribute_mappings": mapping,
		})
	}

	if Hash(mapped("email")) == Hash(mapped("username")) {
		t.Fatal("expected a changed mapping to change the hash, or no plan can ever see it")
	}
	if Hash(mapped("email")) != Hash(mapped("email")) {
		t.Fatal("expected the same parameter to hash the same")
	}

	withValues := func(values ...interface{}) map[string]interface{} {
		return stateParam(map[string]interface{}{"param_key_name": "groups", "default_values": values})
	}
	if Hash(withValues("one")) == Hash(withValues("one", "two")) {
		t.Fatal("expected a changed list to change the hash")
	}

	// param_id is Computed only. It is never in a configuration, so it
	// cannot be part of what identifies a parameter.
	withID := mapped("email")
	withID["param_id"] = 369278
	if Hash(withID) != Hash(mapped("email")) {
		t.Fatal("expected param_id to be left out of the hash")
	}
}

// TestRetainManaged covers the parameters a connector brings with it. The AWS
// connector alone adds saml_username, Role and RoleSessionName; they are in no
// configuration, and the app PUT merges rather than replaces, so recording them
// is a removal every plan proposes and no apply can settle.
func TestRetainManaged(t *testing.T) {
	fromAPI := []map[string]interface{}{
		{"param_key_name": "email", "label": "Email"},
		{"param_key_name": "firstname", "label": "First Name"},
		{"param_key_name": "saml_username", "label": "Amazon Username"},
		{"param_key_name": "https://aws.amazon.com/SAML/Attributes/Role", "label": "Role"},
	}

	t.Run("keeps only the configured parameters", func(t *testing.T) {
		out := RetainManaged(paramSet("email", "firstname"), fromAPI)

		if len(out) != 2 {
			t.Fatalf("expected the two configured parameters, got %v", out)
		}
		for _, param := range out {
			if name := param["param_key_name"]; name == "saml_username" {
				t.Fatalf("expected the connector's own parameters to be left out, got %v", name)
			}
		}
	})

	t.Run("drops a configured parameter the app no longer has", func(t *testing.T) {
		out := RetainManaged(paramSet("email", "department"), fromAPI)

		if len(out) != 1 || out[0]["param_key_name"] != "email" {
			t.Fatalf("expected only the parameter the API still reports, got %v", out)
		}
	})

	t.Run("writes everything on import", func(t *testing.T) {
		for _, prior := range []interface{}{paramSet(), nil, "not a set"} {
			out := RetainManaged(prior, fromAPI)

			if len(out) != len(fromAPI) {
				t.Fatalf("expected every parameter for prior %#v, got %v", prior, out)
			}
		}
	})
}
