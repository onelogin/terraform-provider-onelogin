package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	userschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user"
)

// TestEmailsFromConfig covers reading the emails filter. A blank entry is
// dropped: an interpolated value that came out empty would otherwise become a
// query with no email set, which matches every user in the tenant.
func TestEmailsFromConfig(t *testing.T) {
	t.Run("reads the emails in the order given", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, userschema.QuerySchema(), map[string]interface{}{
			"emails": []interface{}{"alice@example.com", "bob@example.com"},
		})

		emails := emailsFromConfig(d)
		if len(emails) != 2 || emails[0] != "alice@example.com" || emails[1] != "bob@example.com" {
			t.Fatalf("expected both emails in order, got %v", emails)
		}
	})

	t.Run("drops blank entries", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, userschema.QuerySchema(), map[string]interface{}{
			"emails": []interface{}{"alice@example.com", ""},
		})

		if emails := emailsFromConfig(d); len(emails) != 1 || emails[0] != "alice@example.com" {
			t.Fatalf("expected the blank to be dropped, got %v", emails)
		}
	})

	t.Run("reports none when unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, userschema.QuerySchema(), map[string]interface{}{})

		if emails := emailsFromConfig(d); len(emails) != 0 {
			t.Fatalf("expected no emails, got %v", emails)
		}
	})
}

// TestUsersQueriesForEmails covers the fan-out. Every other filter has to
// survive onto each query, or "these people, in this directory" would silently
// widen to "these people, anywhere".
func TestUsersQueriesForEmails(t *testing.T) {
	directory := "42"
	base := &models.UserQuery{DirectoryID: &directory}

	t.Run("returns the query untouched when there are no emails", func(t *testing.T) {
		queries := usersQueriesForEmails(base, nil)

		if len(queries) != 1 || queries[0] != base {
			t.Fatalf("expected the original query, got %v", queries)
		}
	})

	t.Run("builds one query per email", func(t *testing.T) {
		queries := usersQueriesForEmails(base, []string{"alice@example.com", "bob@example.com"})

		if len(queries) != 2 {
			t.Fatalf("expected 2 queries, got %d", len(queries))
		}
		for i, want := range []string{"alice@example.com", "bob@example.com"} {
			if queries[i].Email == nil || *queries[i].Email != want {
				t.Fatalf("query %d: expected email %q, got %v", i, want, queries[i].Email)
			}
		}
	})

	t.Run("carries the other filters onto every query", func(t *testing.T) {
		queries := usersQueriesForEmails(base, []string{"alice@example.com", "bob@example.com"})

		for i, q := range queries {
			if q.DirectoryID == nil || *q.DirectoryID != directory {
				t.Fatalf("query %d lost the directory filter: %v", i, q.DirectoryID)
			}
		}
	})

	t.Run("leaves the base query unmodified", func(t *testing.T) {
		usersQueriesForEmails(base, []string{"alice@example.com"})

		if base.Email != nil {
			t.Fatalf("expected the base query to be untouched, got email %v", *base.Email)
		}
	})
}
