package onelogin

import "testing"

// TestSamlAppSSOStaysConfigurable pins sso as Optional until a major release.
//
// The distinction that matters is Computed-only, not Computed. Three schemas
// are possible and only one of them breaks anyone:
//
//	Optional only          samlAppRead populates sso from the API, so a
//	                       configuration that omits it has an empty value
//	                       against populated state: a diff on every plan.
//	Optional and Computed  configuration is still accepted, and an omitted
//	                       attribute keeps what the API returned.
//	Computed only          configuration that sets sso is rejected, so a plan
//	                       that used to succeed now fails.
//
// #236 asks for the last, because sso is read-only at the API and
// appschema.Inflate has never sent it. That is breaking and waits for the next
// major. Until then Optional must hold, so this test guards the one property
// that would break a practitioner — not Computed, which is wanted.
//
// When the major comes: drop Optional, set Computed only, invert this test,
// and re-close #236.
func TestSamlAppSSOStaysConfigurable(t *testing.T) {
	sso := SAMLApps().Schema["sso"]

	if sso == nil {
		t.Fatal("expected the saml app schema to have an sso attribute")
	}

	// The whole guarantee. Optional false with Computed true is the breaking
	// shape; Computed on its own is not.
	if !sso.Optional {
		t.Fatal("sso stopped being Optional, which rejects any configuration that sets it — that belongs in a major release (#236)")
	}

	// Not a breaking change, and wanted: without it, a configuration that omits
	// sso diffs against the value samlAppRead writes into state.
	if !sso.Computed {
		t.Fatal("expected sso to be Computed alongside Optional, so omitting it does not diff against what read populates")
	}
}
