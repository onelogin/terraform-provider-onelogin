package models

// PrivilegeQuery represents available query parameters
type PrivilegeQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// Privilege represents the Role resource in OneLogin
type Privilege struct {
	ID          *string        `json:"id,omitempty"`
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Privilege   *PrivilegeData `json:"privilege,omitempty"`
	UserIDs     []int          `json:"user_ids,omitempty"`
	RoleIDs     []int          `json:"role_ids,omitempty"`
}

// PrivilegeData represents the group of statements and statement versions pertinent to a privilege
type PrivilegeData struct {
	Version   *string         `json:"version,omitempty"`
	Statement []StatementData `json:"Statement"`
}

// StatementData represents the actions and scope of a given privilege
type StatementData struct {
	Effect *string  `json:"Effect,omitempty"`
	Action []string `json:"Action"`
	Scope  []string `json:"Scope"`
}

func (p *Privilege) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":  validateString,
		"page":   validateString,
		"cursor": validateString,
	}
}
