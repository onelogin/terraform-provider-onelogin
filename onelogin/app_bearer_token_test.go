package onelogin

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestSAMLAppBearerTokenSchema pins what makes the token safe to configure:
// write-only, so it never reaches plan or state, and not Computed, which a
// write-only attribute cannot be and which would let a read invent a value.
func TestSAMLAppBearerTokenSchema(t *testing.T) {
	attributes := SAMLApps().CoreConfigSchema().Attributes

	token, ok := attributes[bearerTokenWOKey]
	if !ok {
		t.Fatalf("%s is missing from the core schema", bearerTokenWOKey)
	}
	if !token.WriteOnly {
		t.Errorf("%s is not WriteOnly, so the token would be kept in state", bearerTokenWOKey)
	}
	if !token.Sensitive {
		t.Errorf("%s is not Sensitive", bearerTokenWOKey)
	}
	if !token.Optional || token.Computed {
		t.Errorf("%s must be Optional and not Computed, got Optional=%v Computed=%v", bearerTokenWOKey, token.Optional, token.Computed)
	}

	version, ok := attributes[bearerTokenWOVersionKey]
	if !ok {
		t.Fatalf("%s is missing from the core schema", bearerTokenWOVersionKey)
	}
	if version.WriteOnly || !version.Optional || version.Computed {
		t.Errorf("%s must be a plain Optional attribute kept in state, got WriteOnly=%v Optional=%v Computed=%v",
			bearerTokenWOVersionKey, version.WriteOnly, version.Optional, version.Computed)
	}
}

// TestSAMLAppBearerTokenVersionNeedsToken covers RequiredWith: a version with
// no token to send would be a change that sends nothing.
func TestSAMLAppBearerTokenVersionNeedsToken(t *testing.T) {
	diags := schema.InternalMap(SAMLApps().Schema).Validate(terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":                  "SCIM app",
		"connector_id":          45504,
		bearerTokenWOVersionKey: 1,
	}))
	if !diags.HasError() {
		t.Fatalf("expected %s without %s to be rejected", bearerTokenWOVersionKey, bearerTokenWOKey)
	}
}

// samlAppBodyJSON runs a configuration and prior state through the resource's
// own diff and samlAppBody, to the JSON the API would receive.
//
// The write-only value is only ever in the raw configuration, which a diff
// takes from the prior state, so that is where it is put.
func samlAppBodyJSON(t *testing.T, state *terraform.InstanceState, config map[string]interface{}, raw map[string]cty.Value, create bool) string {
	t.Helper()

	r := SAMLApps()
	if state == nil {
		state = &terraform.InstanceState{}
	}
	state.RawConfig = cty.ObjectVal(raw)

	d, err := schema.InternalMap(r.Schema).Data(state, appDiff(t, r, state, config))
	if err != nil {
		t.Fatalf("data: %v", err)
	}

	app, err := samlAppBody(d, create)
	if err != nil {
		t.Fatalf("body: %v", err)
	}
	body, err := json.Marshal(app)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(body)
}

// scimAppState is the state a read leaves behind for an app created with
// bearer_token_wo_version = 1. The token itself is never in it.
func scimAppState(t *testing.T) *terraform.InstanceState {
	t.Helper()

	d := SAMLApps().Data(nil)
	d.SetId("1504547")
	for key, value := range map[string]interface{}{
		"name":                  "SCIM app",
		"connector_id":          45504,
		"configuration":         map[string]interface{}{"signature_algorithm": "SHA-256"},
		bearerTokenWOVersionKey: 1,
	} {
		if err := d.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
	return d.State()
}

func TestSAMLAppBearerTokenBody(t *testing.T) {
	t.Run("create sends the token alongside the configuration", func(t *testing.T) {
		got := samlAppBodyJSON(t, nil, map[string]interface{}{
			"name":          "SCIM app",
			"connector_id":  45504,
			"configuration": map[string]interface{}{"signature_algorithm": "SHA-256"},
		}, map[string]cty.Value{
			bearerTokenWOKey: cty.StringVal("secret-1"),
		}, true)

		for _, want := range []string{`"bearer_token":"secret-1"`, `"signature_algorithm":"SHA-256"`} {
			if !strings.Contains(got, want) {
				t.Errorf("expected %s in %s", want, got)
			}
		}
	})

	// Inflate cannot tell what a configuration without signature_algorithm is
	// and builds an empty one. The token still has to go, and it goes in the
	// only place the app endpoint looks for it.
	t.Run("create sends the token with no other configuration", func(t *testing.T) {
		got := samlAppBodyJSON(t, nil, map[string]interface{}{
			"name":         "SCIM app",
			"connector_id": 45504,
		}, map[string]cty.Value{
			bearerTokenWOKey: cty.StringVal("secret-1"),
		}, true)

		if !strings.Contains(got, `"configuration":{"bearer_token":"secret-1"}`) {
			t.Fatalf("expected the token in configuration, got %s", got)
		}
	})

	t.Run("create without a token sends none", func(t *testing.T) {
		got := samlAppBodyJSON(t, nil, map[string]interface{}{
			"name":          "SCIM app",
			"connector_id":  45504,
			"configuration": map[string]interface{}{"signature_algorithm": "SHA-256"},
		}, map[string]cty.Value{
			bearerTokenWOKey: cty.NullVal(cty.String),
		}, true)

		if strings.Contains(got, "bearer_token") {
			t.Fatalf("expected no bearer_token, got %s", got)
		}
	})

	// The provider cannot see whether the token changed, so an update that
	// leaves the version alone leaves the token alone. Sending it anyway would
	// put back a token someone rotated in the admin UI on every apply.
	t.Run("update with the same version does not send the token", func(t *testing.T) {
		got := samlAppBodyJSON(t, scimAppState(t), map[string]interface{}{
			"name":                  "SCIM app renamed",
			"connector_id":          45504,
			"configuration":         map[string]interface{}{"signature_algorithm": "SHA-256"},
			bearerTokenWOVersionKey: 1,
		}, map[string]cty.Value{
			bearerTokenWOKey: cty.StringVal("secret-1"),
		}, false)

		if strings.Contains(got, "bearer_token") {
			t.Fatalf("expected no bearer_token, got %s", got)
		}
		if !strings.Contains(got, `"name":"SCIM app renamed"`) {
			t.Fatalf("expected the rename to be sent, got %s", got)
		}
	})

	t.Run("update with a new version sends the token", func(t *testing.T) {
		got := samlAppBodyJSON(t, scimAppState(t), map[string]interface{}{
			"name":                  "SCIM app",
			"connector_id":          45504,
			"configuration":         map[string]interface{}{"signature_algorithm": "SHA-256"},
			bearerTokenWOVersionKey: 2,
		}, map[string]cty.Value{
			bearerTokenWOKey: cty.StringVal("secret-2"),
		}, false)

		if !strings.Contains(got, `"bearer_token":"secret-2"`) {
			t.Fatalf("expected the new token to be sent, got %s", got)
		}
	})
}

// TestSAMLAppBearerTokenNotRead guards the other half of write-only: a read
// must not put anything in state for it, whatever the API returns.
func TestSAMLAppBearerTokenNotRead(t *testing.T) {
	for _, key := range appBasicReadFields {
		if key == bearerTokenWOKey || key == bearerTokenWOVersionKey {
			t.Fatalf("%s is in appBasicReadFields, so a read would overwrite it", key)
		}
	}
}

func TestWithBearerTokenRejectsTypedConfiguration(t *testing.T) {
	if _, err := withBearerToken(struct{}{}, "secret"); err == nil {
		t.Fatal("expected a configuration that is not a map to be rejected")
	}
}
