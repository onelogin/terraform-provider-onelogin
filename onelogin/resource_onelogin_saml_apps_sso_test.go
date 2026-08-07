package onelogin

import "testing"

// TestSamlAppSSOStaysConfigurableUntilTheNextMajor pins the deferral of #236.
//
// sso should be Computed: OneLogin documents every app sso attribute as
// read-only, and appschema.Inflate has never sent the field, so configuring it
// writes something the provider silently discards.
//
// Making it Computed rejects any configuration that sets it, which fails a plan
// that previously succeeded. It is held for the next major so it arrives
// alongside the other breaking change rather than on its own, and this test
// exists so the deferral is deliberate rather than forgotten.
//
// When that major comes: invert this test, and re-close #236.
func TestSamlAppSSOStaysConfigurableUntilTheNextMajor(t *testing.T) {
	sso := SAMLApps().Schema["sso"]

	if sso == nil {
		t.Fatal("expected the saml app schema to have an sso attribute")
	}
	if !sso.Optional {
		t.Fatal("expected sso to still be Optional; making it Computed is a breaking change held for the next major (#236)")
	}
	if sso.Computed {
		t.Fatal("sso became Computed, which fails any plan that configures it — that belongs in a major release (#236)")
	}
}
