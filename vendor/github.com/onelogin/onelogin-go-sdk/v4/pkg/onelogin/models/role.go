package models

import "encoding/json"

// RoleQuery represents available query parameters
type RoleQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// initSlice ensures that a nil slice is initialized as an empty slice
// This helper is needed for the Role.MarshalJSON method
func initSlice(s []int32) []int32 {
	if s == nil {
		return []int32{}
	}
	return s
}

// Role represents the Role resource in OneLogin
type Role struct {
	ID     *int32  `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Admins []int32 `json:"admins"`
	Apps   []int32 `json:"apps"`
	Users  []int32 `json:"users"`
}

// MarshalJSON provides custom JSON marshaling to ensure empty arrays are included in the JSON output
func (r *Role) MarshalJSON() ([]byte, error) {
	// Create a map to hold the serialized fields
	m := make(map[string]interface{})
	
	// Add ID and Name if they are not nil
	if r.ID != nil {
		m["id"] = *r.ID
	}
	if r.Name != nil {
		m["name"] = *r.Name
	}
	
	// Always include the arrays, even if they're nil (as empty arrays)
	m["admins"] = initSlice(r.Admins)
	m["apps"] = initSlice(r.Apps)
	m["users"] = initSlice(r.Users)
	
	return json.Marshal(m)
}

func (r *Role) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":  validateString,
		"page":   validateString,
		"cursor": validateString,
	}
}
