package policy

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	s := Schema()

	t.Run("name and kind are required, and kind cannot be changed in place", func(t *testing.T) {
		assert.True(t, s["name"].Required)
		assert.True(t, s["kind"].Required)
		assert.True(t, s["kind"].ForceNew, "the API refuses to change kind, so a change has to replace the policy")
	})

	t.Run("is_default is read-only", func(t *testing.T) {
		assert.True(t, s["is_default"].Computed)
		assert.False(t, s["is_default"].Optional, "which policy is the default is set through another endpoint")
	})

	t.Run("every plain field is optional and computed", func(t *testing.T) {
		for _, f := range fields {
			attribute, ok := s[f.Name]
			if assert.True(t, ok, "%s is missing from the schema", f.Name) {
				assert.True(t, attribute.Optional, "%s should be optional", f.Name)
				assert.True(t, attribute.Computed, "%s should be computed, or the API's own default shows as a diff", f.Name)
				assert.Equal(t, f.Type, attribute.Type, "%s has the wrong type", f.Name)
			}
		}
	})

	t.Run("every field has a description", func(t *testing.T) {
		for name, attribute := range s {
			assert.NotEmpty(t, attribute.Description, "%s has no description", name)
		}
	})
}

func TestConfiguredKeys(t *testing.T) {
	t.Run("returns the attributes that are not null", func(t *testing.T) {
		keys := ConfiguredKeys(cty.ObjectVal(map[string]cty.Value{
			"name":                     cty.StringVal("Engineering"),
			"kind":                     cty.StringVal("user"),
			"otp_auth_enabled":         cty.False,
			"password_expiration_days": cty.NullVal(cty.Number),
		}))

		assert.Equal(t, map[string]bool{"name": true, "kind": true, "otp_auth_enabled": true}, keys)
	})

	t.Run("returns nil when the raw configuration is unavailable", func(t *testing.T) {
		// nil rather than an empty map: the caller has to be able to tell
		// "nothing was configured" from "the configuration could not be read".
		assert.Nil(t, ConfiguredKeys(cty.NullVal(cty.EmptyObject)))
		assert.Nil(t, ConfiguredKeys(cty.UnknownVal(cty.EmptyObject)))
		assert.Nil(t, ConfiguredKeys(cty.StringVal("not an object")))
	})
}

func TestFieldsNotApplicableTo(t *testing.T) {
	configured := func(names ...string) map[string]bool {
		keys := map[string]bool{}
		for _, name := range names {
			keys[name] = true
		}
		return keys
	}

	t.Run("app-only fields do not belong on a user policy", func(t *testing.T) {
		assert.Equal(t,
			[]string{"app_otp_offset", "force_authn", "gdt_required"},
			FieldsNotApplicableTo("user", configured("name", "force_authn", "gdt_required", "app_otp_offset")),
		)
	})

	t.Run("user-only fields do not belong on an app policy", func(t *testing.T) {
		assert.Equal(t,
			[]string{"password_expiration_days", "secure_area_otp_timeout_minutes"},
			FieldsNotApplicableTo("app", configured("password_expiration_days", "secure_area_otp_timeout_minutes")),
		)
	})

	t.Run("reset password factors and terms belong to user policies only", func(t *testing.T) {
		names := configured("reset_password_authentication_factor_ids", "terms_and_conditions")
		assert.Equal(t,
			[]string{"reset_password_authentication_factor_ids", "terms_and_conditions"},
			FieldsNotApplicableTo("app", names),
		)
		assert.Empty(t, FieldsNotApplicableTo("user", names))
	})

	t.Run("shared fields belong on both kinds", func(t *testing.T) {
		shared := configured(
			"otp_auth_enabled", "ip_addr_restriction", "ignore_xff", "browser_cert_required",
			"third_party_device_trust", "enable_smart_access", "smart_access_risk_threshold",
			"authentication_factor_ids",
		)
		assert.Empty(t, FieldsNotApplicableTo("user", shared))
		assert.Empty(t, FieldsNotApplicableTo("app", shared))
	})

	t.Run("an unrecognised kind has nothing inapplicable", func(t *testing.T) {
		// The kind validator reports that; reporting every field as well would
		// bury it.
		assert.Empty(t, FieldsNotApplicableTo("", configured("force_authn", "password_expiration_days")))
	})
}

