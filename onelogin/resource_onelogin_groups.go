package onelogin

import (
	"context"
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
	_, err := groupschema.Inflate(map[string]interface{}{
		"name":      d.Get("name"),
		"reference": d.Get("reference"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryCreate, "Group", "")
	}

	// Note: The OneLogin SDK doesn't currently have a CreateGroup method
	// This is a placeholder for when that functionality is added
	// For now, we'll return an error indicating this isn't implemented yet
	return diag.Errorf("Creating groups is not yet supported by the OneLogin API")

	// Once the SDK supports group creation, the code would look something like:
	/*
		client := m.(*onelogin.OneloginSDK)
		tflog.Info(ctx, "[CREATE] Creating group", map[string]interface{}{
			"name": group.Name,
		})

		result, err := client.CreateGroup(group)
		if err != nil {
			return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Group", "")
		}

		// Extract group ID from the result
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
			"id": groupID,
		})

		d.SetId(fmt.Sprintf("%d", groupID))
		return groupRead(ctx, d, m)
	*/
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

	result, err := client.GetGroupByID(groupID)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Group", d.Id())
	}

	// Check if group exists
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] Group not found", map[string]interface{}{
			"id": groupID,
		})
		d.SetId("")
		return nil
	}

	// Parse the group map from the result
	groupMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse group response")
	}

	// Set basic fields
	if name, ok := groupMap["name"].(string); ok {
		d.Set("name", name)
	}

	if ref, ok := groupMap["reference"].(string); ok {
		d.Set("reference", ref)
	}

	return nil
}

// groupUpdate updates a OneLogin Group
func groupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = groupschema.Inflate(map[string]interface{}{
		"id":        groupID,
		"name":      d.Get("name"),
		"reference": d.Get("reference"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryUpdate, "Group", d.Id())
	}

	// Note: The OneLogin SDK doesn't currently have an UpdateGroup method
	// This is a placeholder for when that functionality is added
	return diag.Errorf("Updating groups is not yet supported by the OneLogin API")

	// Once the SDK supports group updates, the code would look something like:
	/*
		client := m.(*onelogin.OneloginSDK)
		tflog.Info(ctx, "[UPDATE] Updating group", map[string]interface{}{
			"id": groupID,
		})

		_, err = client.UpdateGroup(groupID, group)
		if err != nil {
			return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Group", d.Id())
		}

		tflog.Info(ctx, "[UPDATED] Updated group", map[string]interface{}{
			"id": groupID,
		})

		return groupRead(ctx, d, m)
	*/
}

// groupDelete deletes a OneLogin Group
func groupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: The OneLogin SDK doesn't currently have a DeleteGroup method
	// This is a placeholder for when that functionality is added
	return diag.Errorf("Deleting groups is not yet supported by the OneLogin API")

	// Once the SDK supports group deletion, the code would look something like:
	/*
		client := m.(*onelogin.OneloginSDK)
		return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
			aid, _ := strconv.Atoi(id)
			return client.DeleteGroup(aid)
		}, "Group")
	*/
}
