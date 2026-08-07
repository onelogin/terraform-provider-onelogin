package onelogin

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	roleschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/role"
)

// TestSkippedMembershipAttrs covers issue #242.
//
// Role membership lives on three paginated sub-endpoints, and a role with
// thousands of users costs a page walk on every refresh whether or not the
// configuration manages them. ignore_changes cannot prevent it: it is applied
// to the diff, long after the read has made the calls.
func TestSkippedMembershipAttrs(t *testing.T) {
	data := func(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
		t.Helper()
		return schema.TestResourceDataRaw(t, Roles().Schema, raw)
	}

	t.Run("skips nothing when unset", func(t *testing.T) {
		// The import case, and the one that must stay safe: an imported
		// resource's state carries no value for this attribute, and it has to
		// mean "fetch everything" rather than "fetch nothing".
		got := skippedMembershipAttrs(data(t, map[string]interface{}{"name": "Engineering"}))

		if len(got) != 0 {
			t.Fatalf("expected nothing to be skipped, got %v", got)
		}
	})

	t.Run("skips the named attribute", func(t *testing.T) {
		got := skippedMembershipAttrs(data(t, map[string]interface{}{
			"name":                    "Engineering",
			"skip_membership_refresh": []interface{}{"users"},
		}))

		if !got["users"] {
			t.Fatalf("expected users to be skipped, got %v", got)
		}
		if got["apps"] || got["admins"] {
			t.Fatalf("expected only users to be skipped, got %v", got)
		}
	})

	t.Run("skips several", func(t *testing.T) {
		got := skippedMembershipAttrs(data(t, map[string]interface{}{
			"name":                    "Engineering",
			"skip_membership_refresh": []interface{}{"users", "admins"},
		}))

		if !got["users"] || !got["admins"] {
			t.Fatalf("expected users and admins to be skipped, got %v", got)
		}
		if got["apps"] {
			t.Fatalf("expected apps to still be fetched, got %v", got)
		}
	})
}

// TestSkipMembershipRefreshSchema pins the attribute's shape. The names it
// accepts have to match the attributes roleRead actually iterates, or a typo
// becomes a silent no-op.
func TestSkipMembershipRefreshSchema(t *testing.T) {
	s := Roles().Schema["skip_membership_refresh"]
	if s == nil {
		t.Fatal("expected the roles schema to have skip_membership_refresh")
	}
	if !s.Optional {
		t.Fatal("expected the attribute to be optional")
	}

	elem, ok := s.Elem.(*schema.Schema)
	if !ok || elem.ValidateFunc == nil {
		t.Fatal("expected the elements to be validated")
	}

	for _, valid := range roleschema.MembershipAttrs {
		if _, errs := elem.ValidateFunc(valid, "skip_membership_refresh"); len(errs) > 0 {
			t.Fatalf("expected %q to be accepted, got %v", valid, errs)
		}
	}

	// A near-miss is the likely mistake, and silently doing nothing would be
	// the worst outcome for an attribute whose whole purpose is to prevent work.
	for _, invalid := range []string{"user", "Users", "everything"} {
		if _, errs := elem.ValidateFunc(invalid, "skip_membership_refresh"); len(errs) == 0 {
			t.Fatalf("expected %q to be rejected", invalid)
		}
	}
}

// TestSkipMembershipRefreshIsNotSentToTheAPI guards against the attribute
// leaking into a role payload, which the API would reject.
func TestSkipMembershipRefreshIsNotSentToTheAPI(t *testing.T) {
	role := roleschema.Inflate(map[string]interface{}{
		"name":                    "Engineering",
		"skip_membership_refresh": []interface{}{"users"},
	})

	b, err := json.Marshal(role)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "skip_membership_refresh") {
		t.Fatalf("provider-only attribute leaked into the API payload: %s", b)
	}
}
