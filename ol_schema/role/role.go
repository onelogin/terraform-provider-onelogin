package roleschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// RoleQuery implements the Queryable interface for role queries
type RoleQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// GetKeyValidators returns the validation functions for the query keys
func (r *RoleQuery) GetKeyValidators() map[string]func(interface{}) bool {
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
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"apps": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
		"users": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
		"admins": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct a Role
func Inflate(s map[string]interface{}) *models.Role {
	var roleID int32
	var roleName string
	out := &models.Role{}

	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			roleID = int32(id)
			out.ID = &roleID
		}
	}

	if name, notNil := s["name"].(string); notNil {
		roleName = name
		out.Name = &roleName
	}

	if s["apps"] != nil {
		out.Apps = make([]int32, len(s["apps"].(*schema.Set).List()))
		for i, appID := range s["apps"].(*schema.Set).List() {
			out.Apps[i] = int32(appID.(int))
		}
	}

	if s["users"] != nil {
		out.Users = make([]int32, len(s["users"].(*schema.Set).List()))
		for i, userID := range s["users"].(*schema.Set).List() {
			out.Users[i] = int32(userID.(int))
		}
	}

	if s["admins"] != nil {
		out.Admins = make([]int32, len(s["admins"].(*schema.Set).List()))
		for i, adminID := range s["admins"].(*schema.Set).List() {
			out.Admins[i] = int32(adminID.(int))
		}
	}

	return out
}
