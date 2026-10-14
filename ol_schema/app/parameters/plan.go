package appparametersschema

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Hash identifies a parameter in the set by everything a practitioner can write
// to it.
//
// It used to be the key name alone. The SDK compares a set by its hash codes
// before it compares a single attribute, and when old and new produce the same
// codes it reports no change without looking inside any element. A key-name
// hash therefore hid every attribute but the key from the plan: a changed
// user_attribute_mappings planned nothing and was never sent (gh issue 274),
// and a mapping changed in OneLogin was never put back.
//
// Hashing the content on its own brings back what the key-name hash was
// solving: a configuration leaves most attributes unset, so its elements hash
// differently from the ones a read wrote, and every plan proposes replacing the
// whole set. CustomizeDiff is what deals with that now, by planning each
// configured parameter as the one already in state with the configuration
// written over it. The codes then differ only where a configured attribute
// differs from state in a way an apply can write.
var Hash = schema.HashResource(&schema.Resource{Schema: Schema()})

// CustomizeDiff plans the parameters of an existing app.
//
// Without it, the SDK compares the configuration with state attribute by
// attribute. Every attribute but the key name is Optional+Computed, and an
// attribute the configuration leaves out reads as its zero value, so a
// parameter that only sets a mapping looks nothing like the same parameter in
// state, label and flags filled in by the API. Plan builds the element the
// practitioner means instead.
func CustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	config := d.GetRawConfig()
	if config.IsNull() || !config.IsKnown() || !config.Type().HasAttribute("parameters") {
		return nil
	}
	configured := config.GetAttr("parameters")
	if configured.IsNull() || !configured.IsKnown() {
		return nil
	}

	if err := uniqueKeyNames(configured); err != nil {
		return err
	}

	// A create has no state to carry forward, and the SDK's own diff of the
	// configuration is already what it should send. No parameters blocks
	// keeps meaning "leave them alone", as it always has. Anything not yet
	// known is left to the SDK's own diff: Terraform plans again during the
	// apply, once it is known, and this runs then.
	if d.Id() == "" || !configured.IsWhollyKnown() || configured.LengthInt() == 0 {
		return nil
	}

	prior, _ := d.GetChange("parameters")
	return d.SetNew("parameters", Plan(prior, configured))
}

// uniqueKeyNames rejects two parameters blocks with the same param_key_name.
//
// OneLogin keeps one parameter per key, so only one of them could ever be
// written. The key-name hash used to fold them into one set element and drop
// the other without a word; with every attribute in the hash, both would plan
// on every run and neither apply could settle it.
func uniqueKeyNames(configured cty.Value) error {
	seen := map[string]bool{}
	for it := configured.ElementIterator(); it.Next(); {
		_, element := it.Element()
		if !element.IsKnown() || element.IsNull() {
			continue
		}
		name := element.GetAttr("param_key_name")
		if !name.IsKnown() || name.IsNull() {
			continue
		}
		if seen[name.AsString()] {
			return fmt.Errorf("parameters: param_key_name %q is configured more than once; an app has one parameter per key", name.AsString())
		}
		seen[name.AsString()] = true
	}
	return nil
}

// Plan returns the parameters a plan should propose: each configured one as it
// stands in prior state, with every attribute the configuration sets written
// over it -- but only where an apply would actually send that value.
//
// What an apply can send is narrower than what a configuration can say. An
// empty string, a false and an empty values list never reach the API: Inflate
// or the model's omitempty drops each of them. Planning them would propose a
// change no apply can make -- a perpetual diff -- and the documentation's own
// examples spell out include_in_saml_assertion = false.
//
// Two exceptions do reach the API, and so are planned. safe_entitlements_enabled
// is a pointer on the model, so its false is sent. An empty default_values goes
// out as an explicit null, which the API takes as clearing it.
//
// A parameter that is not yet in state is planned from the configuration
// alone, and one in state but not in the configuration is left out, the same
// as RetainManaged does on a read.
func Plan(prior interface{}, configured cty.Value) []interface{} {
	inState := map[string]map[string]interface{}{}
	if set, ok := prior.(*schema.Set); ok {
		for _, raw := range set.List() {
			if param, ok := raw.(map[string]interface{}); ok {
				if name, ok := param["param_key_name"].(string); ok {
					inState[name] = param
				}
			}
		}
	}

	out := make([]interface{}, 0, configured.LengthInt())
	for it := configured.ElementIterator(); it.Next(); {
		_, element := it.Element()
		name := element.GetAttr("param_key_name")
		if name.IsNull() || !name.IsKnown() {
			continue
		}

		planned := map[string]interface{}{}
		for attr, value := range inState[name.AsString()] {
			planned[attr] = value
		}
		planned["param_key_name"] = name.AsString()
		for attr, value := range sentValues(element) {
			planned[attr] = value
		}
		out = append(out, planned)
	}
	return out
}

// sentWhenEmpty names the attributes whose zero value still reaches the API,
// the two exceptions Plan describes. TestSentValuesMatchesInflate holds this
// to what Inflate and the model actually put in the request body.
var sentWhenEmpty = map[string]bool{
	"safe_entitlements_enabled": true,
	"default_values":            true,
}

// sentValues returns the attributes of one configured parameter that an apply
// would send: set, and not a value the model drops.
func sentValues(element cty.Value) map[string]interface{} {
	out := map[string]interface{}{}
	for attr, s := range Schema() {
		if attr == "param_key_name" || attr == "param_id" || !element.Type().HasAttribute(attr) {
			continue
		}
		value := element.GetAttr(attr)
		if value.IsNull() || !value.IsKnown() {
			continue
		}

		switch s.Type {
		case schema.TypeString:
			if str := value.AsString(); str != "" {
				out[attr] = str
			}
		case schema.TypeBool:
			if value.True() || sentWhenEmpty[attr] {
				out[attr] = value.True()
			}
		case schema.TypeList:
			list := make([]interface{}, 0, value.LengthInt())
			for it := value.ElementIterator(); it.Next(); {
				_, item := it.Element()
				if item.IsKnown() && !item.IsNull() {
					list = append(list, item.AsString())
				}
			}
			if len(list) > 0 || sentWhenEmpty[attr] {
				out[attr] = list
			}
		}
	}
	return out
}
