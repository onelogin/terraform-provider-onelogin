package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	groupschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/group"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// OneLoginGroups returns a resource with the OneLogin Groups schema
func resourceOneLoginGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: groupCreate,
		ReadContext:   groupRead,
		UpdateContext: groupUpdate,
		DeleteContext: groupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The API accepts reference in a create or update payload and
			// then neither stores nor returns it. Computed keeps an omitted
			// value from diffing, and the read below leaves a set one alone
			// rather than clearing it, so a configuration that uses the
			// attribute still settles. Whether it should exist at all, given
			// the API does nothing with it, is a question for the next major
			// alongside the sso change held in #236.
			"reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The user policy applied to this group's members. Optional and
			// deliberately not Computed: Computed reads an empty
			// configuration as "keep what is in state", which would leave no
			// way to say "no policy" once one had been set.
			//
			// Only user policies are assignable; an app policy is refused
			// with 422 "Policy must reference a user policy".
			"policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// groupCreate creates a new OneLogin Group
func groupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	inflateMap := map[string]interface{}{
		"name": d.Get("name"),
	}
	if ref, ok := d.GetOk("reference"); ok {
		inflateMap["reference"] = ref
	}
	// A new group with no policy simply does not mention one; sending a 0
	// would store that 0 rather than leaving the column null.
	if policyID, ok := d.GetOk("policy_id"); ok {
		inflateMap["policy_id"] = policyID
	}

	group, err := groupschema.Inflate(inflateMap)
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryCreate, "Group", "")
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating group", map[string]interface{}{
		"name": group.Name,
	})

	result, err := client.CreateGroupWithContext(ctx, &group)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Group", "")
	}

	groupMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse group creation response")
	}

	id, ok := groupMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract group ID from response")
	}

	groupID := int(id)
	tflog.Info(ctx, "[CREATED] Created group", map[string]interface{}{
		"id":   groupID,
		"name": group.Name,
	})

	d.SetId(fmt.Sprintf("%d", groupID))
	return groupRead(ctx, d, m)
}

// groupRead reads a OneLogin Group by ID
func groupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "[READ] Reading group", map[string]interface{}{
		"id": groupID,
	})

	result, err := client.GetGroupByIDV2WithContext(ctx, groupID)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Group", d.Id())
	}

	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] Group not found", map[string]interface{}{
			"id": groupID,
		})
		d.SetId("")
		return nil
	}

	groupMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse group response")
	}

	if err := d.Set("name", groupMap["name"]); err != nil {
		return diag.FromErr(err)
	}
	// Absent stays absent. The API drops reference rather than storing it, so
	// it comes back null on every read, and writing that null cleared whatever
	// the configuration had asked for -- a removal every plan proposed and the
	// next read produced again.
	if reference, present := groupMap["reference"]; present && reference != nil {
		if err := d.Set("reference", reference); err != nil {
			return diag.FromErr(err)
		}
	}

	// An unassigned group reports policy_id as null and a cleared one reports
	// it as 0. Both land on the zero value, which is also what an absent
	// attribute holds, so dropping policy_id from a configuration converges on
	// the next read instead of proposing the same change forever.
	policyID := 0
	if raw, present := groupMap["policy_id"]; present && raw != nil {
		if number, ok := raw.(float64); ok {
			policyID = int(number)
		}
	}
	if err := d.Set("policy_id", policyID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// groupUpdate updates a OneLogin Group
func groupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	inflateMap := map[string]interface{}{
		"id":   groupID,
		"name": d.Get("name"),
	}
	if ref, ok := d.GetOk("reference"); ok {
		inflateMap["reference"] = ref
	}
	// Sent only when it actually changed, and then whatever it now is --
	// including the 0 that clears the assignment. Leaving it out of every
	// other update is what stops a rename from disturbing the policy.
	if d.HasChange("policy_id") {
		inflateMap["policy_id"] = d.Get("policy_id")
	}

	group, err := groupschema.Inflate(inflateMap)
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryUpdate, "Group", d.Id())
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[UPDATE] Updating group", map[string]interface{}{
		"id": groupID,
	})

	_, err = client.UpdateGroupWithContext(ctx, groupID, &group)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Group", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated group", map[string]interface{}{
		"id": groupID,
	})

	return groupRead(ctx, d, m)
}

// groupDelete deletes a OneLogin Group
func groupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		aid, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		return client.DeleteGroupWithContext(ctx, aid)
	}, "Group")
}
