package privilegeschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func statement(effect string, actions ...string) map[string]interface{} {
	acts := make([]interface{}, 0, len(actions))
	for _, action := range actions {
		acts = append(acts, action)
	}
	return map[string]interface{}{
		"effect": effect,
		"action": acts,
		"scope":  []interface{}{"*"},
	}
}

// TestOrderLikePrior covers the value reordering on its own, including the
// leftover branch that OrderStatementsLikeState cannot reach: a matched
// statement always has the same set of values, so the branch is defensive.
func TestOrderLikePrior(t *testing.T) {
	for _, tc := range []struct {
		name     string
		prior    []string
		returned []string
		want     []string
	}{
		{"restores order", []string{"b", "a"}, []string{"a", "b"}, []string{"b", "a"}},
		{"already in order", []string{"a", "b"}, []string{"a", "b"}, []string{"a", "b"}},
		{"unknown values keep API order", []string{"b"}, []string{"c", "a", "b"}, []string{"b", "c", "a"}},
		{"prior names values that are gone", []string{"x", "a"}, []string{"a"}, []string{"a"}},
		{"no prior", nil, []string{"b", "a"}, []string{"b", "a"}},
		{"duplicates survive", []string{"a", "b"}, []string{"b", "a", "a"}, []string{"a", "b", "a"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := orderLikePrior(tc.prior, tc.returned)
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("got %v, want %v", got, tc.want)
				}
			}
		})
	}
}

// priorPrivilege builds the value d.Get("privilege") returns: a set of one
// document holding the statements in the order the configuration has them.
func priorPrivilege(statements ...map[string]interface{}) *schema.Set {
	list := make([]interface{}, 0, len(statements))
	for _, s := range statements {
		list = append(list, s)
	}
	return schema.NewSet(
		func(interface{}) int { return 0 },
		[]interface{}{map[string]interface{}{"version": DefaultVersion, "statement": list}},
	)
}

