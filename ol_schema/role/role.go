package roleschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// MembershipAttrs are the role sub-resources that live on their own paginated
// endpoints, and so are the ones worth being able to leave un-refreshed.
var MembershipAttrs = []string{"apps", "users", "admins"}

func validMembershipAttr(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOfValue(key, val, MembershipAttrs)
}

// RoleQuery implements the Queryable interface for role queries.
//
// The V2 API treats cursor and limit/page as mutually exclusive — supplying a
// cursor alongside either is rejected — so callers clear Limit and Page once
// they have a cursor from the After-Cursor response header.
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
		// Membership lives on three paginated sub-endpoints, and a role with
		// thousands of users costs a page walk on every refresh whether or not
		// the configuration manages them. ignore_changes cannot help: it is
		// applied to the diff, long after the read has made the calls.
		//
		// Naming this as what to skip rather than what to manage is deliberate.
		// A "manage_users" flag defaulting to true would read as false on an
		// imported resource, whose state carries no value for it, and silently
		// skip the very fetch import needs.
		"skip_membership_refresh": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validMembershipAttr,
			},
			Description: "Membership attributes to leave un-refreshed: any of apps, users, admins. Their sub-endpoints are not queried during read, and state keeps whatever it already held. Use for membership managed outside Terraform, alongside a lifecycle ignore_changes block for the same attribute.",
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct a Role
func Inflate(s map[string]interface{}) *models.Role {
	var roleID int32
	var roleName string
	// Initialize the role with empty arrays by default
	// Each array will only be populated if it exists in the input map
	out := &models.Role{}

	// We don't include ID in API payloads:
	// - For create: The API will generate an ID
	// - For update: The ID is already in the URL path
	//
	// The ID is only needed for the Terraform state
	if s["include_id_in_output"] == true && s["id"] != nil {
		// Handle both string and int inputs
		var id int
		var err error

		switch v := s["id"].(type) {
		case string:
			id, err = strconv.Atoi(v)
			if err == nil && id > 0 {
				roleID = int32(id)
				out.ID = &roleID
			}
		case int:
			if v > 0 {
				roleID = int32(v)
				out.ID = &roleID
			}
		}
	}

	if name, notNil := s["name"].(string); notNil {
		roleName = name
		out.Name = &roleName
	}

	// Only populate apps if provided in the input map
	if s["apps"] != nil {
		// Handle both schema.Set and []int inputs
		var appList []interface{}
		switch apps := s["apps"].(type) {
		case *schema.Set:
			appList = apps.List()
		case []int:
			appList = make([]interface{}, len(apps))
			for i, id := range apps {
				appList[i] = id
			}
		}

		if appList != nil {
			out.Apps = make([]int32, len(appList))
			for i, appID := range appList {
				out.Apps[i] = int32(appID.(int))
			}
		}
	}

	// Only populate users if provided in the input map
	if s["users"] != nil {
		// Handle both schema.Set and []int inputs
		var userList []interface{}
		switch users := s["users"].(type) {
		case *schema.Set:
			userList = users.List()
		case []int:
			userList = make([]interface{}, len(users))
			for i, id := range users {
				userList[i] = id
			}
		}

		if userList != nil {
			out.Users = make([]int32, len(userList))
			for i, userID := range userList {
				out.Users[i] = int32(userID.(int))
			}
		}
	}

	// Only populate admins if provided in the input map
	if s["admins"] != nil {
		// Handle both schema.Set and []int inputs
		var adminList []interface{}
		switch admins := s["admins"].(type) {
		case *schema.Set:
			adminList = admins.List()
		case []int:
			adminList = make([]interface{}, len(admins))
			for i, id := range admins {
				adminList[i] = id
			}
		}

		if adminList != nil {
			out.Admins = make([]int32, len(adminList))
			for i, adminID := range adminList {
				out.Admins[i] = int32(adminID.(int))
			}
		}
	}

	return out
}
