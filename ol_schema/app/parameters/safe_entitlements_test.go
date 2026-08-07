package appparametersschema

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// TestSafeEntitlementsEnabled covers issue #235.
//
// The schema had declared safe_entitlements_enabled for some time, but nothing
// carried it: the SDK model had no such field, so the value was discarded on
// the way out and never came back on read.
func TestSafeEntitlementsEnabled(t *testing.T) {
	// Fails on the marshal error rather than asserting against an empty
	// string, which would report a missing key when the real fault was the
	// encoding.
	marshal := func(t *testing.T, s map[string]interface{}) string {
		t.Helper()
		b, err := json.Marshal(Inflate(s))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		return string(b)
	}

	t.Run("sends true", func(t *testing.T) {
		got := marshal(t, map[string]interface{}{"label": "Groups", "safe_entitlements_enabled": true})

		if !strings.Contains(got, `"safe_entitlements_enabled":true`) {
			t.Fatalf("expected true to be sent, got %s", got)
		}
	})

	t.Run("sends false rather than dropping it", func(t *testing.T) {
		// The model field is a pointer for exactly this: the flags beside it
		// are bare bools and omitempty makes their false unsendable.
		got := marshal(t, map[string]interface{}{"label": "Groups", "safe_entitlements_enabled": false})

		if !strings.Contains(got, `"safe_entitlements_enabled":false`) {
			t.Fatalf("expected false to be sent, got %s", got)
		}
	})

	t.Run("omitted when absent", func(t *testing.T) {
		got := marshal(t, map[string]interface{}{"label": "Groups"})

		if strings.Contains(got, `"safe_entitlements_enabled":`) {
			t.Fatalf("expected the key to be absent, got %s", got)
		}
	})

	t.Run("survives the round trip through Flatten", func(t *testing.T) {
		enabled := true
		out := Flatten(map[string]models.Parameter{
			"groups": {Label: "Groups", SafeEntitlementsEnabled: &enabled},
		})

		if len(out) != 1 || out[0]["safe_entitlements_enabled"] != true {
			t.Fatalf("expected the value to come back, got %#v", out)
		}
	})

	t.Run("Flatten leaves it out when the API said nothing", func(t *testing.T) {
		out := Flatten(map[string]models.Parameter{"groups": {Label: "Groups"}})

		if _, present := out[0]["safe_entitlements_enabled"]; present {
			t.Fatalf("expected unset to stay absent, got %#v", out[0])
		}
	})

	t.Run("FlattenV4 reads it from the API response", func(t *testing.T) {
		out := FlattenV4(map[string]interface{}{
			"groups": map[string]interface{}{"label": "Groups", "safe_entitlements_enabled": true},
		})

		if len(out) != 1 || out[0]["safe_entitlements_enabled"] != true {
			t.Fatalf("expected the value to be read, got %#v", out)
		}
	})

	t.Run("FlattenV4 leaves a null out", func(t *testing.T) {
		out := FlattenV4(map[string]interface{}{
			"groups": map[string]interface{}{"label": "Groups", "safe_entitlements_enabled": nil},
		})

		if _, present := out[0]["safe_entitlements_enabled"]; present {
			t.Fatalf("expected a null to stay absent, got %#v", out[0])
		}
	})
}
