package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	userschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// UsersSimplified returns a user resource with CRUD methods and the appropriate schemas
func UsersSimplified() *schema.Resource {
	return &schema.Resource{
		CreateContext: userCreateSimplified,
		ReadContext:   userReadSimplified,
		UpdateContext: userUpdateSimplified,
		DeleteContext: userDeleteSimplified,
		Importer:      &schema.ResourceImporter{},
		Schema:        userschema.Schema(),
	}
}

// userCreateSimplified creates a new user in OneLogin
func userCreateSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user, err := userschema.Inflate(map[string]interface{}{
		"username":           d.Get("username"),
		"email":              d.Get("email"),
		"firstname":          d.Get("firstname"),
		"lastname":           d.Get("lastname"),
		"title":              d.Get("title"),
		"department":         d.Get("department"),
		"company":            d.Get("company"),
		"directory_id":       d.Get("directory_id"),
		"distinguished_name": d.Get("distinguished_name"),
		"external_id":        d.Get("external_id"),
		"manager_ad_id":      d.Get("manager_ad_id"),
		"manager_user_id":    d.Get("manager_user_id"),
		"member_of":          d.Get("member_of"),
		"phone":              d.Get("phone"),
		"samaccountname":     d.Get("samaccountname"),
		"userprincipalname":  d.Get("userprincipalname"),
		"state":              d.Get("state"),
		"status":             d.Get("status"),
		"group_id":           d.Get("group_id"),
		"role_ids":           d.Get("role_ids"),
		"custom_attributes":  d.Get("custom_attributes"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryCreate, "User", "")
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[CREATE] Creating user", map[string]interface{}{
		"username": d.Get("username").(string),
	})

	result, err := client.CreateUser(user)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "User", "")
	}

	// Extract user ID from the result
	userMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse user creation response")
	}

	id, ok := userMap["id"].(float64)
	if !ok {
		return diag.Errorf("failed to extract user ID from response")
	}

	userID := int(id)
	tflog.Info(ctx, "[CREATED] Created user", map[string]interface{}{
		"id":       userID,
		"username": d.Get("username").(string),
	})

	d.SetId(fmt.Sprintf("%d", userID))
	return userReadSimplified(ctx, d, m)
}

// userReadSimplified gets a user by ID from OneLogin
func userReadSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)
	uid, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "[READ] Reading user", map[string]interface{}{
		"id": uid,
	})

	result, err := client.GetUserByID(uid, &userschema.UserQueryable{})
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "User", d.Id())
	}

	// Check if user exists
	if result == nil {
		tflog.Info(ctx, "[NOT FOUND] User not found", map[string]interface{}{
			"id": uid,
		})
		d.SetId("")
		return nil
	}

	// Parse the user from the result
	userMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse user response")
	}

	// Set basic user fields
	basicFields := []string{
		"username", "email", "firstname", "lastname", "title",
		"department", "company", "status", "state", "phone",
		"group_id", "directory_id", "distinguished_name", "external_id",
		"manager_ad_id", "manager_user_id", "samaccountname", "userprincipalname",
		"member_of", "created_at", "updated_at", "activated_at", "last_login",
	}
	utils.SetResourceFields(d, userMap, basicFields)

	// Handle custom attributes if they exist
	if v, ok := userMap["custom_attributes"]; ok {
		if attrs, ok := v.(map[string]interface{}); ok {
			d.Set("custom_attributes", attrs)
		}
	}

	// Handle role IDs if they exist
	if v, ok := userMap["role_ids"]; ok {
		if roleIDs, ok := v.([]interface{}); ok {
			d.Set("role_ids", roleIDs)
		}
	}

	return nil
}

// userUpdateSimplified updates a user by ID in OneLogin
func userUpdateSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uid, _ := strconv.Atoi(d.Id())

	user, err := userschema.Inflate(map[string]interface{}{
		"id":                 d.Id(),
		"username":           d.Get("username"),
		"email":              d.Get("email"),
		"firstname":          d.Get("firstname"),
		"lastname":           d.Get("lastname"),
		"title":              d.Get("title"),
		"department":         d.Get("department"),
		"company":            d.Get("company"),
		"directory_id":       d.Get("directory_id"),
		"distinguished_name": d.Get("distinguished_name"),
		"external_id":        d.Get("external_id"),
		"manager_ad_id":      d.Get("manager_ad_id"),
		"manager_user_id":    d.Get("manager_user_id"),
		"member_of":          d.Get("member_of"),
		"phone":              d.Get("phone"),
		"samaccountname":     d.Get("samaccountname"),
		"userprincipalname":  d.Get("userprincipalname"),
		"state":              d.Get("state"),
		"status":             d.Get("status"),
		"group_id":           d.Get("group_id"),
		"role_ids":           d.Get("role_ids"),
		"custom_attributes":  d.Get("custom_attributes"),
	})
	if err != nil {
		return utils.HandleSchemaError(ctx, err, utils.ErrorCategoryUpdate, "User", d.Id())
	}

	client := m.(*onelogin.OneloginSDK)
	tflog.Info(ctx, "[UPDATE] Updating user", map[string]interface{}{
		"id": uid,
	})

	_, err = client.UpdateUser(uid, user)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "User", d.Id())
	}

	tflog.Info(ctx, "[UPDATED] Updated user", map[string]interface{}{
		"id": uid,
	})

	return userReadSimplified(ctx, d, m)
}

// userDeleteSimplified deletes a user by ID from OneLogin
func userDeleteSimplified(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		uid, _ := strconv.Atoi(id)
		return client.DeleteUser(uid)
	}, "User")
}
