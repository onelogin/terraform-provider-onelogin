package onelogin

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	appparametersschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
)

// These run a plan the way Terraform does, through the provider's gRPC
// PlanResourceChange, against state built by each resource's own read path.
// gh issue 274 lived in the SDK's set comparison and was invisible below this
// level -- every function involved did what it said -- so this is the level
// the tests sit at. TestAccSAMLApp_parameterChanges covers the same ground
// against the API; these need no credentials, so they run in CI.

// apiParam is one parameter as GET /api/2/apps/{id} returns it, with unset
// strings and lists as null.
func apiParam(id float64, set map[string]interface{}) map[string]interface{} {
	param := map[string]interface{}{
		"id":                         id,
		"label":                      nil,
		"values":                     nil,
		"default_values":             nil,
		"user_attribute_mappings":    nil,
		"user_attribute_macros":      nil,
		"attributes_transformations": nil,
		"provisioned_entitlements":   false,
		"skip_if_blank":              false,
	}
	for k, v := range set {
		param[k] = v
	}
	return param
}

// planSAMLApp plans a SAML app whose state is what samlAppRead writes for the
// given parameters, against a configuration with the given parameters blocks.
func planSAMLApp(t *testing.T, fromAPI map[string]interface{}, configured ...map[string]cty.Value) (prior, planned cty.Value) {
	t.Helper()
	return plannedValues(t, planApp(t, "onelogin_saml_apps",
		map[string]interface{}{
			"sso":           map[string]interface{}{"issuer": "https://example.onelogin.com/saml/metadata/1234"},
			"configuration": map[string]interface{}{"signature_algorithm": "SHA-256"},
			"parameters":    appparametersschema.FlattenV4(fromAPI),
		},
		map[string]cty.Value{
			"configuration": cty.MapVal(map[string]cty.Value{"signature_algorithm": cty.StringVal("SHA-256")}),
		},
		configured...,
	))
}

// planGenericApp is planSAMLApp for onelogin_apps, whose read builds the
// parameters through its own conversion and Flatten rather than FlattenV4.
func planGenericApp(t *testing.T, fromAPI map[string]interface{}, configured ...map[string]cty.Value) (prior, planned cty.Value) {
	t.Helper()
	return plannedValues(t, planApp(t, "onelogin_apps",
		map[string]interface{}{"parameters": appparametersschema.Flatten(appParametersFromAPI(fromAPI))},
		nil,
		configured...,
	))
}

// appPlan is a plan request and what came back.
type appPlan struct {
	ty    cty.Type
	prior cty.Value
	resp  *tfprotov5.PlanResourceChangeResponse
}

// planApp plans the named app resource. Its state holds the attributes every
// app read writes plus the given ones, the parameters among them as that
// resource's read flattens them -- with no prior state, RetainManaged keeps
// them all. Its configuration sets name and connector_id, the given
// attributes, and the given parameters blocks.
func planApp(t *testing.T, typeName string, state map[string]interface{}, config map[string]cty.Value, configured ...map[string]cty.Value) appPlan {
	t.Helper()
	provider := Provider()
	res := provider.ResourcesMap[typeName]
	block := res.CoreConfigSchema()
	ty := block.ImpliedType()

	d := res.Data(nil)
	d.SetId("1234")
	attrs := map[string]interface{}{
		"name":                 "App",
		"connector_id":         110016,
		"visible":              true,
		"allow_assumed_signin": false,
		"role_ids":             []interface{}{},
		"provisioning":         map[string]interface{}{"enabled": false},
	}
	for k, v := range state {
		attrs[k] = v
	}
	for k, v := range attrs {
		if err := d.Set(k, v); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	prior, err := d.State().AttrsAsObjectValue(ty)
	if err != nil {
		t.Fatal(err)
	}

	paramTy := ty.AttributeType("parameters").ElementType()
	elements := make([]cty.Value, 0, len(configured))
	for _, p := range configured {
		elements = append(elements, objectWithNulls(paramTy, p))
	}
	configAttrs := map[string]cty.Value{
		"name":         cty.StringVal("App"),
		"connector_id": cty.NumberIntVal(110016),
		"parameters":   setOf(paramTy, elements),
	}
	for k, v := range config {
		configAttrs[k] = v
	}
	configVal := objectWithNulls(ty, configAttrs)

	// What Terraform proposes: the configuration, with anything it leaves
	// unset that is Computed carried over from prior state. A parameters
	// block is matched to one in state only when it sets nothing but its key
	// -- core compares blocks with their Optional+Computed attributes nulled
	// on the state side only -- so a block that sets anything is proposed as
	// written.
	proposedAttrs := map[string]cty.Value{}
	for name := range ty.AttributeTypes() {
		value := configVal.GetAttr(name)
		if attr, ok := block.Attributes[name]; ok && attr.Computed && value.IsNull() {
			value = prior.GetAttr(name)
		}
		proposedAttrs[name] = value
	}
	proposedParams := make([]cty.Value, 0, len(elements))
	for _, element := range elements {
		proposedParams = append(proposedParams, matchKeyOnlyBlock(prior.GetAttr("parameters"), element))
	}
	proposedAttrs["parameters"] = setOf(paramTy, proposedParams)
	proposed := cty.ObjectVal(proposedAttrs)

	resp, err := schema.NewGRPCProviderServer(provider).PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         typeName,
		PriorState:       dynamicValue(t, prior, ty),
		ProposedNewState: dynamicValue(t, proposed, ty),
		Config:           dynamicValue(t, configVal, ty),
	})
	if err != nil {
		t.Fatal(err)
	}
	return appPlan{ty: ty, prior: prior, resp: resp}
}

