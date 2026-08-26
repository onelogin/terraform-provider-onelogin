package onelogin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	appschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// The three app resources share appschema.Schema, and so share policy_id.
// Everything below is asserted against all three rather than against the OIDC
// resource the issue happened to be reported on.
func appResourcesUnderTest() map[string]func() *schema.Resource {
	return map[string]func() *schema.Resource{
		"onelogin_apps":      Apps,
		"onelogin_oidc_apps": OIDCApps,
		"onelogin_saml_apps": SAMLApps,
	}
}

// TestAppPolicyIDIsConfigurable covers issue #260.
//
// The assertion is on CoreConfigSchema rather than on the Schema map because
// that is the schema Terraform is actually handed. policy_id was Computed and
// not Optional, so core saw an attribute it alone decides and rejected any
// configuration mentioning it with "Can't configure a value for policy_id".
// The rejection happened during validation, before a single provider function
// ran, which is why no amount of CRUD-side work could have made the reported
// configuration apply.
//
// Computed stays on alongside Optional; TestAppPolicyIDLeftAloneWhenUnset
// covers why.
func TestAppPolicyIDIsConfigurable(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			attribute, ok := newResource().CoreConfigSchema().Attributes["policy_id"]
			if !ok {
				t.Fatal("policy_id is missing from the core schema")
			}
			if !attribute.Optional {
				t.Error("policy_id is not Optional, so a configuration cannot set it")
			}
			if !attribute.Computed {
				t.Error("policy_id is not Computed, so a policy assigned outside Terraform would be cleared")
			}
		})
	}
}

// appState builds the state a read leaves behind for an app.
func appState(t *testing.T, r *schema.Resource, policyID int) *terraform.InstanceState {
	t.Helper()

	d := r.Data(nil)
	d.SetId("1234567")
	for key, value := range map[string]interface{}{
		"name":         "my OIDC APP",
		"connector_id": 38568,
		"policy_id":    policyID,
	} {
		if err := d.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
	return d.State()
}

// appDiff runs a configuration against prior state through the resource's own
// diff, which is the same path a plan takes.
func appDiff(t *testing.T, r *schema.Resource, state *terraform.InstanceState, config map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()

	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	return diff
}

// TestAppPolicyIDAssignable is the configuration from the issue: a policy
// created elsewhere in the same run, attached to a new app.
func TestAppPolicyIDAssignable(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, nil, map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"policy_id":    955633,
			})

			change, ok := diff.Attributes["policy_id"]
			if !ok {
				t.Fatalf("expected policy_id to be planned, got %v", diff.Attributes)
			}
			if change.New != "955633" {
				t.Fatalf("expected policy_id to be planned as 955633, got %q", change.New)
			}
		})
	}
}

// TestAppPolicyIDLeftAloneWhenUnset is the reason policy_id keeps Computed now
// that it is Optional.
//
// policy_id has been Computed since the resource shipped, so state already
// holds whatever policy the app was given in the OneLogin UI. Optional on its
// own would read a configuration that never mentions policy_id as asking for
// 0, and the first apply after upgrading the provider would quietly unassign
// the policy.
func TestAppPolicyIDLeftAloneWhenUnset(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, appState(t, r, 955633), map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
			})

			if diff == nil {
				return
			}
			if change, ok := diff.Attributes["policy_id"]; ok {
				t.Fatalf("expected policy_id to be left alone, got %q -> %q", change.Old, change.New)
			}
		})
	}
}

// TestAppPolicyIDSettles guards the other direction: an assignment a
// configuration asked for and got must not be proposed again on the next plan.
func TestAppPolicyIDSettles(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, appState(t, r, 955633), map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"policy_id":    955633,
			})

			if diff == nil {
				return
			}
			if change, ok := diff.Attributes["policy_id"]; ok {
				t.Fatalf("expected policy_id to settle, got %q -> %q", change.Old, change.New)
			}
		})
	}
}

// TestAppPolicyIDClearable covers unassigning.
//
// Because policy_id is Computed as well as Optional, dropping the attribute
// keeps the last value, so unassigning has to be said out loud. An explicit 0
// is still a value in the configuration, so it diffs against state -- and
// Inflate turns it into the null the API wants, which the body tests check.
func TestAppPolicyIDClearable(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, appState(t, r, 955633), map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"policy_id":    0,
			})

			change, ok := diff.Attributes["policy_id"]
			if !ok {
				t.Fatalf("expected clearing policy_id to propose a change, got %v", diff.Attributes)
			}
			if change.Old != "955633" || change.New != "0" {
				t.Fatalf("expected 955633 -> 0, got %q -> %q", change.Old, change.New)
			}
		})
	}
}

