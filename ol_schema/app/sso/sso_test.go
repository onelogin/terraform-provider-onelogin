package appssoschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestFlattenOIDCSSO(t *testing.T) {
	tests := map[string]struct {
		InputData      models.SSOOpenId
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SSO fields from an OIDC app": {
			InputData: models.SSOOpenId{
				ClientID: "test",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "test",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenOIDC(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAMLCert(t *testing.T) {
	tests := map[string]struct {
		InputData      models.SSOSAML
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SAML SSO Certificate fields for the given SAML app": {
			InputData: models.SSOSAML{
				MetadataURL: "test",
				AcsURL:      "test",
				SlsURL:      "test",
				Issuer:      "test",
				Certificate: models.Certificate{
					Name:  "test",
					ID:    123,
					Value: "test",
				},
			},
			ExpectedOutput: map[string]interface{}{
				"name":  "test",
				"value": "test",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAMLCert(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAML(t *testing.T) {
	tests := map[string]struct {
		InputData      models.SSOSAML
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SSO fields for a SAML app": {
			InputData: models.SSOSAML{
				MetadataURL: "test",
				AcsURL:      "test",
				SlsURL:      "test",
				Issuer:      "test",
				Certificate: models.Certificate{
					Name:  "test",
					ID:    123,
					Value: "test",
				},
			},
			ExpectedOutput: map[string]interface{}{
				"metadata_url": "test",
				"acs_url":      "test",
				"sls_url":      "test",
				"issuer":       "test",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAML(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenOIDCCredentials(t *testing.T) {
	tests := map[string]struct {
		InputData      map[string]interface{}
		ExpectedOutput map[string]interface{}
	}{
		"returns both credential fields when the API reports them": {
			InputData: map[string]interface{}{
				"client_id":     "test-client-id",
				"client_secret": "test-client-secret",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     "test-client-id",
				"client_secret": "test-client-secret",
			},
		},
		// A public client (PKCE, token_endpoint_auth_method none) has no
		// secret, and a read-only API credential may have the key withheld.
		// Either way the key is absent from the map rather than empty, so
		// "no secret reported" stays distinguishable from "empty secret".
		"omits client_secret when the API does not report one": {
			InputData: map[string]interface{}{
				"client_id": "test-client-id",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "test-client-id",
			},
		},
		// Same contract as the case above, for an API that reports the field
		// present-but-empty rather than absent. Storing "" would put a key in
		// state that reads as a real value at plan time.
		"omits fields the API reports as empty strings": {
			InputData: map[string]interface{}{
				"client_id":     "test-client-id",
				"client_secret": "",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "test-client-id",
			},
		},
		"returns an empty map for an empty sso object": {
			InputData:      map[string]interface{}{},
			ExpectedOutput: map[string]interface{}{},
		},
		// Guards the type assertions: a non-string value is skipped, not
		// written through as a wrong-typed entry.
		"skips fields the API did not return as strings": {
			InputData: map[string]interface{}{
				"client_id":     "test-client-id",
				"client_secret": 12345,
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "test-client-id",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenOIDCCredentials(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestRetainSecret(t *testing.T) {
	tests := map[string]struct {
		Prior          interface{}
		Flattened      map[string]interface{}
		ExpectedOutput map[string]interface{}
	}{
		// The normal refresh: the read endpoint never returns the secret, so
		// the value captured at create time has to survive every read.
		"carries a prior client_secret forward when the response omits one": {
			Prior: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "captured-at-create",
			},
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "captured-at-create",
			},
		},
		// An imported app has no create response, so there is nothing to
		// carry forward and the secret is simply unavailable.
		"omits client_secret when neither the response nor prior state has one": {
			Prior: map[string]interface{}{
				"client_id": "cid",
			},
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "cid",
			},
		},
		// If OneLogin ever starts returning the secret, the API wins.
		"prefers the response's client_secret over the prior one": {
			Prior: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "stale",
			},
			Flattened: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "fresh",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "fresh",
			},
		},
		// client_id always comes from the response; only the secret is retained.
		"drops the prior secret when the response reports a different client_id": {
			// The app's OIDC client was re-issued, so the captured secret belongs to
			// credentials that no longer exist. Callers persist this value into other
			// systems, so an absent secret is safer than a confidently wrong one.
			Prior: map[string]interface{}{
				"client_id":     "old-cid",
				"client_secret": "s",
			},
			Flattened: map[string]interface{}{
				"client_id": "new-cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "new-cid",
			},
		},
		"retains the prior secret when the response reports the same client_id": {
			Prior: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "s",
			},
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "s",
			},
		},
		"keeps a known client_id the response did not report": {
			// d.Set replaces the whole map, so a response without client_id would
			// otherwise erase it from state with no way to get it back.
			Prior: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "s",
			},
			Flattened: map[string]interface{}{},
			ExpectedOutput: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "s",
			},
		},
		"does not make an empty prior secret permanent": {
			Prior: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "",
			},
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "cid",
			},
		},
		"drops a prior secret that has no client_id to confirm it against": {
			// Legacy flatmap state can carry a secret with no client_id beside it.
			// Retaining it would pair it with whatever the response reports, which
			// is the mismatched pair this function exists to prevent -- and it
			// would happen silently, since a secret would be present either way.
			Prior: map[string]interface{}{
				"client_secret": "s",
			},
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "cid",
			},
		},
		"retains a secret when neither side reports a client_id": {
			// Nothing contradicts the pairing, and there is no id to erase.
			Prior: map[string]interface{}{
				"client_secret": "s",
			},
			Flattened: map[string]interface{}{},
			ExpectedOutput: map[string]interface{}{
				"client_secret": "s",
			},
		},
		"handles a nil prior": {
			Prior: nil,
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "cid",
			},
		},
		"handles a prior of an unexpected type": {
			Prior: "not a map",
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "cid",
			},
		},
		"ignores a non-string prior secret": {
			Prior: map[string]interface{}{
				"client_secret": 12345,
			},
			Flattened: map[string]interface{}{
				"client_id": "cid",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "cid",
			},
		},
		// The secret exists in no response but the create one, so an empty
		// string on a later read must not be treated as OneLogin actively
		// clearing it -- that would destroy the captured secret with no way
		// to recover it.
		"retains the prior secret when the response carries an empty one": {
			Prior: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "captured-at-create",
			},
			Flattened: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     "cid",
				"client_secret": "captured-at-create",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := RetainSecret(test.Prior, test.Flattened)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
