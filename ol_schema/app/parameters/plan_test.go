package appparametersschema

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// configParam is one parameters block as Terraform hands it over: every
// attribute present, and null where the configuration says nothing.
func configParam(set map[string]cty.Value) cty.Value {
	ty := (&schema.Resource{Schema: Schema()}).CoreConfigSchema().ImpliedType()
	attrs := map[string]cty.Value{}
	for name, attrTy := range ty.AttributeTypes() {
		if v, ok := set[name]; ok {
			attrs[name] = v
		} else {
			attrs[name] = cty.NullVal(attrTy)
		}
	}
	return cty.ObjectVal(attrs)
}

func configParams(params ...map[string]cty.Value) cty.Value {
	elements := make([]cty.Value, 0, len(params))
	for _, p := range params {
		elements = append(elements, configParam(p))
	}
	return cty.SetVal(elements)
}

func planned(t *testing.T, plan []interface{}, key string) map[string]interface{} {
	t.Helper()
	for _, raw := range plan {
		param := raw.(map[string]interface{})
		if param["param_key_name"] == key {
			return param
		}
	}
	t.Fatalf("expected %s to be planned, got %v", key, plan)
	return nil
}

// The parameter samlAppRead writes for a configured email mapping: the label
// and flag the configuration set, and an ID the API assigned.
func emailInState(mapping string) map[string]interface{} {
	return stateParam(map[string]interface{}{
		"param_key_name":            "email",
		"param_id":                  369278,
		"label":                     "Email",
		"user_attribute_mappings":   mapping,
		"include_in_saml_assertion": true,
	})
}

