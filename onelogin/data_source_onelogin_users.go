package onelogin

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	users "github.com/onelogin/onelogin-go-sdk/pkg/services/users"
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
	query, _ := userschema.QueryInflate(map[string]interface{}{
		"username":       d.Get("username"),
		"email":          d.Get("email"),
		"firstname":      d.Get("firstname"),
		"lastname":       d.Get("lastname"),
		"samaccountname": d.Get("samaccountname"),
		"external_id":    d.Get("external_id"),
		"directory_id":   d.Get("directory_id"),
	})
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
	if len(users) == 0 {
		log.Printf("[WARNING] No users returned by the query")
		d.SetId("")
		return nil
	}

	log.Printf("[READ] %d user returned", len(users))

	userIds := make([]string, 0)
	for _, user := range users {
		userIds = append(userIds, fmt.Sprintf("%d", *(user.ID)))
	}

	d.SetId(fmt.Sprintf("%d", HashQuery(&query)))
	d.Set("ids", userIds)

	return nil
}

func HashQuery(query *users.UserQuery) [16]byte {
	bytes, _ := json.Marshal(query)
	return md5.Sum(bytes)
}
