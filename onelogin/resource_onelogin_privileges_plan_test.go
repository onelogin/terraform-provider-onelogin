package onelogin

import (
	"context"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	privilegeschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/privilege"
)

// A perpetual diff is normally only visible in an acceptance test, as "After
// applying this test step, the plan was not empty" -- which needs credentials
// and a tenant, and in this resource's case only showed up on the runs where
// the API happened to return a different order.
//
// These tests run the same check offline: build the state privilegeRead would
// write from a given response, then ask the SDK for a plan against the
// configuration. What they cannot do is predict the order the API returns, so
// each one names it explicitly.

// planAfterRead returns the diff a plan produces, given the privilege already
// in ResourceData -- prior state on a refresh, the configuration at the end of
// a create or update -- and the statements a read returned.
func planAfterRead(t *testing.T, prior map[string]interface{}, returned []map[string]interface{}, config map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()
	r := Privileges()

	d := r.Data(nil)
	d.SetId("role_id")
	if err := d.Set("name", "role_name"); err != nil {
		t.Fatal(err)
	}
	if prior != nil {
		if err := d.Set("privilege", []map[string]interface{}{prior}); err != nil {
			t.Fatal(err)
		}
	}

	// The write privilegeRead performs, with the HTTP call already made.
	if err := d.Set("privilege", []map[string]interface{}{{
		"version":   privilegeschema.DefaultVersion,
		"statement": privilegeschema.OrderStatementsLikeState(d.Get("privilege"), returned),
	}}); err != nil {
		t.Fatal(err)
	}

	diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":      "role_name",
		"privilege": []interface{}{config},
	}), nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	return diff
}

func privilegeDocument(action ...string) map[string]interface{} {
	return map[string]interface{}{
		"version":   privilegeschema.DefaultVersion,
		"statement": []interface{}{privilegeStatement(action...)},
	}
}

// privilegeResponse is the statement list privilegeRead builds from the JSON.
func privilegeResponse(action ...string) []map[string]interface{} {
	return []map[string]interface{}{privilegeStatement(action...)}
}

func privilegeStatement(action ...string) map[string]interface{} {
	acts := make([]interface{}, 0, len(action))
	for _, a := range action {
		acts = append(acts, a)
	}
	return map[string]interface{}{
		"effect": "Allow",
		"action": acts,
		"scope":  []interface{}{"users/*"},
	}
}

