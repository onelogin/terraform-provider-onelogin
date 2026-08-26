package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	appschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app"
	appconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
	appparametersschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	appprovisioningschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
	appssoschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/sso"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// OIDCApps attaches additional configuration schema and
// returns a resource with the CRUD methods and Terraform Schema defined
func OIDCApps() *schema.Resource {
	appSchema := appschema.Schema()
	appSchema["configuration"] = &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	// Computed only, so a configuration cannot set it, and Sensitive because
	// the map carries the client secret. Every sso value is API-supplied --
	// appschema.Inflate has never sent any of them. TypeMap cannot mark a
	// single key sensitive, so the whole map is redacted; client_id needs
	// nonsensitive() to reach a non-sensitive output.
	appSchema["sso"] = &schema.Schema{
		Type:      schema.TypeMap,
		Computed:  true,
		Elem:      &schema.Schema{Type: schema.TypeString},
		Sensitive: true,
	}
	return &schema.Resource{
		CreateContext: oidcAppCreate,
		ReadContext:   oidcAppRead,
		UpdateContext: oidcAppUpdate,
		DeleteContext: oidcAppDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        appSchema,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    oidcAppsV0().CoreConfigSchema().ImpliedType(),
				Upgrade: appparametersschema.UpgradeParameterValuesV0,
			},
		},
	}
}

// oidcAppsV0 describes state written before a parameter's values and
// default_values became lists. configuration is unchanged and repeated here
// only so the shape matches what was written.
//
// sso is declared here defensively, not because decoding requires it.
// StateUpgrader.Type is only consulted on the flatmap (Terraform 0.11) upgrade
// path; JSON state goes straight into Upgrade as a raw map, and whatever
// survives is then filtered against the *current* schema. Declaring sso in the
// V0 schema keeps it out of that flatmap path's blind spot, so pre-refactor
// state -- which carried it as a TypeMap -- still decodes cleanly instead of
// relying on the current schema alone.
func oidcAppsV0() *schema.Resource {
	appSchema := appschema.SchemaV0()
	appSchema["configuration"] = &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	appSchema["sso"] = &schema.Schema{
		Type:      schema.TypeMap,
		Computed:  true,
		Elem:      &schema.Schema{Type: schema.TypeString},
		Sensitive: true,
	}
	return &schema.Resource{Schema: appSchema}
}

// oidcAppCreate creates an OIDC app with all sub-resources
func oidcAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	inflateMap := map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
	}
	addAppPolicyIDForCreate(d, inflateMap)

	oidcApp, err := appschema.Inflate(inflateMap)
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryCreate, "OIDC App", "")
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating OIDC app", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	result, err := client.CreateApp(oidcApp)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "OIDC App", "")
	}

	// Extract app ID from the result
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse OIDC app creation response")
	}

	id, ok := appMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract OIDC app ID from response")
	}

	appID := int(id)
	tflog.Info(ctx, "[CREATED] Created OIDC app", map[string]interface{}{
		"id": appID,
	})

	d.SetId(fmt.Sprintf("%d", appID))

	// The client secret appears in exactly one response: this one. OneLogin
	// returns sso.client_secret only when an app is created -- the app read
	// endpoint returns sso.client_id alone
	// (https://developers.onelogin.com/api-docs/2/apps/app-resource). Capture
	// it here or lose it permanently; oidcAppRead then preserves it on every
	// later refresh via appssoschema.RetainSecret.
	//
	// Both the missing-sso case and a failed Set are reported rather than passed
	// over: this is the only moment the secret is recoverable, so a silent miss
	// leaves the practitioner with an app whose secret cannot be retrieved and no
	// indication of why.
	ssoData, ok := appMap["sso"].(map[string]interface{})
	if !ok {
		tflog.Warn(ctx, "[CREATED] OIDC app response carried no sso object; the client secret could not be captured and cannot be retrieved later", map[string]interface{}{
			"id": appID,
		})
		return oidcAppRead(ctx, d, m)
	}

	if err := d.Set("sso", appssoschema.FlattenOIDCCredentials(ssoData)); err != nil {
		return diag.FromErr(err)
	}

	return oidcAppRead(ctx, d, m)
}

