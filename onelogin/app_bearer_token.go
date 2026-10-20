package onelogin

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// The SCIM bearer token of a provisioning connector, issue #29.
//
// OneLogin stores it as a provisioning-level parameter of the app. The app
// endpoint takes it under configuration on a create or an update, and stores
// it, but never returns it on a read: only app-level parameters come back in
// configuration. A configuration key the API does not report is a diff every
// plan proposes and no apply can settle, so the token cannot be managed there.
//
// It is write-only instead. Terraform keeps neither the value nor a hash of it
// in plan or state, so there is nothing to compare a read against and nothing
// for a read to undo. That also means the provider cannot tell when the token
// has changed, which is what bearer_token_wo_version is for: a create sends the
// token whenever it is set, and an update sends it only when the version
// changes. A token changed in the OneLogin admin UI is not drift.
const (
	bearerTokenWOKey        = "bearer_token_wo"
	bearerTokenWOVersionKey = "bearer_token_wo_version"

	// The parameter's name on the SCIM connectors that take a bearer token, and
	// so the configuration key the app endpoint expects it under.
	bearerTokenConfigKey = "bearer_token"
)

func addBearerTokenSchema(s map[string]*schema.Schema) {
	s[bearerTokenWOKey] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		WriteOnly:   true,
		Sensitive:   true,
		Description: "SCIM bearer token for the app's provisioning connector. Write-only: never stored in state. Requires Terraform 1.11 or later.",
	}
	// At least 1, because an unset TypeInt reads as 0: after an import, or on a
	// resource applied before this attribute existed, a version of 0 would be
	// no change and the token would never be sent.
	s[bearerTokenWOVersionKey] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		RequiredWith: []string{bearerTokenWOKey},
		ValidateFunc: validation.IntAtLeast(1),
		Description:  "Change this to send bearer_token_wo again on an update. Starts at 1.",
	}
}

// bearerTokenToSend reports the token a create or update should send, if any.
//
// It has to come from the raw configuration: d.Get always reads a write-only
// attribute as its zero value, because the value is never in the plan.
func bearerTokenToSend(d *schema.ResourceData, create bool) (string, bool, error) {
	if !create && !d.HasChange(bearerTokenWOVersionKey) {
		return "", false, nil
	}

	value, diags := d.GetRawConfigAt(cty.GetAttrPath(bearerTokenWOKey))
	if diags.HasError() {
		return "", false, fmt.Errorf("reading %s: %s", bearerTokenWOKey, diags[0].Summary)
	}
	if value.IsNull() || !value.IsKnown() || !value.Type().Equals(cty.String) {
		return "", false, nil
	}
	return value.AsString(), true, nil
}

// withBearerToken adds the token to the configuration appschema.Inflate built.
//
// A SAML app's configuration is always a map, empty when the configuration sets
// nothing. The token goes in either way, because that is the only place the app
// endpoint looks for it. A typed OIDC configuration has nowhere to put it.
func withBearerToken(configuration interface{}, token string) (interface{}, error) {
	if configuration == nil {
		return map[string]interface{}{bearerTokenConfigKey: token}, nil
	}

	conf, ok := configuration.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s cannot be sent with a %T configuration", bearerTokenWOKey, configuration)
	}
	conf[bearerTokenConfigKey] = token
	return conf, nil
}

// addBearerToken puts the token on the app a create or update will send, when
// there is one to send.
func addBearerToken(d *schema.ResourceData, app *models.App, create bool) error {
	token, ok, err := bearerTokenToSend(d, create)
	if err != nil || !ok {
		return err
	}
	app.Configuration, err = withBearerToken(app.Configuration, token)
	return err
}
