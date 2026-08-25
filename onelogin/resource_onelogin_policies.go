package onelogin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	policyschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/policy"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// policiesPath is the collection endpoint. The SDK has no policy methods, so
// this resource builds the requests itself on the SDK's authenticated client.
const policiesPath = "/api/2/policies"

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

	path := policiesPath
	resp, err := client.Client.Post(&path, body)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Policy", "")
	}

	policy, err := decodePolicy(resp)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Policy", "")
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

	path := fmt.Sprintf("%s/%s", policiesPath, d.Id())
	resp, err := client.Client.Get(&path, nil)

	var policy map[string]interface{}
	if err == nil {
		policy, err = decodePolicy(resp)
	}

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

	path := fmt.Sprintf("%s/%s", policiesPath, d.Id())
	resp, err := client.Client.Put(&path, body)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Policy", d.Id())
	}
	if _, err := decodePolicy(resp); err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Policy", d.Id())
	}

	return policyRead(ctx, d, m)
}

func policyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		path := fmt.Sprintf("%s/%s", policiesPath, id)
		resp, err := client.Client.Delete(&path)
		if err != nil {
			return nil, err
		}
		// A successful delete is a 204 with no body; only the status matters.
		_, err = policyResponseBody(resp)
		return nil, err
	}, "Policy")
}

// decodePolicy turns a policy response into a map. The endpoint returns the
// policy at the top level rather than wrapped in a key, so the decoded body is
// the policy itself.
func decodePolicy(resp *http.Response) (map[string]interface{}, error) {
	body, err := policyResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var policy map[string]interface{}
	if err := json.Unmarshal(body, &policy); err != nil {
		return nil, fmt.Errorf("expected a policy object in the response: %w", err)
	}
	return policy, nil
}

// policyResponseBody reads a response and turns a failing status into an error
// carrying what the API said.
//
// The SDK's CheckHTTPResponse discards the error body, and this endpoint
// validates hard -- a field belonging to the other kind of policy, a password
// length outside the allowed set, an invite expiry of zero -- answering with a
// message that names the field. "request failed with status: 422" on its own
// leaves the practitioner guessing which of seventy-odd arguments it meant.
func policyResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if closeErr := resp.Body.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// "status: %d" is the wording utils.IsNotFoundError matches on, and
		// both policyRead and the delete helper depend on recognising a 404,
		// so it stays even with a message appended.
		if message := apiErrorMessage(body); message != "" {
			return nil, fmt.Errorf("request failed with status: %d: %s", resp.StatusCode, message)
		}
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return body, nil
}

// apiErrorMessage pulls the explanation out of an error body, which OneLogin
// sends as {"name":..., "message":..., "statusCode":...}. An unrecognisable
// body gives an empty string, leaving the caller with the status alone.
func apiErrorMessage(body []byte) string {
	var apiError struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiError); err != nil {
		return ""
	}
	return apiError.Message
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
