package onelogin

import (
	"context"
	"fmt"
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

	// api/1, not api/2: api/2/groups/{id} leaves reference out of its
	// response even when the group has one, so only api/1 can fill it in.
	resp, err := client.GetGroupByID(groupID)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "OneLogin Group", strconv.Itoa(groupID))
	}

	group, err := groupFromResponse(resp, groupID)
	if err != nil {
		return diag.Errorf("failed to read OneLogin Group %d: %s", groupID, err)
	}

	d.SetId(strconv.Itoa(groupID))

	if err := d.Set("name", group["name"]); err != nil {
		return diag.FromErr(err)
	}

	// A group without a reference reports it as null, which sets "".
	reference, _ := group["reference"].(string)
	if err := d.Set("reference", reference); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// groupFromResponse pulls one group out of the api/1/groups/{id} response,
// which wraps it in a single-element list:
//
//	{"status": {...}, "data": [{"id": 1, "name": "...", "reference": null}]}
//
// Looking for name at the top level of that wrapper found nothing, so the data
// source returned an empty name for every group (#282). A bare group object is
// still accepted, in case the SDK moves to api/2, which returns one.
//
// The group has to be the one asked for and has to have a name. Anything else
// is an error rather than an empty attribute a configuration might act on.
func groupFromResponse(resp interface{}, groupID int) (map[string]interface{}, error) {
	body, ok := resp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response shape %T", resp)
	}

	group := body
	if data, wrapped := body["data"]; wrapped {
		items, ok := data.([]interface{})
		if !ok || len(items) != 1 {
			return nil, fmt.Errorf("expected one group in the response, got %v", data)
		}
		if group, ok = items[0].(map[string]interface{}); !ok {
			return nil, fmt.Errorf("unexpected group shape %T", items[0])
		}
	}

	if id, ok := group["id"].(float64); !ok || int(id) != groupID {
		return nil, fmt.Errorf("the response is for group %v, not %d", group["id"], groupID)
	}
	if name, ok := group["name"].(string); !ok || name == "" {
		return nil, fmt.Errorf("the response has no group name")
	}

	return group, nil
}
