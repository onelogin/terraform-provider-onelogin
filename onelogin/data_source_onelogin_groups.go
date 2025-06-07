package onelogin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// OneLoginGroups returns a resource with the OneLogin Groups schema
func dataSourceOneLoginGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOneLoginGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
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
				},
			},
		},
	}
}

// dataSourceOneLoginGroupsRead fetches all groups from the OneLogin API
func dataSourceOneLoginGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	tflog.Info(ctx, "[READ] Reading OneLogin Groups")
	resp, err := client.GetGroups()
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "OneLogin Groups", "")
	}

	// Parse the response into a slice of Group models
	groupsData, ok := resp.([]interface{})
	if !ok {
		return diag.Errorf("failed to parse groups response")
	}

	groups := make([]models.Group, 0, len(groupsData))
	for _, groupData := range groupsData {
		if groupMap, ok := groupData.(map[string]interface{}); ok {
			var group models.Group

			// Extract ID
			if id, ok := groupMap["id"].(float64); ok {
				group.ID = int(id)
			}

			// Extract Name
			if name, ok := groupMap["name"].(string); ok {
				group.Name = name
			}

			// Extract Reference (if present)
			if ref, ok := groupMap["reference"].(string); ok {
				refPtr := ref
				group.Reference = &refPtr
			}

			groups = append(groups, group)
		}
	}

	// Flatten the groups for Terraform
	flattenedGroups := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		flatGroup := map[string]interface{}{
			"id":   group.ID,
			"name": group.Name,
		}
		if group.Reference != nil {
			flatGroup["reference"] = *group.Reference
		}
		flattenedGroups = append(flattenedGroups, flatGroup)
	}

	d.SetId(fmt.Sprintf("%d", len(flattenedGroups)))
	if err := d.Set("groups", flattenedGroups); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
