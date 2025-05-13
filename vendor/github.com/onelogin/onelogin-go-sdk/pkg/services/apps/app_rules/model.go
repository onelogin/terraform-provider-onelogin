package apprules

type AppRuleQuery struct {
	AppID string
}

// AppRule is the contract for App Rules.
type AppRule struct {
	ID         *int32              `json:"id,omitempty"`
	AppID      *int32              `json:"app_id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Match      *string             `json:"match,omitempty"`
	Enabled    *bool               `json:"enabled,omitempty"`
	Position   *int32              `json:"position,omitempty"`
	Conditions []AppRuleConditions `json:"conditions"`
	Actions    []AppRuleActions    `json:"actions"`
}

type AppRuleConditions struct {
	Source   *string `json:"source,omitempty"`
	Operator *string `json:"operator,omitempty"`
	Value    *string `json:"value,omitempty"`
}

type AppRuleActions struct {
	Action     *string  `json:"action,omitempty"`
	Value      []string `json:"value,omitempty"`
	Expression *string  `json:"expression,omitempty"`
}
