package onelogin

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/terraform-provider-onelogin/ol_schema/user"
)

// Apps returns a resource with the CRUD methods and Terraform Schema defined
func Users() *schema.Resource {
	return &schema.Resource{
		Create:   usersCreate,
		Read:     usersRead,
		Update:   usersUpdate,
		Delete:   usersDelete,
		Importer: &schema.ResourceImporter{},
		Schema:   userschema.Schema(),
	}
}

func usersCreate(d *schema.ResourceData, m interface{}) error {
	user, _ := userschema.Inflate(map[string]interface{}{
		"username":            d.Get("username"),
		"email":               d.Get("email"),
		"firstname":           d.Get("firstname"),
		"lastname":            d.Get("lastname"),
		"distinguished_name":  d.Get("distinguished_name"),
		"samaccountname":      d.Get("samaccountname"),
		"user_principal_name": d.Get("user_principal_name"),
		"member_of":           d.Get("member_of"),
		"phone":               d.Get("phone"),
		"title":               d.Get("title"),
		"company":             d.Get("company"),
		"department":          d.Get("department"),
		"comment":             d.Get("comment"),
		"state":               d.Get("state"),
		"status":              d.Get("status"),
		"group_id":            d.Get("group_id"),
		"directory_id":        d.Get("directory_id"),
		"trusted_idp_id":      d.Get("trusted_idp_id"),
		"manager_ad_id":       d.Get("manager_ad_id"),
		"manager_user_id":     d.Get("manager_user_id"),
		"external_id":         d.Get("external_id"),
	})
	client := m.(*client.APIClient)
	err := client.Services.UsersV2.Create(&user)
	if err != nil {
		log.Println("[ERROR] There was a problem creating the user!", err)
		return err
	}
	log.Printf("[CREATED] Created user with %d", *(user.ID))

	d.SetId(fmt.Sprintf("%d", *(user.ID)))
	return usersRead(d, m)
}

func usersUpdate(d *schema.ResourceData, m interface{}) error {
	uid, _ := strconv.Atoi(d.Id())
	user, _ := userschema.Inflate(map[string]interface{}{
		"username":            d.Get("username"),
		"email":               d.Get("email"),
		"firstname":           d.Get("firstname"),
		"lastname":            d.Get("lastname"),
		"distinguished_name":  d.Get("distinguished_name"),
		"samaccountname":      d.Get("samaccountname"),
		"user_principal_name": d.Get("user_principal_name"),
		"member_of":           d.Get("member_of"),
		"phone":               d.Get("phone"),
		"title":               d.Get("title"),
		"company":             d.Get("company"),
		"department":          d.Get("department"),
		"comment":             d.Get("comment"),
		"state":               d.Get("state"),
		"status":              d.Get("status"),
		"group_id":            d.Get("group_id"),
		"directory_id":        d.Get("directory_id"),
		"trusted_idp_id":      d.Get("trusted_idp_id"),
		"manager_ad_id":       d.Get("manager_ad_id"),
		"manager_user_id":     d.Get("manager_user_id"),
		"external_id":         d.Get("external_id"),
	})
	client := m.(*client.APIClient)
	err := client.Services.UsersV2.Update(int32(uid), &user)
	if err != nil {
		log.Println("[ERROR] There was a problem updating the user!", err)
		return err
	}
	log.Printf("[CREATED] Updated user with %d", *(user.ID))

	d.SetId(fmt.Sprintf("%d", *(user.ID)))
	return usersRead(d, m)
}

func usersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.APIClient)
	uid, _ := strconv.Atoi(d.Id())
	user, err := client.Services.UsersV2.GetOne(int32(uid))
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the user!")
		log.Println(err)
		return err
	}
	if user == nil {
		d.SetId("")
		return nil
	}
	log.Printf("[READ] Reading user with %d", *(user.ID))

	d.Set("username", user.Username)
	d.Set("email", user.Email)
	d.Set("firstname", user.Firstname)
	d.Set("lastname", user.Lastname)
	d.Set("distinguished_name", user.DistinguishedName)
	d.Set("samaccountname", user.Samaccountname)
	d.Set("user_principal_name", user.UserPrincipalName)
	d.Set("member_of", user.MemberOf)
	d.Set("phone", user.Phone)
	d.Set("title", user.Title)
	d.Set("company", user.Company)
	d.Set("department", user.Department)
	d.Set("comment", user.Comment)
	d.Set("state", user.State)
	d.Set("status", user.Status)
	d.Set("group_id", user.GroupID)
	d.Set("directory_id", user.DirectoryID)
	d.Set("trusted_idp_id", user.TrustedIDPID)
	d.Set("manager_ad_id", user.ManagerADID)
	d.Set("manager_user_id", user.ManagerUserID)
	d.Set("external_id", user.ExternalID)

	return nil
}

func usersDelete(d *schema.ResourceData, m interface{}) error {
	uid, _ := strconv.Atoi(d.Id())
	client := m.(*client.APIClient)

	err := client.Services.UsersV2.Destroy(int32(uid))
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the user!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted user with %d", uid)
		d.SetId("")
	}

	return nil
}