// oidcAppRead reads an OIDC app with all sub-resources
func oidcAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*onelogin.OneloginSDK)
	aid, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "[READ] Reading OIDC app", map[string]interface{}{
		"id": aid,
	})

	result, err := client.GetAppByID(aid, nil)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "OIDC App", d.Id())
	}

	// Check if app exists
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] OIDC app not found", map[string]interface{}{
			"id": aid,
		})
		d.SetId("")
		return nil
	}

	// Parse the app map from the result
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse OIDC app response")
	}

	// Set basic fields
	basicFields := []string{
		"name", "visible", "description", "notes", "icon_url",
		"auth_method", "policy_id", "allow_assumed_signin", "tab_id",
		"brand_id", "connector_id", "created_at", "updated_at",
	}
	utils.SetResourceFields(d, appMap, basicFields)

	// Handle parameters if they exist
	if v, ok := appMap["parameters"]; ok {
		if params, ok := v.(map[string]interface{}); ok {
			d.Set("parameters", appparametersschema.RetainManaged(d.Get("parameters"), appparametersschema.FlattenV4(params)))
		}
	}

	// Handle provisioning if it exists
	if v, ok := appMap["provisioning"]; ok {
		if provData, ok := v.(map[string]interface{}); ok {
			d.Set("provisioning", appprovisioningschema.FlattenMap(provData))
		}
	}

	// Handle configuration if it exists
	if v, ok := appMap["configuration"]; ok {
		if configData, ok := v.(map[string]interface{}); ok {
			d.Set("configuration", appconfigurationschema.RetainManaged(d.Get("configuration"), appconfigurationschema.Flatten(configData)))
		}
	}

	// Handle sso if it exists. client_id comes from the response; client_secret
	// is retained from state, because OneLogin returns it only at create time.
	if sso, ok := mergeSSOFromRead(d.Get("sso"), appMap); ok {
		diags = append(diags, warnIfSecretDropped(ctx, aid, d.Get("sso"), sso)...)

		if err := d.Set("sso", sso); err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// mergeSSOFromRead combines the sso object from a read response with what state
// already holds, returning false when the response carries no usable sso object.
//
// This is a named function rather than an inline block so that the resource and
// its tests exercise the same code. A test that re-implements the merge would
// keep passing if the RetainSecret call were removed from the read path, which
// is precisely the regression worth guarding.
func mergeSSOFromRead(prior interface{}, appMap map[string]interface{}) (map[string]interface{}, bool) {
	ssoData, ok := appMap["sso"].(map[string]interface{})
	if !ok {
		return nil, false
	}

	return appssoschema.RetainSecret(prior, appssoschema.FlattenOIDCCredentials(ssoData)), true
}

// warnIfSecretDropped reports the one case worth interrupting a plan for: state
// held a client_secret and this read discarded it, because the app's OIDC client
// was re-issued and the retained secret no longer belongs to it.
//
// It deliberately stays quiet when there was never a secret to begin with. An
// imported app and a public client (PKCE, token_endpoint_auth_method "none")
// both legitimately have none, and warning on every plan for the lifetime of
// those resources would be noise -- and for a public client the advice would be
// wrong, since it has no secret to recover. Those cases are a debug line.
func warnIfSecretDropped(ctx context.Context, aid int, prior interface{}, merged map[string]interface{}) diag.Diagnostics {
	if secret, _ := merged["client_secret"].(string); secret != "" {
		return nil
	}

	priorMap, _ := prior.(map[string]interface{})
	priorSecret, _ := priorMap["client_secret"].(string)

	if priorSecret == "" {
		tflog.Debug(ctx, "[READ] OIDC app has no client_secret in state; it can only be captured when the app is created", map[string]interface{}{
			"id": aid,
		})
		return nil
	}

	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "OIDC app client_secret discarded because the app's credentials were re-issued",
		Detail: "The client_id returned by OneLogin no longer matches the one the stored client_secret was captured for, so the " +
			"secret has been removed from state rather than paired with credentials it does not belong to. OneLogin returns " +
			"sso.client_secret only when an app is created, so it cannot be read back -- recreate the app if Terraform needs a " +
			"valid secret in state. Anything already given the old secret must be updated separately.",
	}}
}

// oidcAppUpdate updates an OIDC app with all sub-resources
func oidcAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aid, _ := strconv.Atoi(d.Id())

	inflateMap := map[string]interface{}{
		"id":                   d.Id(),
		"name":                 d.Get("name"),
		"description":          d.Get("description"),
		"notes":                d.Get("notes"),
		"connector_id":         d.Get("connector_id"),
		"visible":              d.Get("visible"),
		"allow_assumed_signin": d.Get("allow_assumed_signin"),
		"parameters":           d.Get("parameters"),
		"provisioning":         d.Get("provisioning"),
		"configuration":        d.Get("configuration"),
	}
	addAppPolicyIDForUpdate(d, inflateMap)

	oidcApp, err := appschema.Inflate(inflateMap)
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryUpdate, "OIDC App", d.Id())
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[UPDATE] Updating OIDC app", map[string]interface{}{
		"id": aid,
	})

	_, err = client.UpdateApp(aid, oidcApp)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "OIDC App", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated OIDC app", map[string]interface{}{
		"id": aid,
	})

	return oidcAppRead(ctx, d, m)
}

// oidcAppDelete deletes an OIDC app
func oidcAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		aid, _ := strconv.Atoi(id)
		return client.DeleteApp(aid)
	}, "OIDC App")
}
