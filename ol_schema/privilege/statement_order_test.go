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
		// Same statement, actions listed the other way round. Order within a
		// statement's action list is the practitioner's, so it is compared
		// order-insensitively but written back untouched.
		configured := statement("Allow", "apps:List", "users:List")
		returned := statement("Allow", "users:List", "apps:List")

		out := OrderStatementsLikeState(priorPrivilege(configured), []map[string]interface{}{returned})

		if len(out) != 1 {
			t.Fatalf("expected the statement to be matched, got %v", out)
		}
		if actions := out[0]["action"].([]interface{}); actions[0] != "users:List" {
			t.Fatalf("expected the API's own values to be written back, got %v", actions)
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
