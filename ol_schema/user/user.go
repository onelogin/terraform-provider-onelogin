package userschema

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/users"
)

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
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

// Inflate takes a map representation of a User and returns a User object
func Inflate(s map[string]interface{}) (users.User, error) {
	out := users.User{
		Username: oltypes.String(s["username"].(string)),
		Email:    oltypes.String(s["email"].(string)),
	}
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			out.ID = oltypes.Int32(int32(id))
		}
	}
	if state, notNil := s["state"].(int); state != 0 && notNil {
		out.State = oltypes.Int32(int32(state))
	}
	if status, notNil := s["status"].(int); status != 0 && notNil {
		out.Status = oltypes.Int32(int32(status))
	}
	if groupid, notNil := s["group_id"].(int); groupid != 0 && notNil {
		out.GroupID = oltypes.Int32(int32(groupid))
	}
	if directoryid, notNil := s["directory_id"].(int); directoryid != 0 && notNil {
		out.DirectoryID = oltypes.Int32(int32(directoryid))
	}
	if trustedidpid, notNil := s["trusted_idp_id"].(int); trustedidpid != 0 && notNil {
		out.TrustedIDPID = oltypes.Int32(int32(trustedidpid))
	}
	if manageradid, notNil := s["manager_ad_id"].(int); manageradid != 0 && notNil {
		out.ManagerADID = oltypes.Int32(int32(manageradid))
	}
	if manageruserid, notNil := s["manager_user_id"].(int); manageruserid != 0 && notNil {
		out.ManagerUserID = oltypes.Int32(int32(manageruserid))
	}
	if externalid, notNil := s["external_id"].(int); externalid != 0 && notNil {
		out.ExternalID = oltypes.Int32(int32(externalid))
	}
	if firstname, notNil := s["firstname"].(string); notNil {
		out.Firstname = oltypes.String(firstname)
	}
	if lastname, notNil := s["lastname"].(string); notNil {
		out.Lastname = oltypes.String(lastname)
	}
	if distinguishedname, notNil := s["distinguished_name"].(string); notNil {
		out.DistinguishedName = oltypes.String(distinguishedname)
	}
	if samaccountname, notNil := s["samaccountname"].(string); notNil {
		out.Samaccountname = oltypes.String(samaccountname)
	}
	if userprincipalname, notNil := s["userprincipalname"].(string); notNil {
		out.UserPrincipalName = oltypes.String(userprincipalname)
	}
	if memberof, notNil := s["member_of"].(string); notNil {
		out.MemberOf = oltypes.String(memberof)
	}
	if phone, notNil := s["phone"].(string); notNil {
		out.Phone = oltypes.String(phone)
	}
	if title, notNil := s["title"].(string); notNil {
		out.Title = oltypes.String(title)
	}
	if company, notNil := s["company"].(string); notNil {
		out.Company = oltypes.String(company)
	}
	if department, notNil := s["department"].(string); notNil {
		out.Department = oltypes.String(department)
	}
	if comment, notNil := s["comment"].(string); notNil {
		out.Comment = oltypes.String(comment)
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
	}
}

func QueryInflate(s map[string]interface{}) (users.UserQuery, error) {
	out := users.UserQuery{}
	if userid, notNil := s["user_id"].(string); notNil {
		out.UserIDs = oltypes.String(fmt.Sprint(userid))
	}
	if username, notNil := s["username"].(string); notNil {
		out.Username = oltypes.String(username)
	}
	if directoryid, notNil := s["directory_id"].(int); directoryid != 0 && notNil {
		out.DirectoryID = oltypes.String(fmt.Sprint(directoryid))
	}
	if externalid, notNil := s["external_id"].(int); externalid != 0 && notNil {
		out.ExternalID = oltypes.String(fmt.Sprint(externalid))
	}
	if firstname, notNil := s["firstname"].(string); notNil {
		out.Firstname = oltypes.String(firstname)
	}
	if lastname, notNil := s["lastname"].(string); notNil {
		out.Lastname = oltypes.String(lastname)
	}
	if samaccountname, notNil := s["samaccountname"].(string); notNil {
		out.Samaccountname = oltypes.String(samaccountname)
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
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