func TestRequestBody(t *testing.T) {
	data := func(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
		t.Helper()
		return schema.TestResourceDataRaw(t, Schema(), raw)
	}

	t.Run("sends only the configured fields", func(t *testing.T) {
		d := data(t, map[string]interface{}{
			"name":                     "Engineering",
			"kind":                     "user",
			"password_expiration_days": 90,
		})

		body := RequestBody(d, map[string]bool{"name": true, "kind": true, "password_expiration_days": true})

		assert.Equal(t, map[string]interface{}{
			"name":                     "Engineering",
			"password_expiration_days": 90,
		}, body)
		assert.NotContains(t, body, "kind", "kind is added by the create path alone")
	})

	t.Run("sends a configured false rather than treating it as unset", func(t *testing.T) {
		// The reason ConfiguredKeys exists: several of these fields default to
		// true, so dropping an explicit false would silently leave them on.
		d := data(t, map[string]interface{}{
			"name":                   "Engineering",
			"enable_password_change": false,
		})

		body := RequestBody(d, map[string]bool{"name": true, "enable_password_change": true})

		assert.Equal(t, false, body["enable_password_change"])
	})

	t.Run("leaves out a field that holds a value but was not configured", func(t *testing.T) {
		// What a refreshed Optional+Computed attribute looks like: a value in
		// state that the configuration never mentioned.
		d := data(t, map[string]interface{}{
			"name":                     "Engineering",
			"password_expiration_days": 90,
		})

		body := RequestBody(d, map[string]bool{"name": true})

		assert.NotContains(t, body, "password_expiration_days")
	})

	t.Run("falls back to non-zero values when the configuration is unavailable", func(t *testing.T) {
		d := data(t, map[string]interface{}{
			"name":                     "Engineering",
			"password_expiration_days": 90,
			"enable_password_change":   false,
		})

		body := RequestBody(d, nil)

		assert.Equal(t, 90, body["password_expiration_days"])
		assert.NotContains(t, body, "enable_password_change", "the fallback cannot tell a configured false from an absent one")
	})

	t.Run("sends factor IDs as a sorted list", func(t *testing.T) {
		// A set has no order of its own, so it is sorted rather than left to
		// the set's internal hash order, which would make the request body
		// differ between runs that configured the same thing.
		d := data(t, map[string]interface{}{
			"name":                      "Engineering",
			"authentication_factor_ids": []interface{}{7, 3, 5},
		})

		body := RequestBody(d, map[string]bool{"name": true, "authentication_factor_ids": true})

		assert.Equal(t, []int{3, 5, 7}, body["authentication_factor_ids"])
	})

	t.Run("sends an empty factor list, which clears the factors", func(t *testing.T) {
		d := data(t, map[string]interface{}{"name": "Engineering"})

		body := RequestBody(d, map[string]bool{"name": true, "authentication_factor_ids": true})

		assert.Equal(t, []int{}, body["authentication_factor_ids"])
	})

	t.Run("sends terms and conditions as an object", func(t *testing.T) {
		d := data(t, map[string]interface{}{
			"name": "Engineering",
			"terms_and_conditions": []interface{}{
				map[string]interface{}{"enabled": true, "content": "Be careful out there."},
			},
		})

		body := RequestBody(d, map[string]bool{"name": true, "terms_and_conditions": true})

		assert.Equal(t, map[string]interface{}{
			"enabled": true,
			"content": "Be careful out there.",
		}, body["terms_and_conditions"])
	})
}

func TestFlatten(t *testing.T) {
	flatten := func(t *testing.T, policy map[string]interface{}) *schema.ResourceData {
		t.Helper()
		d := schema.TestResourceDataRaw(t, Schema(), map[string]interface{}{})
		if err := Flatten(d, policy); err != nil {
			t.Fatal(err)
		}
		return d
	}

	t.Run("writes the fields the response carries", func(t *testing.T) {
		d := flatten(t, map[string]interface{}{
			"id":                       float64(42),
			"name":                     "Engineering",
			"kind":                     "user",
			"is_default":               true,
			"password_expiration_days": float64(90),
			"enable_password_change":   false,
			"ip_addr_restriction":      "10.0.0.0/8",
		})

		assert.Equal(t, "Engineering", d.Get("name"))
		assert.Equal(t, "user", d.Get("kind"))
		assert.Equal(t, true, d.Get("is_default"))
		// JSON numbers decode as float64; an int attribute has to be given an int.
		assert.Equal(t, 90, d.Get("password_expiration_days"))
		assert.Equal(t, false, d.Get("enable_password_change"))
		assert.Equal(t, "10.0.0.0/8", d.Get("ip_addr_restriction"))
	})

	t.Run("leaves a field the response omits alone", func(t *testing.T) {
		// The API leaves out every field belonging to the other kind, so an
		// app policy carries no password_expiration_days. Writing a zero there
		// would put a value in state for something the policy does not have.
		d := schema.TestResourceDataRaw(t, Schema(), map[string]interface{}{
			"password_expiration_days": 90,
		})

		if err := Flatten(d, map[string]interface{}{"name": "Zoom", "kind": "app"}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 90, d.Get("password_expiration_days"))
	})

	t.Run("writes both factor lists", func(t *testing.T) {
		d := flatten(t, map[string]interface{}{
			"name":                      "Engineering",
			"kind":                      "user",
			"authentication_factor_ids": []interface{}{float64(3), float64(5)},
			"reset_password_authentication_factor_ids": []interface{}{float64(7)},
		})

		assert.ElementsMatch(t, []interface{}{3, 5}, d.Get("authentication_factor_ids").(*schema.Set).List())
		assert.ElementsMatch(t, []interface{}{7}, d.Get("reset_password_authentication_factor_ids").(*schema.Set).List())
	})

	t.Run("writes terms and conditions", func(t *testing.T) {
		d := flatten(t, map[string]interface{}{
			"name": "Engineering",
			"kind": "user",
			"terms_and_conditions": map[string]interface{}{
				"enabled": true,
				"content": "Be careful out there.",
			},
		})

		terms := d.Get("terms_and_conditions").([]interface{})
		if assert.Len(t, terms, 1) {
			assert.Equal(t, map[string]interface{}{
				"enabled": true,
				"content": "Be careful out there.",
			}, terms[0])
		}
	})

	t.Run("leaves terms and conditions alone when the policy has none", func(t *testing.T) {
		d := flatten(t, map[string]interface{}{
			"name":                 "Engineering",
			"kind":                 "user",
			"terms_and_conditions": nil,
		})

		assert.Empty(t, d.Get("terms_and_conditions"))
	})
}
