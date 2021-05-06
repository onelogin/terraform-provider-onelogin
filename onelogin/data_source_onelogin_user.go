package onelogin

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
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
	client := m.(*client.APIClient)
	query, _ := userschema.QueryInflate(map[string]interface{}{
		"username": d.Get("username"),
		"user_id":  d.Get("user_id"),
	})

	if *(query.UserIDs) == "" && *(query.Username) == "" {
		return fmt.Errorf("At least one of either username or user_id must be defined")
	}

	users, err := client.Services.UsersV2.Query(&query)

	if err != nil {
		log.Printf("[ERROR] There was a problem reading the user!")
		log.Println(err)
		return err
	}
	if users == nil {
		log.Printf("[WARNING] Nil users returned by the query")
		d.SetId("")
		return nil
	}
	if len(users) != 1 {
		log.Printf("[WARNING] %d user returned by the query", len(users))
		d.SetId("")
		return fmt.Errorf("Your query returned more than one result. Usernames and IDs should be unique")
	}

	d.SetId(fmt.Sprintf("%d", *(users[0].ID)))
	if users[0].Username != nil {
		d.Set("username", *(users[0].Username))
	}
	if users[0].Email != nil {
		d.Set("email", *(users[0].Email))
	}
	if users[0].Firstname != nil {
		d.Set("firstname", *(users[0].Firstname))
	}
	if users[0].Lastname != nil {
		d.Set("lastname", *(users[0].Lastname))
	}
	if users[0].DistinguishedName != nil {
		d.Set("distinguished_name", *(users[0].DistinguishedName))
	}
	if users[0].Samaccountname != nil {
		d.Set("samaccountname", *(users[0].Samaccountname))
	}
	if users[0].UserPrincipalName != nil {
		d.Set("user_principal_name", *(users[0].UserPrincipalName))
	}
	if users[0].MemberOf != nil {
		d.Set("member_of", *(users[0].MemberOf))
	}
	if users[0].Phone != nil {
		d.Set("phone", *(users[0].Phone))
	}
	if users[0].Title != nil {
		d.Set("title", *(users[0].Title))
	}
	if users[0].Company != nil {
		d.Set("company", *(users[0].Company))
	}
	if users[0].Department != nil {
		d.Set("department", *(users[0].Department))
	}
	if users[0].Comment != nil {
		d.Set("comment", *(users[0].Comment))
	}
	if users[0].State != nil {
		d.Set("state", *(users[0].State))
	}
	if users[0].Status != nil {
		d.Set("status", *(users[0].Status))
	}
	if users[0].GroupID != nil {
		d.Set("group_id", *(users[0].GroupID))
	}
	if users[0].DirectoryID != nil {
		d.Set("directory_id", *(users[0].DirectoryID))
	}
	if users[0].TrustedIDPID != nil {
		d.Set("trusted_idp_id", *(users[0].TrustedIDPID))
	}
	if users[0].ManagerADID != nil {
		d.Set("manager_ad_id", *(users[0].ManagerADID))
	}
	if users[0].ManagerUserID != nil {
		d.Set("manager_user_id", *(users[0].ManagerUserID))
	}
	if users[0].ExternalID != nil {
		d.Set("external_id", *(users[0].ExternalID))
	}
	d.Set("custom_attributes", users[0].CustomAttributes)

	return nil
}
