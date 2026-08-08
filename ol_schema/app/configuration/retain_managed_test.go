package appconfigurationschema

import (
	"testing"
)

// TestRetainManaged covers the perpetual diff a TypeMap produces when the API
// answers with keys the practitioner never wrote. There is no per-key Computed
// on a map, so every one of those keys is a removal the plan proposes and the
// next read puts straight back.
func TestRetainManaged(t *testing.T) {
	// What the API returns for an OIDC app configured with four keys: its own
	// defaults come back whether or not they were sent.
	fromAPI := map[string]interface{}{
		"redirect_uri":               "https://localhost:3000/callback",
		"login_url":                  "https://www.test.com",
		"oidc_application_type":      "0",
		"token_endpoint_auth_method": "1",
		"include_amr_claims":         "false",
	}

	t.Run("drops keys the configuration does not have", func(t *testing.T) {
		prior := map[string]interface{}{
			"redirect_uri":               "https://localhost:3000/callback",
			"login_url":                  "https://www.test.com",
			"oidc_application_type":      "0",
			"token_endpoint_auth_method": "1",
		}

		out := RetainManaged(prior, fromAPI)

		if _, present := out["include_amr_claims"]; present {
			t.Fatal("expected include_amr_claims to be left out: the API returns it for every OIDC app, and recording one the configuration never mentioned is a diff no apply can settle")
		}
		if len(out) != len(prior) {
			t.Fatalf("expected the four managed keys, got %v", out)
		}
	})

	t.Run("records drift in a key the configuration does have", func(t *testing.T) {
		out := RetainManaged(map[string]interface{}{
			"login_url": "https://www.stale.com",
		}, fromAPI)

		if out["login_url"] != "https://www.test.com" {
			t.Fatalf("expected the value from the API, got %v", out["login_url"])
		}
	})

	t.Run("keeps a managed key the API stopped reporting out of state", func(t *testing.T) {
		out := RetainManaged(map[string]interface{}{"post_logout_redirect_uri": "https://gone"}, fromAPI)

		if _, present := out["post_logout_redirect_uri"]; present {
			t.Fatal("expected a key absent from the response to be absent from state: the provider has no value to claim for it")
		}
	})

	t.Run("writes everything on import", func(t *testing.T) {
		// An empty prior state is an import: there is no configuration to
		// respect yet, so the whole object is recorded.
		for _, prior := range []interface{}{map[string]interface{}{}, nil, "not a map"} {
			out := RetainManaged(prior, fromAPI)

			if len(out) != len(fromAPI) {
				t.Fatalf("expected the whole configuration for prior %#v, got %v", prior, out)
			}
		}
	})
}
