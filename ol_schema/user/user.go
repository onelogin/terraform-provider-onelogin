package userschema

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// UserQueryable implements the Queryable interface for user queries
type UserQueryable struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// GetKeyValidators returns the validation functions for the query keys
func (u *UserQueryable) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":  validateString,
		"page":   validateString,
		"cursor": validateString,
	}
}

// validateString ensures a value is a string
func validateString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

// Schema returns a key/value map of the various fields that make up a OneLogin User.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"email": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"firstname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"lastname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"distinguished_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"samaccountname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"userprincipalname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"member_of": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"phone": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"company": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"department": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"comment": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"status": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"directory_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"trusted_idp_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"manager_ad_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"manager_user_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"external_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"custom_attributes": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Map of custom attribute key/value pairs. This field is being deprecated in favor of the onelogin_user_custom_attributes resource.",
		},
	}
}

// Inflate takes a map representation of a User and returns a User object
func Inflate(s map[string]interface{}) (models.User, error) {
	// In v4 SDK, fields are directly assigned without wrapper functions
	var userID int32
	out := models.User{
		Username: s["username"].(string),
		Email:    s["email"].(string),
	}

	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			userID = int32(id)
			out.ID = userID
		}
	}

	if state, notNil := s["state"].(int); state != 0 && notNil {
		stateInt32 := int32(state)
		out.State = stateInt32
	}

	if status, notNil := s["status"].(int); status != 0 && notNil {
		statusInt32 := int32(status)
		out.Status = statusInt32
	}

	if groupid, notNil := s["group_id"].(int); groupid != 0 && notNil {
		groupInt32 := int32(groupid)
		out.GroupID = groupInt32
	}

	if directoryid, notNil := s["directory_id"].(int); directoryid != 0 && notNil {
		dirInt32 := int32(directoryid)
		out.DirectoryID = dirInt32
	}

	if trustedidpid, notNil := s["trusted_idp_id"].(int); trustedidpid != 0 && notNil {
		idpInt32 := int32(trustedidpid)
		out.TrustedIDPID = idpInt32
	}

	if manageradid, notNil := s["manager_ad_id"].(int); manageradid != 0 && notNil {
		managerAdInt32 := int32(manageradid)
		out.ManagerADID = managerAdInt32
	}

	if manageruserid, notNil := s["manager_user_id"].(int); manageruserid != 0 && notNil {
		managerUserInt32 := int32(manageruserid)
		out.ManagerUserID = managerUserInt32
	}

	if externalid, notNil := s["external_id"].(int); externalid != 0 && notNil {
		out.ExternalID = strconv.Itoa(externalid)
	}

	if firstname, notNil := s["firstname"].(string); notNil {
		out.Firstname = firstname
	}

	if lastname, notNil := s["lastname"].(string); notNil {
		out.Lastname = lastname
	}

	if distinguishedname, notNil := s["distinguished_name"].(string); notNil {
		out.DistinguishedName = distinguishedname
	}

	if samaccountname, notNil := s["samaccountname"].(string); notNil {
		out.Samaccountname = samaccountname
	}

	if userprincipalname, notNil := s["user_principal_name"].(string); notNil {
		out.UserPrincipalName = userprincipalname
	}

	if memberof, notNil := s["member_of"].(string); notNil {
		out.MemberOf = []string{memberof}
	}

	if phone, notNil := s["phone"].(string); notNil {
		out.Phone = phone
	}

	if title, notNil := s["title"].(string); notNil {
		out.Title = title
	}

	if company, notNil := s["company"].(string); notNil {
		out.Company = company
	}

	if department, notNil := s["department"].(string); notNil {
		out.Department = department
	}

	if comment, notNil := s["comment"].(string); notNil {
		out.Comment = comment
	}

	if custom_attributes, notNil := s["custom_attributes"].(map[string]interface{}); notNil {
		out.CustomAttributes = custom_attributes
	}

	return out, nil
}

func QuerySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"firstname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"lastname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"samaccountname": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"directory_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"external_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"users": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "id",
						Type:        schema.TypeInt,
						Computed:    true,
					},
					"username": {
						Description: "username",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"email": {
						Description: "email",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"firstname": {
						Description: "firstname",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"lastname": {
						Description: "lastname",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"samaccountname": {
						Description: "samaccountname",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"external_id": {
						Description: "external_id",
						Type:        schema.TypeInt,
						Computed:    true,
					},
					"directory_id": {
						Description: "directory_id",
						Type:        schema.TypeInt,
						Computed:    true,
					},
					"last_login": {
						Description: "last_login",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

// UserQuery represents the query parameters for searching users
// In v4 SDK, we need to define this since it's different from the previous SDK structure
type UserQuery struct {
	UserIDs        string
	Username       string
	DirectoryID    string
	ExternalID     string
	Firstname      string
	Lastname       string
	Samaccountname string
}

func QueryInflate(s map[string]interface{}) (UserQuery, error) {
	out := UserQuery{}
	if userid, notNil := s["user_id"].(string); notNil {
		out.UserIDs = fmt.Sprint(userid)
	}
	if username, notNil := s["username"].(string); notNil {
		out.Username = username
	}
	if directoryid, notNil := s["directory_id"].(int); directoryid != 0 && notNil {
		out.DirectoryID = fmt.Sprint(directoryid)
	}
	if externalid, notNil := s["external_id"].(int); externalid != 0 && notNil {
		out.ExternalID = fmt.Sprint(externalid)
	}
	if firstname, notNil := s["firstname"].(string); notNil {
		out.Firstname = firstname
	}
	if lastname, notNil := s["lastname"].(string); notNil {
		out.Lastname = lastname
	}
	if samaccountname, notNil := s["samaccountname"].(string); notNil {
		out.Samaccountname = samaccountname
	}

	return out, nil
}

func ReadSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"email": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"firstname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"lastname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"distinguished_name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"samaccountname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"userprincipalname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"member_of": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"phone": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"company": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"department": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"comment": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"directory_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"trusted_idp_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"manager_ad_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"manager_user_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"external_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"custom_attributes": &schema.Schema{
			Type:        schema.TypeMap,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Map of custom attribute key/value pairs. This field is being deprecated in favor of the onelogin_user_custom_attributes resource.",
		},
	}
}
