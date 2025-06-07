package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	selfregistrationprofileschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/self_registration_profile"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// SelfRegistrationProfiles returns a resource with the CRUD methods and Terraform Schema defined
func SelfRegistrationProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: selfRegistrationProfileCreate,
		ReadContext:   selfRegistrationProfileRead,
		UpdateContext: selfRegistrationProfileUpdate,
		DeleteContext: selfRegistrationProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: selfregistrationprofileschema.Schema(),
	}
}

// selfRegistrationProfileCreate creates a new self-registration profile in OneLogin
func selfRegistrationProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	profile := models.SelfRegistrationProfile{
		Name:                 d.Get("name").(string),
		URL:                  d.Get("url").(string),
		Enabled:              d.Get("enabled").(bool),
		Moderated:            d.Get("moderated").(bool),
		Helptext:             d.Get("helptext").(string),
		ThankyouMessage:      d.Get("thankyou_message").(string),
		DomainBlacklist:      d.Get("domain_blacklist").(string),
		DomainWhitelist:      d.Get("domain_whitelist").(string),
		DomainListStrategy:   int32(d.Get("domain_list_strategy").(int)),
		EmailVerificationType: d.Get("email_verification_type").(string),
	}

	if v, ok := d.GetOk("default_role_id"); ok {
		profile.DefaultRoleID = int32(v.(int))
	}

	if v, ok := d.GetOk("default_group_id"); ok {
		profile.DefaultGroupID = int32(v.(int))
	}

	tflog.Info(ctx, "[CREATE] Creating self-registration profile", map[string]interface{}{
		"name": profile.Name,
	})

	result, err := client.CreateSelfRegistrationProfile(profile)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Self-Registration Profile", "")
	}

	// Extract profile ID from the result
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse self-registration profile creation response")
	}

	profileMap, ok := resultMap["self_registration_profile"].(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to extract self-registration profile from response")
	}

	id, ok := profileMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract self-registration profile ID from response")
	}

	profileID := int(id)
	d.SetId(fmt.Sprintf("%d", profileID))

	// Handle fields if they exist
	if v, ok := d.GetOk("fields"); ok {
		fields := v.(*schema.Set).List()
		for _, field := range fields {
			fieldMap := field.(map[string]interface{})
			customAttributeID := int(fieldMap["custom_attribute_id"].(int))

			_, err := client.CreateSelfRegistrationProfileField(profileID, customAttributeID)
			if err != nil {
				return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Self-Registration Profile Field", "")
			}
		}
	}

	tflog.Info(ctx, "[CREATED] Created self-registration profile", map[string]interface{}{
		"id":   profileID,
		"name": profile.Name,
	})

	return selfRegistrationProfileRead(ctx, d, m)
}

// selfRegistrationProfileRead reads a self-registration profile from OneLogin
func selfRegistrationProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	profileID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "[READ] Reading self-registration profile", map[string]interface{}{
		"id": profileID,
	})

	result, err := client.GetSelfRegistrationProfile(profileID)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Self-Registration Profile", d.Id())
	}

	// Parse the profile from the result
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse self-registration profile response")
	}

	profileMap, ok := resultMap["self_registration_profile"].(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to extract self-registration profile from response")
	}

	// Set basic profile fields
	d.Set("name", profileMap["name"])
	d.Set("url", profileMap["url"])
	d.Set("enabled", profileMap["enabled"])
	d.Set("moderated", profileMap["moderated"])
	d.Set("helptext", profileMap["helptext"])
	d.Set("thankyou_message", profileMap["thankyou_message"])
	d.Set("domain_blacklist", profileMap["domain_blacklist"])
	d.Set("domain_whitelist", profileMap["domain_whitelist"])
	d.Set("domain_list_strategy", profileMap["domain_list_strategy"])
	d.Set("email_verification_type", profileMap["email_verification_type"])

	if v, ok := profileMap["default_role_id"]; ok && v != nil {
		d.Set("default_role_id", v)
	}

	if v, ok := profileMap["default_group_id"]; ok && v != nil {
		d.Set("default_group_id", v)
	}

	// Handle fields if they exist
	if fields, ok := profileMap["fields"].([]interface{}); ok {
		fieldSet := schema.NewSet(schema.HashResource(&schema.Resource{
			Schema: selfregistrationprofileschema.Schema()["fields"].Elem.(*schema.Resource).Schema,
		}), []interface{}{})

		for _, field := range fields {
			fieldMap := field.(map[string]interface{})
			fieldSet.Add(map[string]interface{}{
				"id":                 int(fieldMap["id"].(float64)),
				"custom_attribute_id": int(fieldMap["custom_attribute_id"].(float64)),
				"name":               fieldMap["name"],
			})
		}

		d.Set("fields", fieldSet)
	}

	return nil
}

