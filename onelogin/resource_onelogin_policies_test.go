package onelogin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"

	policyschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/policy"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

func TestPolicyID(t *testing.T) {
	for _, tc := range []struct {
		name   string
		policy map[string]interface{}
		want   string
		ok     bool
	}{
		// What a decoded JSON response actually carries.
		{name: "float64", policy: map[string]interface{}{"id": float64(42)}, want: "42", ok: true},
		{name: "int", policy: map[string]interface{}{"id": 42}, want: "42", ok: true},
		{name: "string", policy: map[string]interface{}{"id": "42"}, want: "42", ok: true},
		{name: "empty string", policy: map[string]interface{}{"id": ""}, ok: false},
		{name: "missing", policy: map[string]interface{}{}, ok: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id, ok := policyID(tc.policy)
			assert.Equal(t, tc.ok, ok)
			if tc.ok {
				assert.Equal(t, tc.want, id)
			}
		})
	}
}

func TestDecodePolicy(t *testing.T) {
	response := func(status int, body string) *http.Response {
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body))}
	}

	t.Run("returns the policy, which the endpoint sends unwrapped", func(t *testing.T) {
		policy, err := decodePolicy(response(http.StatusOK, `{"id":42,"name":"Engineering","kind":"user"}`))

		assert.NoError(t, err)
		assert.Equal(t, "Engineering", policy["name"])
	})

	t.Run("reports an error status", func(t *testing.T) {
		_, err := decodePolicy(response(http.StatusNotFound, `{}`))

		assert.Error(t, err)
		// policyRead relies on this wording to tell a deleted policy from a
		// failure, so it is worth pinning down here rather than in an
		// acceptance test that needs a tenant.
		assert.True(t, utils.IsNotFoundError(err), "a 404 should be recognisable as one, got %q", err)
	})

	t.Run("reports what the API said about a rejected write", func(t *testing.T) {
		// The whole reason this does not use the SDK's CheckHTTPResponse: the
		// status alone does not say which of seventy-odd arguments was wrong.
		_, err := decodePolicy(response(http.StatusUnprocessableEntity,
			`{"name":"UnprocessableEntityError","message":"Password expiration days is not applicable to app policies","statusCode":422}`))

		assert.ErrorContains(t, err, "status: 422")
		assert.ErrorContains(t, err, "Password expiration days is not applicable to app policies")
	})

	t.Run("falls back to the status when the body explains nothing", func(t *testing.T) {
		_, err := decodePolicy(response(http.StatusInternalServerError, `<html>gateway</html>`))

		assert.ErrorContains(t, err, "status: 500")
	})

	t.Run("rejects a response that is not a policy object", func(t *testing.T) {
		_, err := decodePolicy(response(http.StatusOK, `[{"id":42}]`))

		assert.Error(t, err)
	})
}

// TestPolicyCheckKind drives the check the way a plan does, through Diff with
// the raw configuration on the state, rather than calling it directly. The
// wiring is the part worth testing: the check reads d.GetRawConfig(), and if
// that were not populated it would pass everything silently.
func TestPolicyCheckKind(t *testing.T) {
	diffFor := func(t *testing.T, config map[string]interface{}) error {
		t.Helper()

		raw := map[string]cty.Value{}
		for name, value := range config {
			switch v := value.(type) {
			case string:
				raw[name] = cty.StringVal(v)
			case bool:
				raw[name] = cty.BoolVal(v)
			case int:
				raw[name] = cty.NumberIntVal(int64(v))
			default:
				t.Fatalf("unsupported config value %T", value)
			}
		}

		_, err := Policies().Diff(
			context.Background(),
			&terraform.InstanceState{RawConfig: cty.ObjectVal(raw)},
			terraform.NewResourceConfigRaw(config),
			nil,
		)
		return err
	}

	t.Run("rejects an app-only field on a user policy", func(t *testing.T) {
		err := diffFor(t, map[string]interface{}{
			"name":        "Engineering",
			"kind":        "user",
			"force_authn": true,
		})

		assert.ErrorContains(t, err, "field is not applicable to user policies: force_authn")
	})

	t.Run("rejects user-only fields on an app policy, naming all of them", func(t *testing.T) {
		err := diffFor(t, map[string]interface{}{
			"name":                     "Zoom",
			"kind":                     "app",
			"password_expiration_days": 90,
			"social_sign_in":           true,
		})

		assert.ErrorContains(t, err, "fields are not applicable to app policies: password_expiration_days, social_sign_in")
	})

	t.Run("accepts a field that belongs to the kind", func(t *testing.T) {
		assert.NoError(t, diffFor(t, map[string]interface{}{
			"name":                     "Engineering",
			"kind":                     "user",
			"password_expiration_days": 90,
		}))
	})

	t.Run("accepts a shared field on either kind", func(t *testing.T) {
		for _, kind := range []string{"user", "app"} {
			assert.NoError(t, diffFor(t, map[string]interface{}{
				"name":                "Shared",
				"kind":                kind,
				"ip_addr_restriction": "10.0.0.0/8",
				"otp_auth_enabled":    true,
			}), "kind %s", kind)
		}
	})
}