// plannedValues returns the prior and planned values of a plan that is
// expected to succeed.
func plannedValues(t *testing.T, plan appPlan) (prior, planned cty.Value) {
	t.Helper()
	for _, diag := range plan.resp.Diagnostics {
		t.Fatalf("plan: %s: %s", diag.Summary, diag.Detail)
	}
	planned, err := msgpack.Unmarshal(plan.resp.PlannedState.MsgPack, plan.ty)
	if err != nil {
		t.Fatal(err)
	}
	return plan.prior, planned
}

func objectWithNulls(ty cty.Type, set map[string]cty.Value) cty.Value {
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

// setOf is cty.SetVal, which panics on no elements, or the empty set.
func setOf(elementTy cty.Type, elements []cty.Value) cty.Value {
	if len(elements) == 0 {
		return cty.SetValEmpty(elementTy)
	}
	return cty.SetVal(elements)
}

func matchKeyOnlyBlock(prior, element cty.Value) cty.Value {
	if prior.IsNull() {
		return element
	}
	for it := element.ElementIterator(); it.Next(); {
		name, value := it.Element()
		if name.AsString() != "param_key_name" && !value.IsNull() {
			return element
		}
	}
	for it := prior.ElementIterator(); it.Next(); {
		_, candidate := it.Element()
		if candidate.GetAttr("param_key_name").RawEquals(element.GetAttr("param_key_name")) {
			return candidate
		}
	}
	return element
}

func dynamicValue(t *testing.T, value cty.Value, ty cty.Type) *tfprotov5.DynamicValue {
	t.Helper()
	b, err := msgpack.Marshal(value, ty)
	if err != nil {
		t.Fatal(err)
	}
	return &tfprotov5.DynamicValue{MsgPack: b}
}

// noChanges is the test Terraform applies to decide a plan is a no-op:
// equality of values, not of their representation. A number that has been
// through msgpack is not RawEquals to the one it started as.
func noChanges(planned, prior cty.Value) bool {
	eq := planned.Equals(prior)
	return eq.IsKnown() && eq.True()
}

// plannedParam returns one attribute of the planned parameter with the given
// key name.
func plannedParam(t *testing.T, planned cty.Value, key, attr string) cty.Value {
	t.Helper()
	for it := planned.GetAttr("parameters").ElementIterator(); it.Next(); {
		_, param := it.Element()
		if param.GetAttr("param_key_name").AsString() == key {
			return param.GetAttr(attr)
		}
	}
	t.Fatalf("expected %s in the planned parameters, got %#v", key, planned.GetAttr("parameters"))
	return cty.NilVal
}

func emailFromAPI(mapping string) map[string]interface{} {
	return map[string]interface{}{
		"email": apiParam(101, map[string]interface{}{
			"label":                     "Email",
			"user_attribute_mappings":   mapping,
			"include_in_saml_assertion": true,
		}),
	}
}

func emailBlock(mapping string) map[string]cty.Value {
	return map[string]cty.Value{
		"param_key_name":            cty.StringVal("email"),
		"label":                     cty.StringVal("Email"),
		"user_attribute_mappings":   cty.StringVal(mapping),
		"include_in_saml_assertion": cty.True,
	}
}

func TestSAMLAppPlan_parameters(t *testing.T) {
	t.Run("plans a changed mapping", func(t *testing.T) {
		// gh issue 274: the mapping is the only change, and this planned
		// nothing.
		prior, planned := planSAMLApp(t, emailFromAPI("email"), emailBlock("username"))

		if noChanges(planned, prior) {
			t.Fatal("expected a plan, got no changes")
		}
		if got := plannedParam(t, planned, "email", "user_attribute_mappings"); !got.RawEquals(cty.StringVal("username")) {
			t.Fatalf("expected the new mapping to be planned, got %#v", got)
		}
		// Carried from state, rather than proposed as unknown.
		if got := plannedParam(t, planned, "email", "label"); !got.RawEquals(cty.StringVal("Email")) {
			t.Fatalf("expected the label to stay as it is, got %#v", got)
		}
	})

	t.Run("plans a mapping changed in OneLogin back", func(t *testing.T) {
		// The same blind spot from the other side: drift in a configured
		// attribute.
		_, planned := planSAMLApp(t, emailFromAPI("firstname"), emailBlock("username"))

		if got := plannedParam(t, planned, "email", "user_attribute_mappings"); !got.RawEquals(cty.StringVal("username")) {
			t.Fatalf("expected the configured mapping to be planned, got %#v", got)
		}
	})

	t.Run("plans no changes when nothing configured has changed", func(t *testing.T) {
		// Two ways a configuration differs from state without meaning
		// to: groups leaves its label to the API, and spells out falses
		// the way the documentation's examples do, one of which the API
		// holds as true -- the model never sends a false, so no apply could
		// change it.
		fromAPI := emailFromAPI("email")
		fromAPI["groups"] = apiParam(102, map[string]interface{}{
			"label":                     "Groups",
			"default_values":            "one",
			"include_in_saml_assertion": true,
		})

		prior, planned := planSAMLApp(t, fromAPI, emailBlock("email"), map[string]cty.Value{
			"param_key_name":            cty.StringVal("groups"),
			"default_values":            cty.ListVal([]cty.Value{cty.StringVal("one")}),
			"include_in_saml_assertion": cty.False,
			"skip_if_blank":             cty.False,
		})

		if !noChanges(planned, prior) {
			t.Fatalf("expected no changes\nprior:   %#v\nplanned: %#v", prior.GetAttr("parameters"), planned.GetAttr("parameters"))
		}
	})

	t.Run("rejects two blocks for one key", func(t *testing.T) {
		// OneLogin keeps one parameter per key, so planning both would be
		// a diff no apply can settle.
		plan := planApp(t, "onelogin_saml_apps",
			map[string]interface{}{"parameters": appparametersschema.FlattenV4(emailFromAPI("email"))},
			nil,
			emailBlock("email"), emailBlock("username"),
		)

		if len(plan.resp.Diagnostics) == 0 || !strings.Contains(plan.resp.Diagnostics[0].Detail+plan.resp.Diagnostics[0].Summary, "configured more than once") {
			t.Fatalf("expected the duplicate key to be rejected, got %v", plan.resp.Diagnostics)
		}
	})
}

func TestAppPlan_parameters(t *testing.T) {
	t.Run("plans a changed mapping", func(t *testing.T) {
		prior, planned := planGenericApp(t, emailFromAPI("email"), emailBlock("username"))

		if noChanges(planned, prior) {
			t.Fatal("expected a plan, got no changes")
		}
		if got := plannedParam(t, planned, "email", "user_attribute_mappings"); !got.RawEquals(cty.StringVal("username")) {
			t.Fatalf("expected the new mapping to be planned, got %#v", got)
		}
	})

	t.Run("plans no changes for safe_entitlements_enabled the API holds", func(t *testing.T) {
		// appRead never copied safe_entitlements_enabled, so state read a
		// configured true as false and every plan proposed it again.
		fromAPI := map[string]interface{}{
			"roles": apiParam(103, map[string]interface{}{"label": "Roles", "safe_entitlements_enabled": true}),
		}

		prior, planned := planGenericApp(t, fromAPI, map[string]cty.Value{
			"param_key_name":            cty.StringVal("roles"),
			"safe_entitlements_enabled": cty.True,
		})

		if !noChanges(planned, prior) {
			t.Fatalf("expected no changes\nprior:   %#v\nplanned: %#v", prior.GetAttr("parameters"), planned.GetAttr("parameters"))
		}
	})
}
