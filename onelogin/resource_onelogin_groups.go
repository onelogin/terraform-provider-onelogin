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
			// then neither stores nor returns it, so a configuration setting
			// it plans the same change on every run. Computed keeps an omitted
			// value from diffing; a set one still will. Whether this attribute
			// should exist at all is a question for the next major, alongside
			// the sso change held in #236.
			"reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	if err := d.Set("reference", groupMap["reference"]); err != nil {
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