// TestPolicyNoPerpetualDiff checks the property that only shows up in an
// acceptance test otherwise -- "the plan was not empty" after an apply -- and
// which every Optional+Computed attribute risks: the API answers with far more
// fields than the configuration set, and each one is written to state.
func TestPolicyNoPerpetualDiff(t *testing.T) {
	config := map[string]interface{}{
		"name":                     "Engineering",
		"kind":                     "user",
		"password_expiration_days": 90,
		"enable_password_change":   false,
	}

	r := Policies()
	d := r.Data(nil)
	d.SetId("42")

	// The state a read leaves behind: the configured fields, plus every default
	// the API filled in and returned.
	if err := policyschema.Flatten(d, map[string]interface{}{
		"id":                        float64(42),
		"name":                      "Engineering",
		"kind":                      "user",
		"is_default":                false,
		"password_expiration_days":  float64(90),
		"enable_password_change":    false,
		"minimum_password_length":   float64(8),
		"passwords_remembered":      float64(3),
		"session_timeout_minutes":   float64(480),
		"ignore_xff":                true,
		"new_portal_setting":        "allowed",
		"authentication_factor_ids": []interface{}{float64(3), float64(5)},
		// A policy with no reset factors and no terms contract. The API sends
		// both keys regardless, the second as null, and both have to reach
		// state as empty or the plan is never empty.
		"reset_password_authentication_factor_ids": []interface{}{},
		"terms_and_conditions":                     nil,
	}); err != nil {
		t.Fatal(err)
	}

	diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}

	if diff != nil && len(diff.Attributes) > 0 {
		for attribute, change := range diff.Attributes {
			t.Errorf("unexpected diff on %s: %q -> %q", attribute, change.Old, change.New)
		}
	}
}

func TestAccOneLoginPolicy_crud(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOneLoginPolicyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_policies.test", "name", "Test User Policy"),
					resource.TestCheckResourceAttr("onelogin_policies.test", "kind", "user"),
					resource.TestCheckResourceAttr("onelogin_policies.test", "password_expiration_days", "90"),
					resource.TestCheckResourceAttr("onelogin_policies.test", "minimum_password_length", "12"),
					// Nobody configured this; it is the API's own default,
					// reported back because the attribute is Computed.
					resource.TestCheckResourceAttrSet("onelogin_policies.test", "is_default"),
				),
			},
			{
				Config: testAccCheckOneLoginPolicyConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_policies.test", "name", "Updated Test User Policy"),
					resource.TestCheckResourceAttr("onelogin_policies.test", "password_expiration_days", "30"),
					resource.TestCheckResourceAttr("onelogin_policies.test", "enable_password_change", "false"),
				),
			},
			{
				ResourceName:      "onelogin_policies.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOneLoginPolicy_appKind(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginAppPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_policies.app", "kind", "app"),
					resource.TestCheckResourceAttr("onelogin_policies.app", "force_authn", "true"),
				),
			},
			{
				// The check that saves an apply: a user-only field on an app
				// policy is refused at plan time rather than by a 422.
				Config:      testAccCheckOneLoginAppPolicyConfigWrongKind,
				ExpectError: regexp.MustCompile(`not applicable to app policies`),
			},
		},
	})
}

func testAccCheckOneLoginPolicyDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*onelogin.OneloginSDK)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "onelogin_policies" {
			continue
		}
		path := fmt.Sprintf("%s/%s", policiesPath, rs.Primary.ID)
		resp, err := client.Client.Get(&path, nil)
		if err != nil {
			continue
		}
		if _, err := decodePolicy(resp); err == nil {
			return fmt.Errorf("policy %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccCheckOneLoginPolicyConfig = `
resource "onelogin_policies" "test" {
  name                     = "Test User Policy"
  kind                     = "user"
  password_expiration_days = 90
  minimum_password_length  = 12
  passwords_remembered     = 3
}
`

const testAccCheckOneLoginPolicyConfigUpdated = `
resource "onelogin_policies" "test" {
  name                     = "Updated Test User Policy"
  kind                     = "user"
  password_expiration_days = 30
  minimum_password_length  = 12
  passwords_remembered     = 3
  enable_password_change   = false
}
`

const testAccCheckOneLoginAppPolicyConfig = `
resource "onelogin_policies" "app" {
  name        = "Test App Policy"
  kind        = "app"
  force_authn = true
}
`

const testAccCheckOneLoginAppPolicyConfigWrongKind = `
resource "onelogin_policies" "app" {
  name                     = "Test App Policy"
  kind                     = "app"
  force_authn              = true
  password_expiration_days = 90
}
`
