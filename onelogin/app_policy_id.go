package onelogin

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// The three app resources -- basic, OIDC and SAML -- share appschema.Schema and
// so share policy_id. They also share the two rules for when to send it, which
// live here rather than being repeated six times.
//
// Both rules end up saying "only a real policy ID is ever sent". OneLogin
// refuses 0 with a 422 and treats a JSON null as "no policy", and a *int tagged
// omitempty cannot produce a null, so 0 is worth nothing to either caller.
// appschema.validAppPolicyID stops a configuration reaching them with one.

// addAppPolicyIDForCreate records the policy a new app should be created with,
// and says nothing when the configuration named none.
//
// Omitting the field is what a create with no policy wants in any case: the app
// is stored with policy_id null, which is exactly the state an app that never
// had a policy should be in.
func addAppPolicyIDForCreate(d *schema.ResourceData, inflateMap map[string]interface{}) {
	if policyID, ok := d.GetOk("policy_id"); ok {
		inflateMap["policy_id"] = policyID
	}
}

// addAppPolicyIDForUpdate records the policy an update should apply, and says
// nothing when it did not change.
//
// The app endpoint takes a PUT but merges it, so a field left out is left
// alone. That is what keeps an update touching only the name from disturbing
// the policy, and it matters more here than it does for a group: policy_id is
// Computed as well as Optional on an app, so state routinely holds a value no
// configuration asked for, and d.Get would echo it back on every unrelated
// apply.
func addAppPolicyIDForUpdate(d *schema.ResourceData, inflateMap map[string]interface{}) {
	if !d.HasChange("policy_id") {
		return
	}
	// A 0 here is a policy going away, which the API cannot be asked to do.
	// validAppPolicyID rejects a configured 0 before a plan gets this far, so
	// what reaches this line is state catching up with an app whose policy was
	// removed elsewhere -- and sending that 0 back would only turn a settled
	// read into a 422.
	if policyID := d.Get("policy_id").(int); policyID > 0 {
		inflateMap["policy_id"] = policyID
	}
}
