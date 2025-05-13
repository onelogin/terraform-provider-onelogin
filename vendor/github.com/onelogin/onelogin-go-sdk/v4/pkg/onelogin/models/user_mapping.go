package models

// UserMappingsQuery represents available query parameters for mappings
type UserMappingsQuery struct {
	Limit            string `json:"limit,omitempty"`
	Page             string `json:"page,omitempty"`
	Cursor           string `json:"cursor,omitempty"`
	HasCondition     string `json:"has_condition,omitempty"`
	HasConditionType string `json:"has_condition_type,omitempty"`
	HasAction        string `json:"has_action,omitempty"`
	HasActionType    string `json:"has_action_type,omitempty"`
	Enabled          string `json:"enabled,omitempty"`
}

// GetKeyValidators implements the Queryable interface
func (q *UserMappingsQuery) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":              validateString,
		"page":               validateString,
		"cursor":             validateString,
		"has_condition":      validateString,
		"has_condition_type": validateString,
		"has_action":         validateString,
		"has_action_type":    validateString,
		"enabled":            validateString,
	}
}

// UserMapping is the contract for User Mappings.
type UserMapping struct {
	ID         *int32                  `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Match      *string                 `json:"match,omitempty"`
	Enabled    *bool                   `json:"enabled,omitempty"`
	Position   *int32                  `json:"position,omitempty"`
	Conditions []UserMappingConditions `json:"conditions"`
	Actions    []UserMappingActions    `json:"actions"`
}

// UserMappingConditions is the contract for User Mapping Conditions.
type UserMappingConditions struct {
	Source   *string `json:"source,omitempty"`
	Operator *string `json:"operator,omitempty"`
	Value    *string `json:"value,omitempty"`
}

// UserMappingActions is the contract for User Mapping Actions.
type UserMappingActions struct {
	Action *string  `json:"action,omitempty"`
	Value  []string `json:"value,omitempty"`
}

func (u *UserMapping) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":            validateString,
		"page":             validateString,
		"cursor":           validateString,
		"has_condition":    validateString,
		"has_condition_id": validateString,
		"has_action":       validateString,
		"has_action_id":    validateString,
		"enabled":          validateBool,
	}
}
