package onelogin

import "testing"

// TestSamlAppSSOIsReadOnly pins sso as Computed only.
//
// This inverts the guard that held the change back. Three schemas are possible
// and the difference between them is who wins when a configuration sets sso:
//
//	Optional only          samlAppRead populates sso from the API, so a
//	                       configuration that omits it has an empty value
//	                       against populated state: a diff on every plan.
//	Optional and Computed  configuration is accepted and then quietly ignored,
//	                       because appschema.Inflate never sends it. What was
//	                       written and what the app has can disagree forever
//	                       with nothing to show for it.
//	Computed only          configuration that sets sso is rejected outright.
//
// #236 asked for the last, and this is the release that can take it. Every sso
// attribute is read-only at the API -- the certificate, the ACS and metadata
// URLs, the issuer -- so a practitioner writing one was describing something
// they never had any way to change.
//
// This is breaking on purpose: a configuration that sets sso now fails to
// validate rather than pretending to work.
func TestSamlAppSSOIsReadOnly(t *testing.T) {
	sso := SAMLApps().Schema["sso"]

	if sso == nil {
		t.Fatal("expected the saml app schema to have an sso attribute")
	}

	if sso.Optional {
		t.Fatal("sso is Optional again, so a configuration can set a value the API will never accept — #236 asked for Computed only")
	}
	if !sso.Computed {
		t.Fatal("expected sso to be Computed: it is read-only, not absent, and read populates it")
	}
	if sso.Required {
		t.Fatal("sso cannot be Required: nothing supplies it but the API")
	}
}