func TestPlan(t *testing.T) {
	t.Run("plans a changed mapping over the parameter in state", func(t *testing.T) {
		// gh issue 274: the mapping is the only change.
		prior := schema.NewSet(Hash, []interface{}{emailInState("email")})
		plan := Plan(prior, configParams(map[string]cty.Value{
			"param_key_name":            cty.StringVal("email"),
			"label":                     cty.StringVal("Email"),
			"user_attribute_mappings":   cty.StringVal("username"),
			"include_in_saml_assertion": cty.True,
		}))

		got := planned(t, plan, "email")
		if want := emailInState("username"); !reflect.DeepEqual(got, want) {
			t.Fatalf("expected the parameter in state with the new mapping\n got: %v\nwant: %v", got, want)
		}
		if Hash(got) == Hash(emailInState("email")) {
			t.Fatal("expected the planned parameter to hash differently from state, or the SDK sees no change")
		}
	})

	t.Run("plans exactly what state holds when nothing configured has changed", func(t *testing.T) {
		// The configuration leaves label and the flags to the API. Read as
		// zero values they would look like changes; carried from state they
		// are not.
		prior := schema.NewSet(Hash, []interface{}{emailInState("email")})
		plan := Plan(prior, configParams(map[string]cty.Value{
			"param_key_name":          cty.StringVal("email"),
			"user_attribute_mappings": cty.StringVal("email"),
		}))

		got := planned(t, plan, "email")
		if !reflect.DeepEqual(got, emailInState("email")) {
			t.Fatalf("expected state unchanged, got %v", got)
		}
		if Hash(got) != Hash(emailInState("email")) {
			t.Fatal("expected the same hash as state, or every plan proposes replacing the set")
		}
	})

	t.Run("does not plan a value an apply cannot send", func(t *testing.T) {
		// Inflate or the model's omitempty drops "", false and an empty
		// values list. The documentation's examples spell out = false, so
		// planning these would be a diff no apply can settle.
		inState := stateParam(map[string]interface{}{
			"param_key_name":            "groups",
			"label":                     "Groups",
			"values":                    []interface{}{"memberOf"},
			"include_in_saml_assertion": true,
			"skip_if_blank":             true,
			"provisioned_entitlements":  true,
		})
		plan := Plan(schema.NewSet(Hash, []interface{}{inState}), configParams(map[string]cty.Value{
			"param_key_name":            cty.StringVal("groups"),
			"label":                     cty.StringVal(""),
			"values":                    cty.ListValEmpty(cty.String),
			"include_in_saml_assertion": cty.False,
			"skip_if_blank":             cty.False,
			"provisioned_entitlements":  cty.False,
		}))

		if got := planned(t, plan, "groups"); !reflect.DeepEqual(got, inState) {
			t.Fatalf("expected state unchanged, got %v", got)
		}
	})

	t.Run("plans an empty default_values, which clears it", func(t *testing.T) {
		// Unlike values, default_values is sent even when empty, as an
		// explicit null, and the API clears it.
		inState := stateParam(map[string]interface{}{"param_key_name": "groups", "default_values": []interface{}{"one"}})
		plan := Plan(schema.NewSet(Hash, []interface{}{inState}), configParams(map[string]cty.Value{
			"param_key_name": cty.StringVal("groups"),
			"default_values": cty.ListValEmpty(cty.String),
		}))

		if got := planned(t, plan, "groups")["default_values"]; !reflect.DeepEqual(got, []interface{}{}) {
			t.Fatalf("expected an empty list to be planned, got %v", got)
		}
	})

	t.Run("plans safe_entitlements_enabled false, which is sent", func(t *testing.T) {
		// A pointer on the model, so unlike the other flags a false reaches
		// the API.
		inState := stateParam(map[string]interface{}{"param_key_name": "roles", "safe_entitlements_enabled": true})
		plan := Plan(schema.NewSet(Hash, []interface{}{inState}), configParams(map[string]cty.Value{
			"param_key_name":            cty.StringVal("roles"),
			"safe_entitlements_enabled": cty.False,
		}))

		if got := planned(t, plan, "roles")["safe_entitlements_enabled"]; got != false {
			t.Fatalf("expected false to be planned, got %v", got)
		}
	})

	t.Run("plans a changed list", func(t *testing.T) {
		inState := stateParam(map[string]interface{}{"param_key_name": "groups", "default_values": []interface{}{"one"}})
		plan := Plan(schema.NewSet(Hash, []interface{}{inState}), configParams(map[string]cty.Value{
			"param_key_name": cty.StringVal("groups"),
			"default_values": cty.ListVal([]cty.Value{cty.StringVal("one"), cty.StringVal("two")}),
		}))

		got := planned(t, plan, "groups")["default_values"]
		if want := []interface{}{"one", "two"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})

	t.Run("plans a new parameter from the configuration alone", func(t *testing.T) {
		plan := Plan(schema.NewSet(Hash, []interface{}{emailInState("email")}), configParams(
			map[string]cty.Value{
				"param_key_name":          cty.StringVal("email"),
				"user_attribute_mappings": cty.StringVal("email"),
			},
			map[string]cty.Value{
				"param_key_name":          cty.StringVal("department"),
				"user_attribute_mappings": cty.StringVal("department"),
				"skip_if_blank":           cty.False,
			},
		))

		want := map[string]interface{}{"param_key_name": "department", "user_attribute_mappings": "department"}
		if got := planned(t, plan, "department"); !reflect.DeepEqual(got, want) {
			t.Fatalf("expected only what the configuration sets\n got: %v\nwant: %v", got, want)
		}
	})

	t.Run("leaves out a parameter the configuration no longer has", func(t *testing.T) {
		prior := schema.NewSet(Hash, []interface{}{
			emailInState("email"),
			stateParam(map[string]interface{}{"param_key_name": "firstname"}),
		})
		plan := Plan(prior, configParams(map[string]cty.Value{
			"param_key_name":          cty.StringVal("email"),
			"user_attribute_mappings": cty.StringVal("email"),
		}))

		if len(plan) != 1 {
			t.Fatalf("expected only the configured parameter, got %v", plan)
		}
	})
}

// TestSentValuesMatchesInflate holds sentValues to the request body. Plan only
// proposes what an apply sends, so for every attribute, set to its zero value
// and to something else, the two must agree on whether it goes out: planning a
// value Inflate drops is a perpetual diff, and dropping one Inflate sends is
// gh issue 274 again. An attribute added to Schema with a type sentValues does
// not handle fails here rather than going silently unplanned.
func TestSentValuesMatchesInflate(t *testing.T) {
	for attr, s := range Schema() {
		if attr == "param_key_name" || attr == "param_id" {
			continue
		}

		type sample struct {
			config cty.Value
			raw    interface{}
		}
		var samples []sample
		switch s.Type {
		case schema.TypeString:
			samples = []sample{{cty.StringVal(""), ""}, {cty.StringVal("x"), "x"}}
		case schema.TypeBool:
			samples = []sample{{cty.False, false}, {cty.True, true}}
		case schema.TypeList:
			samples = []sample{
				{cty.ListValEmpty(cty.String), []interface{}{}},
				{cty.ListVal([]cty.Value{cty.StringVal("x")}), []interface{}{"x"}},
			}
		default:
			t.Fatalf("%s: sentValues does not handle %s", attr, s.Type)
		}

		for _, sample := range samples {
			_, planned := sentValues(configParam(map[string]cty.Value{
				"param_key_name": cty.StringVal("k"),
				attr:             sample.config,
			}))[attr]

			b, err := json.Marshal(Inflate(map[string]interface{}{"param_key_name": "k", attr: sample.raw}))
			if err != nil {
				t.Fatal(err)
			}
			var body map[string]interface{}
			if err := json.Unmarshal(b, &body); err != nil {
				t.Fatal(err)
			}
			_, sent := body[attr]

			if planned != sent {
				t.Errorf("%s = %#v: planned %v, but sent %v (body %s)", attr, sample.raw, planned, sent, b)
			}
		}
	}
}

// TestUniqueKeyNames covers two blocks for one key. OneLogin keeps one
// parameter per key, so planning both would be a diff no apply can settle.
func TestUniqueKeyNames(t *testing.T) {
	block := func(name, mapping string) map[string]cty.Value {
		return map[string]cty.Value{
			"param_key_name":          cty.StringVal(name),
			"user_attribute_mappings": cty.StringVal(mapping),
		}
	}

	err := uniqueKeyNames(configParams(block("email", "email"), block("email", "username")))
	if err == nil || !strings.Contains(err.Error(), `"email"`) {
		t.Fatalf("expected an error naming the key, got %v", err)
	}

	if err := uniqueKeyNames(configParams(block("email", "email"), block("firstname", "firstname"))); err != nil {
		t.Fatalf("expected distinct keys to pass, got %v", err)
	}

	unknown := configParams(block("email", "email"), map[string]cty.Value{"param_key_name": cty.UnknownVal(cty.String)})
	if err := uniqueKeyNames(unknown); err != nil {
		t.Fatalf("expected a key not yet known to be skipped, got %v", err)
	}
}
