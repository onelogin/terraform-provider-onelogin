package onelogin

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// brandState builds the state a read leaves behind for an app carrying a brand.
func brandState(t *testing.T, r *schema.Resource, brandID int) *terraform.InstanceState {
	t.Helper()

	d := r.Data(nil)
	d.SetId("1504547")
	for key, value := range map[string]interface{}{
		"name":         "my OIDC APP",
		"connector_id": 38568,
		"brand_id":     brandID,
	} {
		if err := d.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
	return d.State()
}

// TestAppBrandIDLeftAloneWhenUnset is the bug this file exists for.
//
// brand_id was Optional but not Computed, while being filled in by every read.
// An app branded in the OneLogin UI therefore had state holding a brand that no
// configuration mentioned, and Terraform read the absence as a request for 0:
// every plan proposed brand_id -> 0.
//
// What happened next depended on the resource, and neither outcome was good.
// onelogin_apps was the only one that passed brand_id to Inflate, so it sent a
// literal 0 and the apply failed with 422 "The associated AccountBrand with ID
// 0 could not be found". onelogin_oidc_apps and onelogin_saml_apps never sent
// the field at all, so the apply did nothing, the next read put the brand back,
// and the same diff returned for ever.
func TestAppBrandIDLeftAloneWhenUnset(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, brandState(t, r, 44530), map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
			})

			if diff == nil {
				return
			}
			if change, ok := diff.Attributes["brand_id"]; ok {
				t.Fatalf("expected brand_id to be left alone, got %q -> %q", change.Old, change.New)
			}
		})
	}
}

// TestAppBrandIDConfigurable covers the wire contract: settable, and still
// computed so that a brand assigned outside Terraform survives.
func TestAppBrandIDConfigurable(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			attribute, ok := newResource().CoreConfigSchema().Attributes["brand_id"]
			if !ok {
				t.Fatal("brand_id is missing from the core schema")
			}
			if !attribute.Optional {
				t.Error("brand_id is not Optional, so a configuration cannot set it")
			}
			if !attribute.Computed {
				t.Error("brand_id is not Computed, so a brand assigned outside Terraform would be cleared")
			}
		})
	}
}

// TestAppBrandIDSettles guards against the reverse: a brand a configuration
// asked for and got must not be proposed again on the next plan.
func TestAppBrandIDSettles(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, brandState(t, r, 44530), map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"brand_id":     44530,
			})

			if diff == nil {
				return
			}
			if change, ok := diff.Attributes["brand_id"]; ok {
				t.Fatalf("expected brand_id to settle, got %q -> %q", change.Old, change.New)
			}
		})
	}
}

// TestAppBrandIDBody covers what actually reaches the API, on every resource.
//
// The second and third cases are the wiring half of the bug: brand_id was
// Optional on all three resources but only ever passed to Inflate by
// onelogin_apps, and then only on update. Setting it on an OIDC or SAML app was
// accepted by the plan and silently discarded.
func TestAppBrandIDBody(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()

			t.Run("create sends the configured brand", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"brand_id":     44530,
				}, "brand_id", addAppAssignmentForCreate)

				if !strings.Contains(got, `"brand_id":44530`) {
					t.Fatalf("expected brand_id to be sent, got %s", got)
				}
			})

			// Same as policy_id: a create has no brand to take off, so an
			// explicit 0 and an absent argument produce the same app.
			t.Run("create omits an explicit zero", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"brand_id":     0,
				}, "brand_id", addAppAssignmentForCreate)

				if strings.Contains(got, `"brand_id"`) {
					t.Fatalf("expected brand_id to be omitted on create, got %s", got)
				}
			})

			t.Run("update sends a changed brand", func(t *testing.T) {
				got := appBody(t, r, brandState(t, r, 44530), map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"brand_id":     44531,
				}, "brand_id", addAppAssignmentForUpdate)

				if !strings.Contains(got, `"brand_id":44531`) {
					t.Fatalf("expected the new brand to be sent, got %s", got)
				}
			})

			// 0 is refused by the API -- 422 "The associated AccountBrand with
			// ID 0 could not be found" -- and null is what unassigns.
			t.Run("update sends null to unassign", func(t *testing.T) {
				got := appBody(t, r, brandState(t, r, 44530), map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"brand_id":     0,
				}, "brand_id", addAppAssignmentForUpdate)

				if !strings.Contains(got, `"brand_id":null`) {
					t.Fatalf("expected brand_id to be sent as null, got %s", got)
				}
				if strings.Contains(got, `"brand_id":0`) {
					t.Fatalf("expected 0 not to be sent literally, got %s", got)
				}
			})

			t.Run("update omits an unchanged brand", func(t *testing.T) {
				got := appBody(t, r, brandState(t, r, 44530), map[string]interface{}{
					"name":         "a renamed app",
					"connector_id": 38568,
				}, "brand_id", addAppAssignmentForUpdate)

				if strings.Contains(got, `"brand_id"`) {
					t.Fatalf("expected an unrelated update to leave brand_id out, got %s", got)
				}
			})
		})
	}
}

// TestAppBrandIDNegativeRejected mirrors the policy_id rule: 0 unassigns, a
// positive assigns, a negative names nothing.
func TestAppBrandIDNegativeRejected(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			diags := schema.InternalMap(newResource().Schema).Validate(terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"brand_id":     -1,
			}))

			if !diags.HasError() {
				t.Fatal("expected a negative brand_id to be rejected")
			}
		})
	}
}

// TestAppResourcesWireBothAssignments is the test the original bug needed and
// did not have.
//
// brand_id was declared on all three app resources and passed to
// appschema.Inflate by exactly one of them, on update only. Nothing noticed,
// because every test that exercised the field called the helper directly rather
// than going through the resource. This walks the maps the resources actually
// build, so a resource that stops passing a field fails here.
func TestAppResourcesWireBothAssignments(t *testing.T) {
	builders := map[string]func(*schema.ResourceData) map[string]interface{}{
		"onelogin_apps/create":      basicAppCreateMap,
		"onelogin_apps/update":      basicAppUpdateMap,
		"onelogin_oidc_apps/create": oidcAppCreateMap,
		"onelogin_oidc_apps/update": oidcAppUpdateMap,
		"onelogin_saml_apps/create": samlAppCreateMap,
		"onelogin_saml_apps/update": samlAppUpdateMap,
	}
	resources := map[string]func() *schema.Resource{
		"onelogin_apps":      Apps,
		"onelogin_oidc_apps": OIDCApps,
		"onelogin_saml_apps": SAMLApps,
	}

	for name, build := range builders {
		t.Run(name, func(t *testing.T) {
			r := resources[strings.SplitN(name, "/", 2)[0]]()

			// A create map is built from configuration alone; an update map
			// reads HasChange, so it needs a diff behind it. Going through the
			// resource's own diff gives both.
			state := appState(t, r, 955633)
			if err := stateSet(state, "brand_id", "44530"); err != nil {
				t.Fatal(err)
			}
			d, err := schema.InternalMap(r.Schema).Data(state, appDiff(t, r, state, map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"policy_id":    955634,
				"brand_id":     44531,
			}))
			if err != nil {
				t.Fatalf("data: %v", err)
			}

			inflateMap := build(d)
			for _, key := range []string{"policy_id", "brand_id"} {
				if _, ok := inflateMap[key]; !ok {
					t.Errorf("%s never reaches Inflate: %v", key, keysOf(inflateMap))
				}
			}
		})
	}
}

func stateSet(state *terraform.InstanceState, key, value string) error {
	state.Attributes[key] = value
	return nil
}

func keysOf(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
