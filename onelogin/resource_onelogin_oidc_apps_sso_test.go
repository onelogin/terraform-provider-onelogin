package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	appssoschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/sso"
)

// TestOIDCAppSSOIsReadOnlyAndSensitive pins sso as Computed only and Sensitive.
//
// Computed only, because every value is API-supplied: appschema.Inflate never
// reads sso, so an Optional attribute would either diff on every plan or
// accept a value and silently discard it. This matches the guarantee #236
// established for the SAML resource.
//
// Sensitive, because the map carries the OIDC client secret. TypeMap cannot
// mark individual keys, so the whole map is redacted and client_id needs
// nonsensitive() to reach a non-sensitive output. That cost is accepted
// deliberately: a leaked client secret in plan output is the worse outcome.
func TestOIDCAppSSOIsReadOnlyAndSensitive(t *testing.T) {
	sso := OIDCApps().Schema["sso"]

	if sso == nil {
		t.Fatal("expected the oidc app schema to have an sso attribute")
	}

	if sso.Optional {
		t.Fatal("sso cannot be Optional: Inflate never sends it, so a configured value would be silently discarded")
	}
	if !sso.Computed {
		t.Fatal("expected sso to be Computed: it is read-only, not absent, and read populates it")
	}
	if sso.Required {
		t.Fatal("sso cannot be Required: nothing supplies it but the API")
	}
	if !sso.Sensitive {
		t.Fatal("expected sso to be Sensitive: the map carries the OIDC client secret")
	}
	if sso.Type != schema.TypeMap {
		t.Fatalf("expected sso to be a TypeMap, got %s", sso.Type)
	}
}

// TestOIDCAppSSOIsDeclaredInV0 guards shape fidelity, not a decode
// requirement.
//
// Provider versions before the SDK v4 refactor wrote state at SchemaVersion 0
// with sso present as a TypeMap. StateUpgrader.Type is only read on the
// flatmap (Terraform 0.11) upgrade path -- JSON state goes straight into
// Upgrade as a raw map, and is filtered against the current schema afterwards.
// Declaring sso in the v0 schema is defensive: it keeps
// pre-refactor state intact on the flatmap path instead of dropping it, even
// though the JSON path already preserves it via the current schema alone.
func TestOIDCAppSSOIsDeclaredInV0(t *testing.T) {
	sso := oidcAppsV0().Schema["sso"]

	if sso == nil {
		t.Fatal("expected the v0 oidc app schema to declare sso so pre-refactor state decodes")
	}
	if sso.Type != schema.TypeMap {
		t.Fatalf("expected v0 sso to be a TypeMap, got %s", sso.Type)
	}
}

// TestOIDCAppSSOStateRoundTrip checks that the flattener's output survives the
// schema. A Sensitive TypeMap of strings should hold both credential keys and
// hand them back through the usual dotted lookup a practitioner would use.
func TestOIDCAppSSOStateRoundTrip(t *testing.T) {
	d := schema.TestResourceDataRaw(t, OIDCApps().Schema, map[string]interface{}{})

	ssoData := map[string]interface{}{
		"client_id":     "test-client-id",
		"client_secret": "test-client-secret",
	}

	if err := d.Set("sso", appssoschema.FlattenOIDCCredentials(ssoData)); err != nil {
		t.Fatalf("setting sso: %v", err)
	}

	assert.Equal(t, "test-client-id", d.Get("sso.client_id"))
	assert.Equal(t, "test-client-secret", d.Get("sso.client_secret"))
}

// TestOIDCAppSSOStateOmitsAbsentSecret covers the public-client case end to
// end: no client_secret in the response means no client_secret key in state.
//
// Asserts on the whole map rather than on d.Get("sso.client_secret"), because
// key-absent and key-empty are exactly the distinction under test and a dotted
// lookup collapses them.
func TestOIDCAppSSOStateOmitsAbsentSecret(t *testing.T) {
	d := schema.TestResourceDataRaw(t, OIDCApps().Schema, map[string]interface{}{})

	ssoData := map[string]interface{}{
		"client_id": "test-client-id",
	}

	if err := d.Set("sso", appssoschema.FlattenOIDCCredentials(ssoData)); err != nil {
		t.Fatalf("setting sso: %v", err)
	}

	stored, ok := d.Get("sso").(map[string]interface{})
	if !ok {
		t.Fatalf("expected sso to read back as a map, got %T", d.Get("sso"))
	}

	assert.Equal(t, "test-client-id", stored["client_id"])
	assert.NotContains(t, stored, "client_secret", "no secret was reported, so no key should be stored")
}

// TestOIDCAppSSOSecretSurvivesRead pins the behaviour the create-time capture
// depends on: a secret already in state is still there after a read whose
// response carries only client_id. Without RetainSecret this is the exact
// sequence that destroyed the captured secret.
func TestOIDCAppSSOSecretSurvivesRead(t *testing.T) {
	d := schema.TestResourceDataRaw(t, OIDCApps().Schema, map[string]interface{}{})

	// What oidcAppCreate captures from the POST response.
	created := map[string]interface{}{
		"client_id":     "cid",
		"client_secret": "captured-at-create",
	}
	if err := d.Set("sso", appssoschema.FlattenOIDCCredentials(created)); err != nil {
		t.Fatalf("setting sso from the create response: %v", err)
	}

	// What every subsequent read sees: client_id only.
	fromRead := map[string]interface{}{
		"client_id": "cid",
	}
	merged := appssoschema.RetainSecret(d.Get("sso"), appssoschema.FlattenOIDCCredentials(fromRead))
	if err := d.Set("sso", merged); err != nil {
		t.Fatalf("setting sso from the read response: %v", err)
	}

	assert.Equal(t, "cid", d.Get("sso.client_id"))
	assert.Equal(t, "captured-at-create", d.Get("sso.client_secret"))
}
