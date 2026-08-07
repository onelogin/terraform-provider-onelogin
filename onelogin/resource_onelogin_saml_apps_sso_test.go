package onelogin

import "testing"

// TestSamlAppSSOIsReadOnly covers issue #236.
//
// OneLogin documents every app sso attribute as read-only, and
// appschema.Inflate has never sent the field. Leaving it Optional meant
// Terraform accepted configuration the provider then silently discarded --
// config that looks effective and does nothing.
func TestSamlAppSSOIsReadOnly(t *testing.T) {
	sso := SAMLApps().Schema["sso"]

	if sso == nil {
		t.Fatal("expected the saml app schema to have an sso attribute")
	}
	if sso.Optional {
		t.Fatal("expected sso to be read-only, the API does not accept it and Inflate never sent it")
	}
	if !sso.Computed {
		t.Fatal("expected sso to be Computed so it is still populated on read and import")
	}
}
