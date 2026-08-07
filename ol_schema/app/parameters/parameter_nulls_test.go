package appparametersschema

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestInflateDoesNotSendEmptyStringsForUnsetFields covers issue #237 for app
// parameters.
//
// These fields are interface{} on the model. An interface holding "" is not
// empty to encoding/json -- only a nil interface is -- so omitempty did not
// protect them, and assigning whatever Terraform reported sent `"values": ""`
// for parameters the API had returned as null.
//
// Asserted on the marshalled body, since that is where the defect lives: the
// struct looks reasonable either way.
func TestInflateDoesNotSendEmptyStringsForUnsetFields(t *testing.T) {
	marshal := func(t *testing.T, s map[string]interface{}) string {
		t.Helper()
		b, err := json.Marshal(Inflate(s))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		return string(b)
	}

	t.Run("omits the string fields when unset", func(t *testing.T) {
		got := marshal(t, map[string]interface{}{"label": "Groups"})

		for _, key := range []string{"values", "user_attribute_macros", "user_attribute_mappings", "attributes_transformations"} {
			if strings.Contains(got, `"`+key+`":""`) {
				t.Fatalf("expected %s to be omitted rather than sent as \"\", got %s", key, got)
			}
		}
	})

	t.Run("sends default_values as null rather than empty string", func(t *testing.T) {
		// default_values carries no omitempty, so it is always serialised. A
		// nil interface makes that an explicit null, which is the state the
		// API was already in -- an empty string is a different value.
		got := marshal(t, map[string]interface{}{"label": "Groups"})

		if strings.Contains(got, `"default_values":""`) {
			t.Fatalf("expected null rather than an empty string, got %s", got)
		}
		if !strings.Contains(got, `"default_values":null`) {
			t.Fatalf("expected an explicit null, got %s", got)
		}
	})

	t.Run("still sends values that are set", func(t *testing.T) {
		got := marshal(t, map[string]interface{}{
			"label":                 "Groups",
			"values":                "memberOf",
			"default_values":        "everyone",
			"user_attribute_macros": "{$user.email}",
		})

		for _, want := range []string{`"values":"memberOf"`, `"default_values":"everyone"`, `"user_attribute_macros":"{$user.email}"`} {
			if !strings.Contains(got, want) {
				t.Fatalf("expected %s in %s", want, got)
			}
		}
	})

	t.Run("label survives, it is not one of the guarded fields", func(t *testing.T) {
		if got := marshal(t, map[string]interface{}{"label": "Groups"}); !strings.Contains(got, `"label":"Groups"`) {
			t.Fatalf("expected the label to be sent, got %s", got)
		}
	})
}
