package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	apprulesschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules"
	appruleactionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/rules/actions"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

const resourceTypeAppRule = "app_rule"

// AppRules returns a resource with the CRUD methods and Terraform Schema defined
// This implementation uses utility functions for better error handling and logging
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
// makes the POST request to OneLogin to create an App Rule
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
	client := m.(*onelogin.OneloginSDK)

	appIDStr := d.Get("app_id").(string)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryCreate,
			resourceTypeAppRule,
			"parsing app_id",
			appIDStr,
			err,
		)
	}

	tflog.Info(ctx, "[CREATE] Creating app rule", map[string]interface{}{
		"app_id": appID,
		"name":   d.Get("name"),
	})

	result, err := client.CreateAppRule(appID, appRule)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, resourceTypeAppRule, appIDStr)
	}

	ruleMap, ok := result.(map[string]interface{})
	if !ok || ruleMap["id"] == nil {
		err := fmt.Errorf("failed to parse app rule creation response or rule ID not found in response")
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryCreate,
			resourceTypeAppRule,
			"parsing response",
			appIDStr,
			err,
		)
	}

	ruleID := int(ruleMap["id"].(float64))
	ruleIDStr := fmt.Sprintf("%d", ruleID)

	tflog.Info(ctx, "[CREATED] Created app rule", map[string]interface{}{
		"app_id":  appID,
		"rule_id": ruleID,
	})

	d.SetId(ruleIDStr)
	return appRuleRead(ctx, d, m)
}

// appRuleRead takes a pointer to the ResourceData Struct and a HTTP client and
// makes the GET request to OneLogin to read an App Rule
func appRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	ruleIDStr := d.Id()
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryRead,
			resourceTypeAppRule,
			"parsing rule_id",
			ruleIDStr,
			err,
		)
	}

	appIDStr := d.Get("app_id").(string)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryRead,
			resourceTypeAppRule,
			"parsing app_id",
			appIDStr,
			err,
		)
	}

	tflog.Info(ctx, "[READ] Reading app rule", map[string]interface{}{
		"app_id":  appID,
		"rule_id": ruleID,
	})

	result, err := client.GetAppRuleByID(appID, ruleID, nil)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, resourceTypeAppRule, ruleIDStr)
	}

	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] App rule not found", map[string]interface{}{
			"app_id":  appID,
			"rule_id": ruleID,
		})
		d.SetId("")
		return nil
	}

	ruleMap, ok := result.(map[string]interface{})
	if !ok {
		err := fmt.Errorf("failed to parse app rule response")
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryRead,
			resourceTypeAppRule,
			"parsing response",
			ruleIDStr,
			err,
		)
	}

	// Set simple fields
	simpleFields := []string{"name", "match", "enabled"}
	utils.SetResourceFields(d, ruleMap, simpleFields)

	// Position is a number that needs special handling
	if pos, ok := ruleMap["position"].(float64); ok {
		d.Set("position", int(pos))
	}

	// Handle conditions
	conditions := []map[string]interface{}{}
	if conditionsList, ok := ruleMap["conditions"].([]interface{}); ok {
		for _, c := range conditionsList {
			if cond, ok := c.(map[string]interface{}); ok {
				condition := map[string]interface{}{
					"source":   cond["source"],
					"operator": cond["operator"],
					"value":    cond["value"],
				}
				conditions = append(conditions, condition)
			}
		}
	}
	d.Set("conditions", conditions)

	// Handle actions
	actions := []map[string]interface{}{}
	if actionsList, ok := ruleMap["actions"].([]interface{}); ok {
		for _, a := range actionsList {
			if act, ok := a.(map[string]interface{}); ok {
				action := map[string]interface{}{
					"action": act["action"],
					"value":  act["value"],
				}

				if act["expression"] != nil && act["expression"] != "" {
					action["expression"] = act["expression"]
				} else {
					action["action"] = fmt.Sprintf("%s%s", act["action"], appruleactionsschema.NO_EXPRESSION_SUFFIX)
				}

				if act["scriplet"] != nil && act["scriplet"] != "" {
					action["scriplet"] = act["scriplet"]
				}

				if act["macro"] != nil && act["macro"] != "" {
					action["macro"] = act["macro"]
				}

				actions = append(actions, action)
			}
		}
	}
	d.Set("actions", actions)

	return nil
}

// appRuleUpdate takes a pointer to the ResourceData Struct and a HTTP client and
// makes the PUT request to OneLogin to update an App Rule
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
	client := m.(*onelogin.OneloginSDK)

	ruleIDStr := d.Id()
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryUpdate,
			resourceTypeAppRule,
			"parsing rule_id",
			ruleIDStr,
			err,
		)
	}

	appIDStr := d.Get("app_id").(string)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryUpdate,
			resourceTypeAppRule,
			"parsing app_id",
			appIDStr,
			err,
		)
	}

	tflog.Info(ctx, "[UPDATE] Updating app rule", map[string]interface{}{
		"app_id":  appID,
		"rule_id": ruleID,
	})

	_, err = client.UpdateAppRule(appID, ruleID, appRule, nil)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, resourceTypeAppRule, ruleIDStr)
	}

	tflog.Info(ctx, "[UPDATED] Updated app rule", map[string]interface{}{
		"app_id":  appID,
		"rule_id": ruleID,
	})

	return appRuleRead(ctx, d, m)
}

// appRuleDelete takes a pointer to the ResourceData Struct and a HTTP client and
// makes the DELETE request to OneLogin to delete an App Rule
func appRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	ruleIDStr := d.Id()
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryDelete,
			resourceTypeAppRule,
			"parsing rule_id",
			ruleIDStr,
			err,
		)
	}

	appIDStr := d.Get("app_id").(string)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		return utils.LogAndReturnError(
			ctx,
			utils.ErrorSeverityError,
			utils.ErrorCategoryDelete,
			resourceTypeAppRule,
			"parsing app_id",
			appIDStr,
			err,
		)
	}

	tflog.Info(ctx, "[DELETE] Deleting app rule", map[string]interface{}{
		"app_id":  appID,
		"rule_id": ruleID,
	})

	_, err = client.DeleteAppRule(appID, ruleID, nil)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryDelete, resourceTypeAppRule, ruleIDStr)
	}

	tflog.Info(ctx, "[DELETED] Deleted app rule", map[string]interface{}{
		"app_id":  appID,
		"rule_id": ruleID,
	})

	d.SetId("")
	return nil
}
