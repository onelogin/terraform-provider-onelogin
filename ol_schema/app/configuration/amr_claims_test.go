package appconfigurationschema

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestIncludeAmrClaims covers issue #238's remaining field.
//
// The configuration block is a map of strings, so the value arrives as "true"
// or "false" and has to leave as a real bool. It is a pointer for the reason
// every other field here is: unset and switched-off are different requests.
func TestIncludeAmrClaims(t *testing.T) {
	inflate := func(t *testing.T, s map[string]interface{}) string {
		t.Helper()
		out, err := Inflate(s)
		if err != nil {
			t.Fatalf("inflate: %v", err)
		}
		b, err := json.Marshal(out)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		return string(b)
	}

	// A configuration with no OIDC hints still has to route down the OIDC
	// branch, which is what redirect_uri signals.
	base := func(extra map[string]interface{}) map[string]interface{} {
		s := map[string]interface{}{"redirect_uri": "https://example.com/cb"}
		for k, v := range extra {
			s[k] = v
		}
		return s
	}

	t.Run("omitted when the key is absent", func(t *testing.T) {
		if got := inflate(t, base(nil)); strings.Contains(got, `"include_amr_claims":`) {
			t.Fatalf("expected the key to be absent, got %s", got)
		}
	})

	t.Run("sends true", func(t *testing.T) {
		if got := inflate(t, base(map[string]interface{}{"include_amr_claims": "true"})); !strings.Contains(got, `"include_amr_claims":true`) {
			t.Fatalf("expected true, got %s", got)
		}
	})

	t.Run("sends false rather than dropping it", func(t *testing.T) {
		// The point of the pointer: with a bare bool this would be
		// indistinguishable from unset and the claim could never be turned off.
		if got := inflate(t, base(map[string]interface{}{"include_amr_claims": "false"})); !strings.Contains(got, `"include_amr_claims":false`) {
			t.Fatalf("expected false to be sent, got %s", got)
		}
	})

	t.Run("rejects a value that is not a boolean", func(t *testing.T) {
		// Dropping it would leave the key in the configuration and out of
		// state: a diff every plan reports and no apply can settle. Erroring
		// says so once, and names the offending value.
		_, err := Inflate(base(map[string]interface{}{"include_amr_claims": "yes please"}))

		if err == nil {
			t.Fatal("expected an unparseable value to be rejected")
		}
		if !strings.Contains(err.Error(), "include_amr_claims") {
			t.Fatalf("expected the field name in the error, got %v", err)
		}
		if !strings.Contains(err.Error(), "yes please") {
			t.Fatalf("expected the offending value in the error, got %v", err)
		}
	})

	t.Run("accepts the spellings ParseBool takes", func(t *testing.T) {
		// Terraform renders a bool into a string map as "true"/"false", but a
		// practitioner may well write "1" or "TRUE" by hand.
		for _, in := range []string{"true", "TRUE", "True", "1", "t"} {
			if _, err := Inflate(base(map[string]interface{}{"include_amr_claims": in})); err != nil {
				t.Fatalf("expected %q to be accepted, got %v", in, err)
			}
		}
	})

	t.Run("flattens back to a string for the map", func(t *testing.T) {
		out := Flatten(map[string]interface{}{
			"redirect_uri":       "https://example.com/cb",
			"include_amr_claims": true,
		})

		if got := out["include_amr_claims"]; got != "true" {
			t.Fatalf("expected the string \"true\", got %#v", got)
		}
	})

	t.Run("leaves the key out when the API returns null", func(t *testing.T) {
		out := Flatten(map[string]interface{}{
			"redirect_uri":       "https://example.com/cb",
			"include_amr_claims": nil,
		})

		if _, present := out["include_amr_claims"]; present {
			t.Fatalf("expected a null to stay absent, got %#v", out["include_amr_claims"])
		}
	})
}