// TestAppPolicyIDNegativeRejected covers the one value that still cannot mean
// anything. 0 unassigns and a positive assigns; a negative names no policy at
// all, and saying so during planning beats an apply that comes back 422.
func TestAppPolicyIDNegativeRejected(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			diags := schema.InternalMap(newResource().Schema).Validate(terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":         "my OIDC APP",
				"connector_id": 38568,
				"policy_id":    -1,
			}))

			if !diags.HasError() {
				t.Fatal("expected a negative policy_id to be rejected")
			}
		})
	}
}

// appBody runs a configuration and prior state all the way to the JSON the API
// would receive, through the resource's own diff and the helper each app
// resource calls. The bug only exists once the request is serialised, so that
// is where the assertions are.
func appBody(t *testing.T, r *schema.Resource, state *terraform.InstanceState, config map[string]interface{}, add func(*schema.ResourceData, map[string]interface{})) string {
	t.Helper()

	d, err := schema.InternalMap(r.Schema).Data(state, appDiff(t, r, state, config))
	if err != nil {
		t.Fatalf("data: %v", err)
	}

	inflateMap := map[string]interface{}{
		"name":         d.Get("name"),
		"connector_id": d.Get("connector_id"),
	}
	add(d, inflateMap)

	app, err := appschema.Inflate(inflateMap)
	if err != nil {
		t.Fatalf("inflate: %v", err)
	}
	body, err := json.Marshal(app)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(body)
}

// TestAppPolicyIDCreateBody covers what a create actually sends.
//
// A create that names no policy omits the field rather than sending 0, because
// the API stores what it is sent: an app that never had a policy is not the
// same record as one whose policy was cleared.
func TestAppPolicyIDCreateBody(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()

			t.Run("sends the configured policy", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"policy_id":    955633,
				}, addAppPolicyIDForCreate)

				if !strings.Contains(got, `"policy_id":955633`) {
					t.Fatalf("expected policy_id to be sent, got %s", got)
				}
			})

			t.Run("omits policy_id when none is configured", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
				}, addAppPolicyIDForCreate)

				if strings.Contains(got, `"policy_id"`) {
					t.Fatalf("expected policy_id to be omitted, got %s", got)
				}
			})

			// An explicit 0 on a create is omitted too, and that is not the
			// same compromise the update path makes. A create has no policy to
			// take off: the app comes back with policy_id null either way, so
			// omitting the field and sending the null produce the same app.
			// GetOk collapsing the two is therefore harmless here, which is
			// what lets this path avoid the deprecated GetOkExists.
			t.Run("omits an explicit zero on create", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"policy_id":    0,
				}, addAppPolicyIDForCreate)

				if strings.Contains(got, `"policy_id"`) {
					t.Fatalf("expected policy_id to be omitted on create, got %s", got)
				}
			})
		})
	}
}

// TestAppPolicyIDUpdateBody covers what an update actually sends.
//
// The app endpoint takes a PUT but merges it, so leaving a field out leaves it
// alone. That is what lets an update touching only the name keep its hands off
// a policy assigned outside Terraform -- and it is why the 0 that clears an
// assignment has to be sent explicitly when it is genuinely being cleared.
func TestAppPolicyIDUpdateBody(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()

			t.Run("sends a changed policy", func(t *testing.T) {
				got := appBody(t, r, appState(t, r, 955633), map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"policy_id":    955634,
				}, addAppPolicyIDForUpdate)

				if !strings.Contains(got, `"policy_id":955634`) {
					t.Fatalf("expected the new policy to be sent, got %s", got)
				}
			})

			// The whole point of the SDK bump: a configured 0 leaves here as a
			// JSON null. Sending it as 0 comes back 422, and omitting it would
			// leave the policy in place while the plan claimed otherwise.
			t.Run("sends null to unassign", func(t *testing.T) {
				got := appBody(t, r, appState(t, r, 955633), map[string]interface{}{
					"name":         "my OIDC APP",
					"connector_id": 38568,
					"policy_id":    0,
				}, addAppPolicyIDForUpdate)

				if !strings.Contains(got, `"policy_id":null`) {
					t.Fatalf("expected policy_id to be sent as null, got %s", got)
				}
				if strings.Contains(got, `"policy_id":0`) {
					t.Fatalf("expected 0 not to be sent literally, got %s", got)
				}
			})

			t.Run("omits policy_id when it did not change", func(t *testing.T) {
				got := appBody(t, r, appState(t, r, 955633), map[string]interface{}{
					"name":         "a renamed app",
					"connector_id": 38568,
				}, addAppPolicyIDForUpdate)

				if strings.Contains(got, `"policy_id"`) {
					t.Fatalf("expected an unrelated update to leave policy_id out, got %s", got)
				}
			})
		})
	}
}

