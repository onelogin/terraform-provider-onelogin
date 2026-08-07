package appschema

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestInflateDoesNotSendEmptyStringsForUnsetFields covers issue #237.
//
// notes and description are *string tagged omitempty, but a pointer to "" is
// not empty to encoding/json -- only a nil pointer is. Taking the address
// unconditionally therefore sent `"notes": ""` for a field the API had returned
// as null, and the API stores what it is sent, so an update touching only the
// description rewrote unrelated fields.
//
// The assertion is on the marshalled body rather than the struct, because the
// bug only exists once it is serialised.
func TestInflateDoesNotSendEmptyStringsForUnsetFields(t *testing.T) {
	marshal := func(t *testing.T, s map[string]interface{}) string {
		t.Helper()
		app, err := Inflate(s)
		if err != nil {
			t.Fatalf("inflate: %v", err)
		}
		b, err := json.Marshal(app)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		return string(b)
	}

	t.Run("omits notes and description when unset", func(t *testing.T) {
		got := marshal(t, map[string]interface{}{"name": "test app"})

		if strings.Contains(got, `"notes"`) {
			t.Fatalf("expected notes to be omitted, got %s", got)
		}
		if strings.Contains(got, `"description"`) {
			t.Fatalf("expected description to be omitted, got %s", got)
		}
	})

	t.Run("omits notes while description is being changed", func(t *testing.T) {
		// The exact shape reported in #237: an update that sets only the
		// description must not carry notes along with it.
		got := marshal(t, map[string]interface{}{
			"name":        "test app",
			"description": "a new description",
		})

		if !strings.Contains(got, `"description":"a new description"`) {
			t.Fatalf("expected the description to be sent, got %s", got)
		}
		if strings.Contains(got, `"notes"`) {
			t.Fatalf("expected notes to stay out of an unrelated update, got %s", got)
		}
	})

	t.Run("still sends both when they are set", func(t *testing.T) {
		got := marshal(t, map[string]interface{}{
			"name":        "test app",
			"description": "d",
			"notes":       "n",
		})

		if !strings.Contains(got, `"description":"d"`) || !strings.Contains(got, `"notes":"n"`) {
			t.Fatalf("expected both values to be sent, got %s", got)
		}
	})

	t.Run("name is always sent", func(t *testing.T) {
		// name is required and carries no omitempty; it is not part of the fix
		// and must not be caught by it.
		if got := marshal(t, map[string]interface{}{"name": "test app"}); !strings.Contains(got, `"name":"test app"`) {
			t.Fatalf("expected name to be sent, got %s", got)
		}
	})
}
