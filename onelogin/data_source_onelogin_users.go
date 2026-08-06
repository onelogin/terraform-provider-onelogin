package onelogin

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
	client := m.(*onelogin.OneloginSDK)
	query, _ := userschema.QueryInflate(map[string]interface{}{
		"username":       d.Get("username"),
		"email":          d.Get("email"),
		"firstname":      d.Get("firstname"),
		"lastname":       d.Get("lastname"),
		"samaccountname": d.Get("samaccountname"),
		"external_id":    d.Get("external_id"),
		"directory_id":   d.Get("directory_id"),
	})

	// Create pointers for non-empty query parameters
	var username, email, firstname, lastname, samaccountname, externalID, directoryID *string

	if query.Username != "" {
		usernameVal := query.Username
		username = &usernameVal
	}

	if query.Email != "" {
		emailVal := query.Email
		email = &emailVal
	}

	if query.Firstname != "" {
		firstnameVal := query.Firstname
		firstname = &firstnameVal
	}

	if query.Lastname != "" {
		lastnameVal := query.Lastname
		lastname = &lastnameVal
	}

	if query.Samaccountname != "" {
		samVal := query.Samaccountname
		samaccountname = &samVal
	}

	if query.ExternalID != "" {
		extIDVal := query.ExternalID
		externalID = &extIDVal
	}

	if query.DirectoryID != "" {
		dirIDVal := query.DirectoryID
		directoryID = &dirIDVal
	}

	// Create the SDK query
	sdkQuery := &models.UserQuery{
		Username:       username,
		Email:          email,
		Firstname:      firstname,
		Lastname:       lastname,
		Samaccountname: samaccountname,
		ExternalID:     externalID,
		DirectoryID:    directoryID,
	}

	// The API matches a single email per request, so a list of them becomes one
	// request each and the results are combined. Every other filter still
	// applies to each, which keeps "these people, in this directory" expressible.
	queries := usersQueriesForEmails(sdkQuery, emailsFromConfig(d))

	data := make([]interface{}, 0)
	seen := make(map[int]bool)
	for _, q := range queries {
		result, err := client.GetUsers(q)
		if err != nil {
			log.Printf("[ERROR] There was a problem reading the users!")
			log.Println(err)
			return err
		}

		// Parse the response
		found, err := usersFromResponse(result)
		if err != nil {
			log.Printf("[WARNING] Invalid response format")
			d.SetId("")
			return err
		}

		// A user matching two of the emails, or two filters, is still one user.
		for _, userData := range found {
			user, ok := userData.(map[string]interface{})
			if !ok {
				continue
			}
			id, ok := user["id"].(float64)
			if !ok {
				continue
			}
			if seen[int(id)] {
				continue
			}
			seen[int(id)] = true
			data = append(data, userData)
		}
	}

	if len(data) == 0 {
		log.Printf("[WARNING] No users returned by the query")
		d.SetId("")
		return nil
	}

	log.Printf("[READ] %d users returned", len(data))

	userIds := make([]string, 0)
	userList := make([]map[string]interface{}, 0)
	for _, userData := range data {
		user, ok := userData.(map[string]interface{})
		if !ok {
			continue
		}

		// Get the user ID
		userID, ok := user["id"].(float64)
		if !ok {
			continue
		}

		userIds = append(userIds, fmt.Sprintf("%d", int(userID)))

		u := make(map[string]interface{})
		u["id"] = int(userID)

		if v, ok := user["username"]; ok {
			u["username"] = v
		}
		if v, ok := user["email"]; ok {
			u["email"] = v
		}
		if v, ok := user["firstname"]; ok {
			u["firstname"] = v
		}
		if v, ok := user["lastname"]; ok {
			u["lastname"] = v
		}
		if v, ok := user["samaccountname"]; ok {
			u["samaccountname"] = v
		}
		if v, ok := user["external_id"]; ok {
			u["external_id"] = v
		}
		if v, ok := user["directory_id"]; ok {
			u["directory_id"] = v
		}
		if v, ok := user["last_login"]; ok {
			// Handle last_login which might be a string or time value
			var lastLoginStr string
			switch lastLogin := v.(type) {
			case string:
				lastLoginStr = lastLogin
			case time.Time:
				lastLoginStr = lastLogin.Format(time.RFC3339)
			}

			// Only set last_login if it's not a "never logged in" placeholder date
			if lastLoginStr != "" && !isNeverLoggedInDate(lastLoginStr) {
				u["last_login"] = lastLoginStr
			}
			// If it is a placeholder date, don't set it (leave it empty/null)
		}

		userList = append(userList, u)
	}

	// Generate a hash for the ID. The emails go in alongside the query: they are
	// not part of it, so two different email lists would otherwise hash alike.
	queryBytes, _ := json.Marshal(struct {
		Query  *models.UserQuery
		Emails []string
	}{sdkQuery, emailsFromConfig(d)})
	queryHash := md5.Sum(queryBytes)
	d.SetId(fmt.Sprintf("%x", queryHash))

	d.Set("ids", userIds)
	d.Set("users", userList)

	return nil
}

// emailsFromConfig reads the "emails" filter, dropping blanks so an interpolated
// value that came out empty does not turn into a query for every user.
func emailsFromConfig(d *schema.ResourceData) []string {
	raw, ok := d.Get("emails").([]interface{})
	if !ok {
		return nil
	}

	emails := make([]string, 0, len(raw))
	for _, v := range raw {
		if email, ok := v.(string); ok && email != "" {
			emails = append(emails, email)
		}
	}
	return emails
}

// usersQueriesForEmails fans a query out over a list of emails, one query each,
// leaving every other filter in place. With no emails the query is returned
// unchanged, so the single-email and unfiltered cases are untouched.
func usersQueriesForEmails(base *models.UserQuery, emails []string) []*models.UserQuery {
	if len(emails) == 0 {
		return []*models.UserQuery{base}
	}

	queries := make([]*models.UserQuery, 0, len(emails))
	for _, email := range emails {
		q := *base
		emailVal := email
		q.Email = &emailVal
		queries = append(queries, &q)
	}
	return queries
}

func HashQuery(query *userschema.UserQuery) [16]byte {
	bytes, _ := json.Marshal(query)
	return md5.Sum(bytes)
}