// TestOrderStatementsLikeState covers the intermittent perpetual diff. The API
// returns the statements in an order of its own, and statement is a TypeList,
// so a swapped pair reads as a changed value.
func TestOrderStatementsLikeState(t *testing.T) {
	apps := statement("Allow", "apps:List")
	users := statement("Allow", "users:List")

	t.Run("restores the configuration's order", func(t *testing.T) {
		out := OrderStatementsLikeState(priorPrivilege(apps, users), []map[string]interface{}{users, apps})

		if len(out) != 2 {
			t.Fatalf("expected both statements, got %v", out)
		}
		if statementKey(out[0]) != statementKey(apps) || statementKey(out[1]) != statementKey(users) {
			t.Fatalf("expected the configuration's order, got %v", out)
		}
	})

	t.Run("keeps a statement the configuration does not have", func(t *testing.T) {
		roles := statement("Deny", "roles:List")
		out := OrderStatementsLikeState(priorPrivilege(apps), []map[string]interface{}{users, roles, apps})

		if len(out) != 3 {
			t.Fatalf("expected every statement to survive, got %v", out)
		}
		if statementKey(out[0]) != statementKey(apps) {
			t.Fatalf("expected the matched statement first, got %v", out[0])
		}
	})

	t.Run("matches on contents, not position", func(t *testing.T) {
		// Same statement, actions listed the other way round. It is still the
		// same statement, and the configuration's order is what reaches state.
		configured := statement("Allow", "apps:List", "users:List")
		returned := statement("Allow", "users:List", "apps:List")

		out := OrderStatementsLikeState(priorPrivilege(configured), []map[string]interface{}{returned})

		if len(out) != 1 {
			t.Fatalf("expected the statement to be matched, got %v", out)
		}
		if actions := out[0]["action"].([]interface{}); actions[0] != "apps:List" {
			t.Fatalf("expected the configuration's order, got %v", actions)
		}
	})

	t.Run("restores the order of the values inside a statement", func(t *testing.T) {
		// gh-254: the API reorders the values within a statement too, and
		// action is a TypeList inside a TypeSet, so a swapped pair changes the
		// privilege element's hash and every plan proposes replacing the block.
		configured := statement("Allow", "users:List", "users:Get")
		returned := statement("Allow", "users:Get", "users:List")
		returned["scope"] = []interface{}{"users/2", "users/1"}
		configured["scope"] = []interface{}{"users/1", "users/2"}

		out := OrderStatementsLikeState(priorPrivilege(configured), []map[string]interface{}{returned})

		if got := out[0]["action"].([]interface{}); got[0] != "users:List" || got[1] != "users:Get" {
			t.Fatalf("expected the configuration's action order, got %v", got)
		}
		if got := out[0]["scope"].([]interface{}); got[0] != "users/1" || got[1] != "users/2" {
			t.Fatalf("expected the configuration's scope order, got %v", got)
		}
	})

	t.Run("a different set of actions is a different statement", func(t *testing.T) {
		// statementKey covers the whole action set, so reordering only ever
		// applies to statements that are genuinely the same one. A statement
		// that gained an action does not match, and keeps what the API sent --
		// the reordering can never mask a real change.
		configured := statement("Allow", "users:List")
		returned := statement("Allow", "users:Get", "users:List")

		out := OrderStatementsLikeState(priorPrivilege(configured), []map[string]interface{}{returned})

		if got := out[0]["action"].([]interface{}); got[0] != "users:Get" || len(got) != 2 {
			t.Fatalf("expected the API's statement untouched, got %v", got)
		}
	})

	t.Run("reorders a statement whose values are []string", func(t *testing.T) {
		// FlattenPrivilegeData builds the values as []string rather than
		// []interface{}. statementKey reads both, so orderValuesLikeState must
		// too -- a statement that matches on its contents and is then skipped
		// would be the perpetual diff all over again, with nothing logged.
		configured := statement("Allow", "users:List", "users:Get")
		returned := map[string]interface{}{
			"effect": "Allow",
			"action": []string{"users:Get", "users:List"},
			"scope":  []string{"*"},
		}

		out := OrderStatementsLikeState(priorPrivilege(configured), []map[string]interface{}{returned})

		got, ok := toStrings(out[0]["action"])
		if !ok {
			t.Fatalf("action was not a list of values: %#v", out[0]["action"])
		}
		if got[0] != "users:List" || got[1] != "users:Get" {
			t.Fatalf("expected the configuration's order, got %v", got)
		}
	})

	t.Run("an unreadable field is not the same as an empty one", func(t *testing.T) {
		// action and scope are Required, so neither case should reach us. If
		// one does, keying them alike would let two statements that have
		// nothing in common match on their shared effect and be paired up.
		missing := map[string]interface{}{"effect": "Allow", "scope": []interface{}{"*"}}
		empty := map[string]interface{}{"effect": "Allow", "action": []interface{}{}, "scope": []interface{}{"*"}}
		text := map[string]interface{}{"effect": "Allow", "action": "users:List", "scope": []interface{}{"*"}}

		keys := map[string]string{
			"missing action":  statementKey(missing),
			"empty action":    statementKey(empty),
			"action a string": statementKey(text),
		}
		seen := map[string]string{}
		for name, key := range keys {
			if other, clash := seen[key]; clash {
				t.Fatalf("%s and %s share a statement key", name, other)
			}
			seen[key] = name
		}
	})

	t.Run("does not mutate the response it was given", func(t *testing.T) {
		configured := statement("Allow", "users:List", "users:Get")
		returned := statement("Allow", "users:Get", "users:List")

		OrderStatementsLikeState(priorPrivilege(configured), []map[string]interface{}{returned})

		if got := returned["action"].([]interface{}); got[0] != "users:Get" {
			t.Fatalf("the caller's statement was modified in place: %v", got)
		}
	})

	t.Run("leaves the API order alone on import", func(t *testing.T) {
		for _, prior := range []interface{}{nil, priorPrivilege(), "not a privilege"} {
			out := OrderStatementsLikeState(prior, []map[string]interface{}{users, apps})

			if len(out) != 2 || statementKey(out[0]) != statementKey(users) {
				t.Fatalf("expected the response untouched for prior %#v, got %v", prior, out)
			}
		}
	})
}
