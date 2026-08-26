package onelogin

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// The three app resources -- basic, OIDC and SAML -- share appschema.Schema and
// so share policy_id. They also share the two rules for when to send it, which
// live here rather than being repeated six times.
//
// Both rules are about the key: appschema.Inflate reads its presence, not the
// truth of its value, because 0 is meaningful. It is how a configuration says
// "no policy", and Inflate turns it into the JSON null the app endpoint takes
// as an unassignment.

// addAppPolicyIDForCreate records the policy a new app should be created with,
// and says nothing when the configuration named none.
//
// GetOk reads an explicit policy_id = 0 as unset, and that is wanted here even
// though 0 is meaningful everywhere else in this file. On an update a 0 is an
// instruction -- take the policy off -- and has to be sent. On a create there
// is nothing to take off: an app created without the field comes back with
// policy_id null, which is the same app a null would have produced. Omitting it
// and sending the null are indistinguishable in the result, so the simpler of
// the two is fine.
//
// GetOkExists would tell the two apart, but it is deprecated in
// terraform-plugin-sdk v2 -- "usage is discouraged due to undefined behaviors"
// -- and GetRawConfig would be the supported way if a real difference ever
// appeared. TestAppPolicyIDCreateBody pins the behaviour either way.
func addAppPolicyIDForCreate(d *schema.ResourceData, inflateMap map[string]interface{}) {
	if policyID, ok := d.GetOk("policy_id"); ok {
		inflateMap["policy_id"] = policyID
	}
}

// addAppPolicyIDForUpdate records the policy an update should apply, and says
// nothing when it did not change.
//
// Sent only when it actually changed, and then whatever it now is -- including
// the 0 that unassigns. The app endpoint takes a PUT but merges it, so a field
// left out is left alone, and that is what keeps an update touching only the
// name from disturbing the policy. It matters more here than it does for a
// group: policy_id is Computed as well as Optional on an app, so state
// routinely holds a value no configuration asked for, and d.Get would echo it
// back on every unrelated apply.
//
// A policy removed outside Terraform does not reach this either. State falls to
// 0 on the next read, but so does the absent attribute it is compared against,
// so HasChange stays false and nothing is sent.
func addAppPolicyIDForUpdate(d *schema.ResourceData, inflateMap map[string]interface{}) {
	if d.HasChange("policy_id") {
		inflateMap["policy_id"] = d.Get("policy_id")
	}
}
