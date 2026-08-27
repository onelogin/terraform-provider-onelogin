package onelogin

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// policy_id and brand_id both name a record an app is assigned to, and behave
// identically: omitted leaves the assignment alone, a positive ID assigns, and
// 0 unassigns. The three app resources -- basic, OIDC and SAML -- share
// appschema.Schema and so share both fields, and these are the two rules for
// when to send one, rather than the same dozen lines written out per resource.
//
// Both rules are about the key. appschema.Inflate reads its presence, not the
// truth of its value, because 0 is meaningful: it is how a configuration says
// "none", and Inflate turns it into the JSON null the app endpoint takes as an
// unassignment.

// addAppAssignmentForCreate records what a new app should be created with, and
// says nothing when the configuration named nothing.
//
// GetOk reads an explicit 0 as unset, and that is wanted here even though 0 is
// meaningful everywhere else in this file. On an update a 0 is an instruction --
// take the assignment off -- and has to be sent. On a create there is nothing to
// take off: an app created without the field comes back with the column null,
// which is the same app a null would have produced. Omitting it and sending the
// null are indistinguishable in the result, so the simpler of the two is fine.
//
// GetOkExists would tell the two apart, but it is deprecated in
// terraform-plugin-sdk v2 -- "usage is discouraged due to undefined behaviors"
// -- and GetRawConfig would be the supported way if a real difference ever
// appeared. The create-body tests pin the behaviour either way.
func addAppAssignmentForCreate(d *schema.ResourceData, inflateMap map[string]interface{}, key string) {
	if id, ok := d.GetOk(key); ok {
		inflateMap[key] = id
	}
}

// addAppAssignmentForUpdate records what an update should apply, and says
// nothing when it did not change.
//
// Sent only when it actually changed, and then whatever it now is -- including
// the 0 that unassigns. The app endpoint takes a PUT but merges it, so a field
// left out is left alone, and that is what keeps an update touching only the
// name from disturbing either assignment. It matters more here than it does for
// a group: both fields are Computed as well as Optional, so state routinely
// holds a value no configuration asked for, and d.Get would echo it back on
// every unrelated apply.
//
// An assignment changed outside Terraform does not reach this either. State
// falls to 0 on the next read, but so does the absent attribute it is compared
// against, so HasChange stays false and nothing is sent.
func addAppAssignmentForUpdate(d *schema.ResourceData, inflateMap map[string]interface{}, key string) {
	if d.HasChange(key) {
		inflateMap[key] = d.Get(key)
	}
}
