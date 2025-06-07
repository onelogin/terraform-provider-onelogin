package onelogin

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// OneLoginGroup returns a resource with the OneLogin Group schema
func dataSourceOneLoginGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOneLoginGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceOneLoginGroupRead fetches a group from the OneLogin API by ID
func dataSourceOneLoginGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	groupID := d.Get("id").(int)

	tflog.Info(ctx, "[READ] Reading OneLogin Group", map[string]interface{}{
		"id": groupID,
	})

	resp, err := client.GetGroupByID(groupID)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "OneLogin Group", strconv.Itoa(groupID))
	}

	// Parse the response into a Group model
	groupData, ok := resp.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse group response")
	}

	// Set the resource ID
	d.SetId(strconv.Itoa(groupID))

	// Set the group attributes
	if name, ok := groupData["name"].(string); ok {
		d.Set("name", name)
	}

	if ref, ok := groupData["reference"].(string); ok {
		d.Set("reference", ref)
	}

	return nil
}
