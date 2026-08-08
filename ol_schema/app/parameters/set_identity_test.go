package appparametersschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func paramSet(names ...string) *schema.Set {
	list := make([]interface{}, 0, len(names))
	for _, name := range names {
		list = append(list, map[string]interface{}{"param_key_name": name})
	}
	return schema.NewSet(HashByKeyName, list)
}

// TestHashByKeyName covers the set identity. Every attribute but the key name
// is Computed, so hashing the whole element makes a parameter whose values are
// not yet known unmatchable against the one already in state, and each plan
// proposes replacing the set.
func TestHashByKeyName(t *testing.T) {
	configured := map[string]interface{}{
		"param_key_name": "email",
		"label":          "Email",
	}
	// The same parameter after a read, carrying what the API filled in.
	fromState := map[string]interface{}{
		"param_key_name":            "email",
		"label":                     "Email",
		"param_id":                  369278,
		"values":                    "user.email",
		"user_attribute_mappings":   "email",
		"include_in_saml_assertion": true,
	}

	if HashByKeyName(configured) != HashByKeyName(fromState) {
		t.Fatal("expected a parameter to hash the same before and after the API fills in its Computed attributes")
	}
	if HashByKeyName(configured) == HashByKeyName(map[string]interface{}{"param_key_name": "firstname"}) {
		t.Fatal("expected different key names to hash differently")
	}
	if HashByKeyName("not a parameter") != 0 {
		t.Fatal("expected an unexpected shape to hash to zero rather than panic")
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
