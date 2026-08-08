package smarthooksschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	smarthookoptions "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/options"
)

// optionsSet builds the shape Terraform actually hands to Inflate for a
// TypeSet of a nested resource. The existing Inflate tests pass a bare map,
// which is the shape a caller building input by hand uses -- and is why the
// original code asserting map[string]interface{} survived: nothing ever gave
// it the set.
func optionsSet(elems ...map[string]interface{}) *schema.Set {
	list := make([]interface{}, 0, len(elems))
	for _, e := range elems {
		list = append(list, e)
	}
	return schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: smarthookoptions.Schema(),
	}), list)
}

// TestInflateOptionsFromSet covers the shape that panicked.
//
// options is a TypeSet, so Terraform produces a *schema.Set. Inflate asserted
// map[string]interface{} on it, which panics -- and no fixture reached that
// code because the fixtures used pre-0.12 syntax and never parsed.
func TestInflateOptionsFromSet(t *testing.T) {
	t.Run("reads the set Terraform produces", func(t *testing.T) {
		out := Inflate(map[string]interface{}{
			"type": "pre-authentication",
			"options": optionsSet(map[string]interface{}{
				"risk_enabled":     true,
				"location_enabled": true,
			}),
		})

		if out.Options == nil {
			t.Fatal("expected options to be read from the set")
		}
		if out.Options.RiskEnabled == nil || !*out.Options.RiskEnabled {
			t.Fatalf("expected risk_enabled true, got %v", out.Options.RiskEnabled)
		}
		if out.Options.LocationEnabled == nil || !*out.Options.LocationEnabled {
			t.Fatalf("expected location_enabled true, got %v", out.Options.LocationEnabled)
		}
	})

	t.Run("still reads a bare map", func(t *testing.T) {
		// Callers that build the input themselves, including the older tests
		// in this package, pass the element directly.
		out := Inflate(map[string]interface{}{
			"type":    "pre-authentication",
			"options": map[string]interface{}{"risk_enabled": true},
		})

		if out.Options == nil || out.Options.RiskEnabled == nil || !*out.Options.RiskEnabled {
			t.Fatalf("expected the map shape to still work, got %v", out.Options)
		}
	})

	t.Run("leaves options unset for an empty set", func(t *testing.T) {
		out := Inflate(map[string]interface{}{
			"type":    "pre-authentication",
			"options": optionsSet(),
		})

		if out.Options != nil {
			t.Fatalf("expected no options, got %v", out.Options)
		}
	})

	t.Run("does not panic on an unexpected shape", func(t *testing.T) {
		// The original failure was a panic, so the guard matters more than the
		// value: a bad shape must be ignored, not crash the provider.
		for _, bad := range []interface{}{42, "options", []string{"x"}} {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("panicked on %#v: %v", bad, r)
					}
				}()

				if out := Inflate(map[string]interface{}{"type": "pre-authentication", "options": bad}); out.Options != nil {
					t.Fatalf("expected no options for %#v", bad)
				}
			}()
		}
	})
}

// TestOptionsSchemaAllowsOneBlock pins MaxItems. The model holds a single
// Options struct, so a second block cannot be represented — without this,
// Inflate would silently keep one and discard the rest.
func TestOptionsSchemaAllowsOneBlock(t *testing.T) {
	options := Schema()["options"]

	if options == nil {
		t.Fatal("expected the smarthook schema to have options")
	}
	if options.MaxItems != 1 {
		t.Fatalf("expected MaxItems 1 so a second block is refused at plan time, got %d", options.MaxItems)
	}
}
