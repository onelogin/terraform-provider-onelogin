package onelogin

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	userschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user"
)

// Users returns a resource with the CRUD methods and Terraform Schema defined
func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceUserRead,
		Schema: userschema.ReadSchema(),
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*onelogin.OneloginSDK)
	query, _ := userschema.QueryInflate(map[string]interface{}{
		"username": d.Get("username"),
		"user_id":  d.Get("user_id"),
	})

	if query.UserIDs == "" && query.Username == "" {
		return fmt.Errorf("At least one of either username or user_id must be defined")
	}

	// In v4 SDK, we need to use the models.UserQuery struct that implements Queryable
	// Create a pointer to a string for each field that has a value
	var username, userIDs *string

	if query.Username != "" {
		usernameVal := query.Username
		username = &usernameVal
	}

	if query.UserIDs != "" {
		userIDsVal := query.UserIDs
		userIDs = &userIDsVal
	}

	// Create the query object using the SDK models
	sdkQuery := &models.UserQuery{
		Username: username,
		UserIDs:  userIDs,
	}

	result, err := client.GetUsers(sdkQuery)
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the user!")
		log.Println(err)
		return err
	}

	// Parse the users from the result
	respMap, ok := result.(map[string]interface{})
	if !ok {
		log.Printf("[WARNING] Invalid response format")
		d.SetId("")
		return fmt.Errorf("Invalid response format from API")
	}

	data, ok := respMap["data"].([]interface{})
	if !ok || len(data) == 0 {
		log.Printf("[WARNING] No users returned by the query")
		d.SetId("")
		return nil
	}

	if len(data) != 1 {
		log.Printf("[WARNING] %d users returned by the query", len(data))
		d.SetId("")
		return fmt.Errorf("Your query returned more than one result. Usernames and IDs should be unique")
	}

	// Get the first user from the data
	user, ok := data[0].(map[string]interface{})
	if !ok {
		log.Printf("[WARNING] Invalid user format")
		d.SetId("")
		return fmt.Errorf("Invalid user format in response")
	}

	// Set the user ID
	userID, ok := user["id"].(float64)
	if !ok {
		log.Printf("[WARNING] Invalid or missing user ID")
		d.SetId("")
		return fmt.Errorf("Invalid or missing user ID in response")
	}

	d.SetId(fmt.Sprintf("%d", int(userID)))

	// Set user fields
	if v, ok := user["username"]; ok {
		d.Set("username", v)
	}
	if v, ok := user["email"]; ok {
		d.Set("email", v)
	}
	if v, ok := user["firstname"]; ok {
		d.Set("firstname", v)
	}
	if v, ok := user["lastname"]; ok {
		d.Set("lastname", v)
	}
	if v, ok := user["distinguished_name"]; ok {
		d.Set("distinguished_name", v)
	}
	if v, ok := user["samaccountname"]; ok {
		d.Set("samaccountname", v)
	}
	if v, ok := user["userprincipalname"]; ok {
		d.Set("user_principal_name", v)
	}
	if v, ok := user["member_of"]; ok {
		memberOf, isList := v.([]interface{})
		if isList && len(memberOf) > 0 {
			// If member_of is a list, use the first item
			d.Set("member_of", memberOf[0])
		} else {
			d.Set("member_of", v)
		}
	}
	if v, ok := user["phone"]; ok {
		d.Set("phone", v)
	}
	if v, ok := user["title"]; ok {
		d.Set("title", v)
	}
	if v, ok := user["company"]; ok {
		d.Set("company", v)
	}
	if v, ok := user["department"]; ok {
		d.Set("department", v)
	}
	if v, ok := user["comment"]; ok {
		d.Set("comment", v)
	}
	if v, ok := user["state"]; ok {
		d.Set("state", v)
	}
	if v, ok := user["status"]; ok {
		d.Set("status", v)
	}
	if v, ok := user["group_id"]; ok {
		d.Set("group_id", v)
	}
	if v, ok := user["directory_id"]; ok {
		d.Set("directory_id", v)
	}
	if v, ok := user["trusted_idp_id"]; ok {
		d.Set("trusted_idp_id", v)
	}
	if v, ok := user["manager_ad_id"]; ok {
		d.Set("manager_ad_id", v)
	}
	if v, ok := user["manager_user_id"]; ok {
		d.Set("manager_user_id", v)
	}
	if v, ok := user["external_id"]; ok {
		d.Set("external_id", v)
	}
	if v, ok := user["custom_attributes"]; ok {
		d.Set("custom_attributes", v)
	}

	return nil
}
