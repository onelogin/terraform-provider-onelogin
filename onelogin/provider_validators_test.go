package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestValidatorsDoNotPanicOnUnexpectedTypes walks every ValidateFunc the
// provider actually wires up and hands each one values it should never see.
//
// The guard belongs in the validators, but the risk is that a new one is added
// later with a bare assertion and nobody notices until a provider crashes in
// somebody's pipeline. Walking the live schema means this covers validators
// that do not exist yet.
func TestValidatorsDoNotPanicOnUnexpectedTypes(t *testing.T) {
	// Types a validator should never receive, plus a string that is simply not
	// an allowed value. The string exercises the ordinary rejection path: a
	// guard that returned early on every input would satisfy the rest of this
	// list while breaking validation entirely.
	unexpected := []interface{}{
		nil, 42, true, 3.14, []string{"x"}, map[string]string{}, struct{}{},
		"definitely-not-an-allowed-value",
	}

	var walk func(t *testing.T, path string, s map[string]*schema.Schema)
	walk = func(t *testing.T, path string, s map[string]*schema.Schema) {
		for name, attr := range s {
			key := path + name

			if attr.ValidateFunc != nil {
				for _, val := range unexpected {
					func() {
						defer func() {
							if r := recover(); r != nil {
								t.Errorf("%s panicked on %#v: %v — a validator must report a bad value, not take the provider down", key, val, r)
							}
						}()
						attr.ValidateFunc(val, key)
					}()
				}
			}

			// Nested resources carry their own validators.
			if res, ok := attr.Elem.(*schema.Resource); ok {
				walk(t, key+".", res.Schema)
			}
			if elem, ok := attr.Elem.(*schema.Schema); ok && elem.ValidateFunc != nil {
				for _, val := range unexpected {
					func() {
						defer func() {
							if r := recover(); r != nil {
								t.Errorf("%s element validator panicked on %#v: %v", key, val, r)
							}
						}()
						elem.ValidateFunc(val, key)
					}()
				}
			}
		}
	}

	p := Provider()
	if len(p.ResourcesMap) == 0 {
		t.Fatal("expected the provider to expose resources")
	}

	// The provider configuration block too, not only resources and data
	// sources. Its attributes take a ValidateFunc like any others, and a panic
	// there is worse than most: it happens before anything can be configured.
	walk(t, "provider.", p.Schema)

	for name, res := range p.ResourcesMap {
		walk(t, name+".", res.Schema)
	}
	for name, ds := range p.DataSourcesMap {
		walk(t, name+".", ds.Schema)
	}
}