// TestAppPolicyIDReadsNullAsZero covers the shape an unassigned app comes back
// in. OneLogin reports policy_id as null rather than 0 or an absent key, and
// the read path hands that straight to d.Set.
//
// It has to land on 0, because 0 is what an app with no policy has to hold for
// a configuration that never mentions policy_id to produce an empty plan.
func TestAppPolicyIDReadsNullAsZero(t *testing.T) {
	r := OIDCApps()

	t.Run("null becomes zero", func(t *testing.T) {
		d := r.Data(nil)
		d.SetId("1504537")
		if err := d.Set("policy_id", 955633); err != nil {
			t.Fatal(err)
		}

		utils.SetResourceFields(d, map[string]interface{}{"policy_id": nil}, []string{"policy_id"})

		if got := d.Get("policy_id").(int); got != 0 {
			t.Fatalf("expected a null policy_id to read as 0, got %d", got)
		}
	})

	// The number arrives as a float64 out of encoding/json, not an int.
	t.Run("an assigned policy survives the decode", func(t *testing.T) {
		d := r.Data(nil)
		d.SetId("1504537")

		utils.SetResourceFields(d, map[string]interface{}{"policy_id": float64(955633)}, []string{"policy_id"})

		if got := d.Get("policy_id").(int); got != 955633 {
			t.Fatalf("expected 955633, got %d", got)
		}
	})
}

// TestAccOIDCApp_policy is the configuration from issue #260 end to end: an app
// policy created in the same run, attached to an OIDC app, swapped for a
// different one, and finally taken back off.
//
// The last step is what onelogin-go-sdk v4.16.0 made possible. It is written as
// policy_id = 0 rather than by dropping the attribute, because policy_id is
// Computed as well as Optional -- see TestAppPolicyIDClearable.
//
// The unit tests above stop at the request body, so this is what covers the
// resources actually calling the helpers, and the API accepting what they send.
func TestAccOIDCApp_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOIDCAppPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"onelogin_oidc_apps.policy_test", "policy_id",
						"onelogin_policies.app_policy", "id",
					),
				),
			},
			{
				// Reassignment, which is the case an update has to send.
				Config: testAccOIDCAppPolicyConfigOther,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"onelogin_oidc_apps.policy_test", "policy_id",
						"onelogin_policies.other_policy", "id",
					),
				),
			},
			{
				// Unassignment. The API takes a JSON null for this and refuses
				// a 0, so what this really exercises is Inflate translating the
				// one into the other.
				Config: testAccOIDCAppPolicyConfigCleared,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_oidc_apps.policy_test", "policy_id", "0"),
				),
			},
		},
	})
}

// Only an app policy resolves; a user policy is refused with the same 422 the
// API gives for an ID that does not exist.
const testAccOIDCAppPolicyPolicies = `
resource "onelogin_policies" "app_policy" {
  name                   = "TF Acc App Policy"
  kind                   = "app"
  force_authn            = true
  app_force_authn_offset = 60
}

resource "onelogin_policies" "other_policy" {
  name        = "TF Acc Other App Policy"
  kind        = "app"
  force_authn = true
}
`

const testAccOIDCAppPolicyApp = `
resource "onelogin_oidc_apps" "policy_test" {
  connector_id = 108419
  name         = "TF Acc Policy OIDC App"
  %s

  configuration = {
    redirect_uri                     = "https://localhost:3000/callback"
    refresh_token_expiration_minutes = 1
    login_url                        = "https://www.test.com"
    oidc_application_type            = 0
    token_endpoint_auth_method       = 1
    access_token_expiration_minutes  = 1
  }
}
`

var (
	testAccOIDCAppPolicyConfig      = testAccOIDCAppPolicyPolicies + fmt.Sprintf(testAccOIDCAppPolicyApp, "policy_id = onelogin_policies.app_policy.id")
	testAccOIDCAppPolicyConfigOther = testAccOIDCAppPolicyPolicies + fmt.Sprintf(testAccOIDCAppPolicyApp, "policy_id = onelogin_policies.other_policy.id")

	// Explicit rather than absent: dropping policy_id would leave the last
	// value in place, because the attribute is Computed as well as Optional.
	testAccOIDCAppPolicyConfigCleared = testAccOIDCAppPolicyPolicies + fmt.Sprintf(testAccOIDCAppPolicyApp, "policy_id = 0")
)
