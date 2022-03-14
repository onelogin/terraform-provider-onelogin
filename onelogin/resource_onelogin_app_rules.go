package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	apprulesschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules"
	appruleactionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/actions"
	appruleconditionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/conditions"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// AppRules returns a resource with the CRUD methods and Terraform Schema defined
func AppRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: appRuleCreate,
		ReadContext:   appRuleRead,
		UpdateContext: appRuleUpdate,
		DeleteContext: appRuleDelete,
		Importer: &schema.ResourceImporter{
			// State is added here, which splits app_id and rules_id
			// and sets them appropriately before reading
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				app_id, rule_id, err := utils.ParseNestedResourceImportId(d.Id())
				if err != nil {
					return nil, err
				}
				d.SetId(rule_id)
				d.Set("app_id", app_id)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: apprulesschema.Schema(),
	}
}

// appRuleCreate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the POST request to OneLogin to create an App with its sub-resources
func appRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	appRule := apprulesschema.Inflate(map[string]interface{}{
		"app_id":     d.Get("app_id"),
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"position":   d.Get("position"),
		"enabled":    d.Get("enabled"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	client := m.(*client.APIClient)
	err := client.Services.AppRulesV2.Create(&appRule)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem creating the app rule! %v", err)
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "[CREATED] Created app rule with", *(appRule.ID))

	d.SetId(fmt.Sprintf("%d", *(appRule.ID)))
	return appRuleRead(ctx, d, m)
}

// appRuleRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an App with its sub-resources
func appRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.APIClient)
	id, _ := strconv.Atoi(d.Id())
	appID, _ := strconv.Atoi(d.Get("app_id").(string))
	app, err := client.Services.AppRulesV2.GetOne(int32(appID), int32(id))
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem reading the app rule!", err)
		return diag.FromErr(err)
	}
	if app == nil {
		d.SetId("")
		return nil
	}
	tflog.Info(ctx, "[READ] Reading app rule with %d", *(app.ID))

	d.Set("name", app.Name)
	d.Set("match", app.Match)
	d.Set("position", app.Position)
	d.Set("enabled", app.Enabled)

	d.Set("conditions", appruleconditionsschema.Flatten(app.Conditions))
	d.Set("actions", appruleactionsschema.Flatten(app.Actions))

	return nil
}

// appRuleUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an App and its sub-resources
func appRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	appRule := apprulesschema.Inflate(map[string]interface{}{
		"id":         d.Id(),
		"app_id":     d.Get("app_id"),
		"name":       d.Get("name"),
		"match":      d.Get("match"),
		"position":   d.Get("position"),
		"enabled":    d.Get("enabled"),
		"conditions": d.Get("conditions"),
		"actions":    d.Get("actions"),
	})
	client := m.(*client.APIClient)

	err := client.Services.AppRulesV2.Update(&appRule)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem updating the app rule!", err)
		return diag.FromErr(err)
	}
	if appRule.ID == nil { // app must be deleted in api so remove from tf state
		d.SetId("")
		return nil
	}
	tflog.Info(ctx, "[UPDATED] Updated app rule with %d", *(appRule.ID))
	d.SetId(fmt.Sprintf("%d", *(appRule.ID)))
	return appRuleRead(ctx, d, m)
}

// appRuleDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an App and its sub-resources
func appRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, _ := strconv.Atoi(d.Id())
	appID, _ := strconv.Atoi(d.Get("app_id").(string))
	client := m.(*client.APIClient)

	err := client.Services.AppRulesV2.Destroy(int32(appID), int32(id))
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem deleting the app rule!", err)
		return diag.FromErr(err)
	} else {
		tflog.Info(ctx, "[DELETED] Deleted app rule with %d", id)
		d.SetId("")
	}

	return nil
}
