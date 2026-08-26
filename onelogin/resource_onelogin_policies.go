package onelogin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	policyschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/policy"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Policies returns the onelogin_policies resource: a OneLogin security policy,
// either the user policy that governs how people sign in or the app policy that
// governs a single application.
func Policies() *schema.Resource {
	return &schema.Resource{
		CreateContext: policyCreate,
		ReadContext:   policyRead,
		UpdateContext: policyUpdate,
		DeleteContext: policyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:        policyschema.Schema(),
		CustomizeDiff: policyCheckKind,
	}
}

// policyCheckKind rejects a configuration that sets a field belonging to the
// other kind of policy.
//
// The API rejects these too, with a 422 naming each field, so this only moves
// the complaint from apply to plan. That is worth doing: the practitioner sees
// it before anything is created, and on an update before a policy that already
// exists is touched.
func policyCheckKind(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	kind, ok := d.Get("kind").(string)
	if !ok {
		return nil
	}

	wrong := policyschema.FieldsNotApplicableTo(kind, policyschema.ConfiguredKeys(d.GetRawConfig()))
	if len(wrong) == 0 {
		return nil
	}

	subject := "fields are"
	if len(wrong) == 1 {
		subject = "field is"
	}
	return fmt.Errorf(
		"%s not applicable to %s policies: %s",
		subject, kind, strings.Join(wrong, ", "),
	)
}

func policyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	body := policyschema.RequestBody(d, policyschema.ConfiguredKeys(d.GetRawConfig()))
	body["kind"] = d.Get("kind")

	tflog.Info(ctx, "[CREATE] Creating policy", map[string]interface{}{
		"name": body["name"],
		"kind": body["kind"],
	})

	result, err := client.CreatePolicyWithContext(ctx, body)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Policy", "")
	}

	policy, ok := policyFromResponse(result)
	if !ok {
		return diag.Errorf("expected a policy object in the create response, got %T", result)
	}

	id, ok := policyID(policy)
	if !ok {
		return diag.Errorf("failed to extract policy ID from the create response")
	}
	d.SetId(id)

	tflog.Info(ctx, "[CREATED] Created policy", map[string]interface{}{"id": id})

	return policyRead(ctx, d, m)
}

func policyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	tflog.Info(ctx, "[READ] Reading policy", map[string]interface{}{"id": d.Id()})

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.GetPolicyByIDWithContext(ctx, id)
	if err != nil {
		// A policy that is gone is not an error on a read. Somebody deleted it
		// outside Terraform; dropping it from state lets the next plan offer
		// to create it again rather than failing every run until it is
		// imported or removed by hand.
		if utils.IsNotFoundError(err) {
			tflog.Warn(ctx, "[READ] Policy is gone, removing it from state", map[string]interface{}{"id": d.Id()})
			d.SetId("")
			return nil
		}
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Policy", d.Id())
	}

	policy, ok := policyFromResponse(result)
	if !ok {
		return diag.Errorf("expected a policy object in the read response, got %T", result)
	}

	if err := policyschema.Flatten(d, policy); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func policyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	// kind is deliberately absent: it is ForceNew, so an update never changes
	// it, and the API rejects a body that names a different one.
	body := policyschema.RequestBody(d, policyschema.ConfiguredKeys(d.GetRawConfig()))

	tflog.Info(ctx, "[UPDATE] Updating policy", map[string]interface{}{"id": d.Id()})

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.UpdatePolicyWithContext(ctx, id, body); err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Policy", d.Id())
	}

	return policyRead(ctx, d, m)
}

func policyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		policyID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		return client.DeletePolicyWithContext(ctx, policyID)
	}, "Policy")
}

// policyFromResponse narrows the SDK's decoded body to the policy itself. The
// endpoint returns the policy at the top level rather than wrapped in a key,
// so the decoded body is the policy.
//
// Reading the response and turning a failing status into an error that carries
// the API's message now happens in the SDK, which this resource used to do for
// itself because there were no policy methods to do it.
func policyFromResponse(result interface{}) (map[string]interface{}, bool) {
	policy, ok := result.(map[string]interface{})
	return policy, ok
}

// policyID reads the ID out of a policy response as the string Terraform keeps.
func policyID(policy map[string]interface{}) (string, bool) {
	switch id := policy["id"].(type) {
	case float64:
		return strconv.Itoa(int(id)), true
	case int:
		return strconv.Itoa(id), true
	case string:
		return id, id != ""
	default:
		return "", false
	}
}