func changedAttributes(diff *terraform.InstanceDiff) []string {
	if diff == nil {
		return nil
	}
	out := make([]string, 0, len(diff.Attributes))
	for k := range diff.Attributes {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// TestPrivilegePlanIsEmptyWhenOnlyTheAPIOrderDiffers covers gh-254. The
// practitioner changed nothing; the API returned the statement's actions in a
// different order from the one they were written in. action is a TypeList
// inside a TypeSet, so that changed the privilege element's hash, and the plan
// proposed removing the whole block and adding it back -- every time, including
// straight after an apply.
func TestPrivilegePlanIsEmptyWhenOnlyTheAPIOrderDiffers(t *testing.T) {
	config := privilegeDocument("users:List", "users:Get")
	response := privilegeResponse("users:Get", "users:List")

	// Create, update and refresh all reach this with the configuration's order
	// in ResourceData: the first two because SDK v2 layers the diff over state
	// and d.Get returns the planned value, refresh because the previous apply
	// already wrote state in that order.
	diff := planAfterRead(t, config, response, config)
	if diff != nil && !diff.Empty() {
		t.Fatalf("plan was not empty for an unchanged privilege: %v", changedAttributes(diff))
	}
}

// TestPrivilegeUpgradeFromApiOrderedState is the path every existing user is
// on: state was written by a version without the value ordering, so it holds
// whatever order the API gave. The first plan still proposes a change -- the
// fix cannot retroactively reorder what is already on disk -- and one apply
// settles it for good.
func TestPrivilegeUpgradeFromApiOrderedState(t *testing.T) {
	config := privilegeDocument("users:List", "users:Get")
	apiOrder := privilegeDocument("users:Get", "users:List")
	response := privilegeResponse("users:Get", "users:List")

	// Before the apply: state is in the API's order, so it is matched against
	// itself and stays there.
	diff := planAfterRead(t, apiOrder, response, config)
	if diff == nil || diff.Empty() {
		t.Fatal("expected the upgrade to plan one change")
	}

	// The apply's own read has the configuration in ResourceData, so state is
	// rewritten in the configuration's order and every plan after is empty.
	if diff = planAfterRead(t, config, response, config); diff != nil && !diff.Empty() {
		t.Fatalf("did not settle after one apply: %v", changedAttributes(diff))
	}
}

// TestPrivilegePlanAgainstRecordedAPIResponse replays the orderings the API
// actually produced, recorded against the development tenant on 2026-08-10 via
// POST/GET/PUT /api/1/privileges. The action orderings below are verbatim; the
// scope is this file's synthetic one, since every recorded statement used "*"
// and a single-element list cannot demonstrate ordering. Scope reordering is
// covered at unit level in ol_schema/privilege/statement_order_test.go.
//
// Sent, in this order:
//
//	statement 0  users:List, users:Get, users:Update, users:Delete
//	statement 1  apps:List, apps:Get
//
// Read back after the create, on four consecutive reads:
//
//	statement 0  users:List, users:Delete, users:Get, users:Update
//	statement 1  apps:List, apps:Get
//
// Then the identical document was PUT back, and the next read returned the
// statements swapped AND the actions in a third order. The same document
// written twice does not come back the same way, so there is no canonical
// order to mirror -- only the configuration's order to restore.
func TestPrivilegePlanAgainstRecordedAPIResponse(t *testing.T) {
	config := map[string]interface{}{
		"version": privilegeschema.DefaultVersion,
		"statement": []interface{}{
			privilegeStatement("users:List", "users:Get", "users:Update", "users:Delete"),
			privilegeStatement("apps:List", "apps:Get"),
		},
	}

	for _, tc := range []struct {
		name     string
		returned []map[string]interface{}
	}{
		{
			"after create: values reordered",
			[]map[string]interface{}{
				privilegeStatement("users:List", "users:Delete", "users:Get", "users:Update"),
				privilegeStatement("apps:List", "apps:Get"),
			},
		},
		{
			"after update: statements swapped and values reordered again",
			[]map[string]interface{}{
				privilegeStatement("apps:Get", "apps:List"),
				privilegeStatement("users:Update", "users:Get", "users:Delete", "users:List"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diff := planAfterRead(t, config, tc.returned, config)
			if diff != nil && !diff.Empty() {
				t.Fatalf("plan was not empty for an unchanged privilege: %v", changedAttributes(diff))
			}
		})
	}
}

// A real edit must still plan, or the fix above would be hiding changes rather
// than settling them.
func TestPrivilegePlanDetectsRealChanges(t *testing.T) {
	prior := privilegeDocument("users:List", "users:Get")
	response := privilegeResponse("users:Get", "users:List")

	for _, tc := range []struct {
		name   string
		config map[string]interface{}
	}{
		{"an action added", privilegeDocument("users:List", "users:Get", "users:Delete")},
		{"an action removed", privilegeDocument("users:List")},
		{"an action swapped", privilegeDocument("users:List", "users:Update")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diff := planAfterRead(t, prior, response, tc.config)
			if diff == nil || diff.Empty() {
				t.Fatal("a real change produced an empty plan")
			}
		})
	}
}

// On import there is no prior order to respond to, so the API's order reaches
// state and the first plan proposes a change. One apply settles it, which is
// the same bargain the statement ordering already makes.
func TestPrivilegePlanAfterImport(t *testing.T) {
	config := privilegeDocument("users:List", "users:Get")

	diff := planAfterRead(t, nil, privilegeResponse("users:Get", "users:List"), config)
	if diff == nil || diff.Empty() {
		t.Fatal("expected import to plan the one-time change")
	}
}
