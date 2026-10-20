package onelogin

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

// TestSAMLAppConfigurationWithoutSignatureAlgorithm covers what onelogin_saml_apps
// sends for a configuration that leaves signature_algorithm out.
//
// It used to send an empty configuration: appschema.Inflate could only tell a
// SAML configuration by that one key. The apply succeeded, the next read did
// not have the values, and every plan proposed them again.
func TestSAMLAppConfigurationWithoutSignatureAlgorithm(t *testing.T) {
	noToken := map[string]cty.Value{bearerTokenWOKey: cty.NullVal(cty.String)}
	configuration := map[string]interface{}{"account_id": "my-domain"}

	t.Run("create", func(t *testing.T) {
		got := samlAppBodyJSON(t, nil, map[string]interface{}{
			"name":          "New Relic",
			"connector_id":  39893,
			"configuration": configuration,
		}, noToken, true)

		if !strings.Contains(got, `"configuration":{"account_id":"my-domain"}`) {
			t.Fatalf("expected the configuration to be sent, got %s", got)
		}
	})

	t.Run("update", func(t *testing.T) {
		got := samlAppBodyJSON(t, scimAppState(t), map[string]interface{}{
			"name":          "SCIM app",
			"connector_id":  45504,
			"configuration": configuration,
		}, noToken, false)

		if !strings.Contains(got, `"configuration":{"account_id":"my-domain"}`) {
			t.Fatalf("expected the configuration to be sent, got %s", got)
		}
	})
}
