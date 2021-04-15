package onelogin

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	userschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user"
)

// Users returns a resource with the CRUD methods and Terraform Schema defined
func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceUsersRead,
		Schema: userschema.QuerySchema(),
	}
}

func dataSourceUsersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	// until this fixed: https://github.com/onelogin/onelogin-go-sdk/pull/29 we can't use proper query
	// we're querying all users and looping through to find the one with matching username only
	// query, _ := userschema.QueryInflate(map[string]interface{}{
	// 	"username":            d.Get("username"),
	// })
	username := d.Get("username")
	users, err := client.Services.UsersV2.Query(nil)

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
	if len(users) == 0 {
		log.Printf("[WARNING] No users returned by the query")
		d.SetId("")
		return nil
	}

	log.Printf("[READ] %d user returned", len(users))
	log.Printf("[READ] looking for %s", username)

	for _, user := range users {
		if user.Username != nil {
			log.Printf("%+v\n", *(user.Username))
			if *(user.Username) == username {
				log.Printf("[READ] found it: %d", *(user.ID))
				d.SetId(fmt.Sprintf("%d", *(user.ID)))
				d.Set("username", *(user.Username))
				if user.Email != nil {
					d.Set("email", *(user.Email))
				}
				if user.Firstname != nil {
					d.Set("firstname", *(user.Firstname))
				}
				if user.Lastname != nil {
					d.Set("lastname", *(user.Lastname))
				}
				if user.DistinguishedName != nil {
					d.Set("distinguished_name", *(user.DistinguishedName))
				}
				if user.Samaccountname != nil {
					d.Set("samaccountname", *(user.Samaccountname))
				}
				if user.UserPrincipalName != nil {
					d.Set("user_principal_name", *(user.UserPrincipalName))
				}
				if user.MemberOf != nil {
					d.Set("member_of", *(user.MemberOf))
				}
				if user.Phone != nil {
					d.Set("phone", *(user.Phone))
				}
				if user.Title != nil {
					d.Set("title", *(user.Title))
				}
				if user.Company != nil {
					d.Set("company", *(user.Company))
				}
				if user.Department != nil {
					d.Set("department", *(user.Department))
				}
				if user.Comment != nil {
					d.Set("comment", *(user.Comment))
				}
				if user.State != nil {
					d.Set("state", *(user.State))
				}
				if user.Status != nil {
					d.Set("status", *(user.Status))
				}
				if user.GroupID != nil {
					d.Set("group_id", *(user.GroupID))
				}
				if user.DirectoryID != nil {
					d.Set("directory_id", *(user.DirectoryID))
				}
				if user.TrustedIDPID != nil {
					d.Set("trusted_idp_id", *(user.TrustedIDPID))
				}
				if user.ManagerADID != nil {
					d.Set("manager_ad_id", *(user.ManagerADID))
				}
				if user.ManagerUserID != nil {
					d.Set("manager_user_id", *(user.ManagerUserID))
				}
				if user.ExternalID != nil {
					d.Set("external_id", *(user.ExternalID))
				}
				d.Set("custom_attributes", user.CustomAttributes)
			}
		}
	}

	return nil
}
