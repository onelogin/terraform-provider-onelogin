package models

import "encoding/json"

// RoleQuery represents available query parameters
type RoleQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// Role represents the Role resource in OneLogin
type Role struct {
	ID     *int32  `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Admins []int32 `json:"admins"`
	Apps   []int32 `json:"apps"`
	Users  []int32 `json:"users"`
}

// MarshalJSON keeps a nil membership slice out of the request body entirely,
// while still sending an explicitly empty one.
//
// The API replaces memberships wholesale for any array it receives, so the two
// have to stay distinct or a caller updating one field silently wipes the rest:
//
//	nil            omit the key; leave those memberships alone
//	[]int32{}      send []; remove every membership
//	[]int32{1, 2}  send the ids; replace the memberships
//
// The struct tags cannot express this on their own. Without omitempty a nil
// slice marshals to null, and with it an explicitly empty slice disappears too,
// which is why this is hand-rolled. The receiver is deliberately a value so that
// json.Marshal treats a Role and a *Role the same way.
func (r Role) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	if r.ID != nil {
		m["id"] = *r.ID
	}
	if r.Name != nil {
		m["name"] = *r.Name
	}

	if r.Admins != nil {
		m["admins"] = r.Admins
	}
	if r.Apps != nil {
		m["apps"] = r.Apps
	}
	if r.Users != nil {
		m["users"] = r.Users
	}

	return json.Marshal(m)
}

func (r *Role) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":  validateString,
		"page":   validateString,
		"cursor": validateString,
	}
}
