package appparametersschema

import (
	"context"
	"reflect"
	"testing"
)

func upgrade(t *testing.T, parameters ...map[string]interface{}) []interface{} {
	t.Helper()

	raw := make([]interface{}, 0, len(parameters))
	for _, p := range parameters {
		raw = append(raw, p)
	}

	out, err := UpgradeParameterValuesV0(context.Background(), map[string]interface{}{
		"id":         "1502449",
		"name":       "An App",
		"parameters": raw,
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	upgraded, ok := out["parameters"].([]interface{})
	if !ok {
		t.Fatalf("expected parameters to survive, got %#v", out["parameters"])
	}
	return upgraded
}

func field(t *testing.T, parameter interface{}, name string) interface{} {
	t.Helper()
	return parameter.(map[string]interface{})[name]
}

// TestUpgradeParameterValuesV0 covers the migration from the string these
// fields used to be to the list they are now.
//
// The rules follow what the API does, not what the string looks like: it keeps
// exactly what it is sent and returns it unchanged, so a stored "a,b" is one
// value that happens to contain a comma.
func TestUpgradeParameterValuesV0(t *testing.T) {
	t.Run("a value becomes a list of one", func(t *testing.T) {
		out := upgrade(t, map[string]interface{}{
			"param_key_name": "groups",
			"values":         "memberOf",
			"default_values": "everyone",
		})

		if got := field(t, out[0], "values"); !reflect.DeepEqual(got, []interface{}{"memberOf"}) {
			t.Fatalf("expected [memberOf], got %#v", got)
		}
		if got := field(t, out[0], "default_values"); !reflect.DeepEqual(got, []interface{}{"everyone"}) {
			t.Fatalf("expected [everyone], got %#v", got)
		}
	})

	t.Run("a comma is part of the value, not a separator", func(t *testing.T) {
		// The API stores "a,b" and returns "a,b". Splitting it here would
		// change what the app sends to its service provider, silently, during
		// an upgrade nobody asked to be creative.
		out := upgrade(t, map[string]interface{}{
			"param_key_name": "groups",
			"default_values": "one,two",
		})

		if got := field(t, out[0], "default_values"); !reflect.DeepEqual(got, []interface{}{"one,two"}) {
			t.Fatalf("expected the comma to survive inside a single value, got %#v", got)
		}
	})

	t.Run("an empty string becomes nothing", func(t *testing.T) {
		// Inflate has always skipped an empty string, so [""] would start
		// sending a value where none was sent before -- and the API keeps an
		// empty string as a value and hands it back.
		out := upgrade(t, map[string]interface{}{
			"param_key_name": "groups",
			"values":         "",
		})

		if got := field(t, out[0], "values"); !reflect.DeepEqual(got, []interface{}{}) {
			t.Fatalf("expected an empty list, got %#v", got)
		}
	})

	t.Run("leaves a field that was never set alone", func(t *testing.T) {
		out := upgrade(t, map[string]interface{}{"param_key_name": "groups"})

		if _, present := out[0].(map[string]interface{})["values"]; present {
			t.Fatal("expected an absent field to stay absent")
		}
	})

	t.Run("is safe to run twice", func(t *testing.T) {
		// Already a list means the state has been through this once, or was
		// written by a newer provider. Either way there is nothing to do.
		out := upgrade(t, map[string]interface{}{
			"param_key_name": "groups",
			"values":         []interface{}{"memberOf"},
		})

		if got := field(t, out[0], "values"); !reflect.DeepEqual(got, []interface{}{"memberOf"}) {
			t.Fatalf("expected the list to be left alone, got %#v", got)
		}
	})

	t.Run("carries the rest of the parameter and the resource across", func(t *testing.T) {
		out := upgrade(t, map[string]interface{}{
			"param_key_name": "groups",
			"label":          "Groups",
			"param_id":       369280,
			"values":         "memberOf",
		})

		if got := field(t, out[0], "label"); got != "Groups" {
			t.Fatalf("expected the label to survive, got %#v", got)
		}
		if got := field(t, out[0], "param_id"); got != 369280 {
			t.Fatalf("expected param_id to survive, got %#v", got)
		}
	})

	t.Run("does not mind a resource with no parameters", func(t *testing.T) {
		for _, state := range []map[string]interface{}{
			{"id": "1", "name": "An App"},
			{"id": "1", "parameters": nil},
			{"id": "1", "parameters": "not a list"},
		} {
			if _, err := UpgradeParameterValuesV0(context.Background(), state, nil); err != nil {
				t.Fatalf("expected %#v to pass through, got %v", state, err)
			}
		}
	})
}

// TestSchemaV0 pins the shape the upgrader decodes old state with. If this
// drifts from what an earlier provider actually wrote, refresh fails on state
// the upgrade was meant to rescue.
func TestSchemaV0(t *testing.T) {
	v0 := SchemaV0()

	for _, field := range []string{"values", "default_values"} {
		if v0[field] == nil {
			t.Fatalf("expected v0 to declare %q", field)
		}
		if got := v0[field].Type.String(); got != "TypeString" {
			t.Fatalf("expected %q to be a string in v0, got %s", field, got)
		}
	}

	// Everything else comes from the current schema, so a field added later is
	// carried into v0 without anyone remembering to.
	if len(v0) != len(Schema()) {
		t.Fatalf("expected v0 to have the same fields as the current schema, got %d against %d", len(v0), len(Schema()))
	}
	if v0["param_key_name"] == nil || !v0["param_key_name"].Required {
		t.Fatal("expected the unchanged fields to come through from the current schema")
	}
}
