package usermappings

// UserMappingsQuery represents available query parameters for mappings
type UserMappingsQuery struct {
	Limit            string
	Page             string
	Cursor           string
	HasCondition     string
	HasConditionType string
	HasAction        string
	HasActionType    string
	Enabled          string
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