// selfRegistrationProfileUpdate updates a self-registration profile in OneLogin
func selfRegistrationProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	profileID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	profile := models.SelfRegistrationProfile{
		Name:                 d.Get("name").(string),
		URL:                  d.Get("url").(string),
		Enabled:              d.Get("enabled").(bool),
		Moderated:            d.Get("moderated").(bool),
		Helptext:             d.Get("helptext").(string),
		ThankyouMessage:      d.Get("thankyou_message").(string),
		DomainBlacklist:      d.Get("domain_blacklist").(string),
		DomainWhitelist:      d.Get("domain_whitelist").(string),
		DomainListStrategy:   int32(d.Get("domain_list_strategy").(int)),
		EmailVerificationType: d.Get("email_verification_type").(string),
	}

	if v, ok := d.GetOk("default_role_id"); ok {
		profile.DefaultRoleID = int32(v.(int))
	}

	if v, ok := d.GetOk("default_group_id"); ok {
		profile.DefaultGroupID = int32(v.(int))
	}

	tflog.Info(ctx, "[UPDATE] Updating self-registration profile", map[string]interface{}{
		"id": profileID,
	})

	_, err = client.UpdateSelfRegistrationProfile(profileID, profile)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Self-Registration Profile", d.Id())
	}

	// Handle fields if they've changed
	if d.HasChange("fields") {
		// Get the current fields from the API
		result, err := client.GetSelfRegistrationProfile(profileID)
		if err != nil {
			return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Self-Registration Profile", d.Id())
		}

		resultMap := result.(map[string]interface{})
		profileMap := resultMap["self_registration_profile"].(map[string]interface{})
		currentFields := profileMap["fields"].([]interface{})

		// Create a map of current field IDs to their custom attribute IDs
		currentFieldMap := make(map[int]int)
		for _, field := range currentFields {
			fieldMap := field.(map[string]interface{})
			fieldID := int(fieldMap["id"].(float64))
			customAttributeID := int(fieldMap["custom_attribute_id"].(float64))
			currentFieldMap[customAttributeID] = fieldID
		}

		// Get the desired fields from the schema
		desiredFields := d.Get("fields").(*schema.Set).List()
		desiredFieldMap := make(map[int]bool)
		for _, field := range desiredFields {
			fieldMap := field.(map[string]interface{})
			customAttributeID := fieldMap["custom_attribute_id"].(int)
			desiredFieldMap[customAttributeID] = true

			// Add fields that don't exist
			if _, exists := currentFieldMap[customAttributeID]; !exists {
				_, err := client.CreateSelfRegistrationProfileField(profileID, customAttributeID)
				if err != nil {
					return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Self-Registration Profile Field", "")
				}
			}
		}

		// Remove fields that are no longer desired
		for customAttributeID, fieldID := range currentFieldMap {
			if _, exists := desiredFieldMap[customAttributeID]; !exists {
				_, err := client.DeleteSelfRegistrationProfileField(profileID, fieldID)
				if err != nil {
					return utils.HandleAPIError(ctx, err, utils.ErrorCategoryDelete, "Self-Registration Profile Field", "")
				}
			}
		}
	}

	tflog.Info(ctx, "[UPDATED] Updated self-registration profile", map[string]interface{}{
		"id": profileID,
	})

	return selfRegistrationProfileRead(ctx, d, m)
}

// selfRegistrationProfileDelete deletes a self-registration profile from OneLogin
func selfRegistrationProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		profileID, _ := strconv.Atoi(id)
		return client.DeleteSelfRegistrationProfile(profileID)
	}, "Self-Registration Profile")
}
