package models

type Condition struct {
	Source   string `json:"source"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Action struct {
	Action     string   `json:"action"`
	Value      []string `json:"value,omitempty"`
	Expression string   `json:"expression,omitempty"`
	Scriplet   string   `json:"scriplet,omitempty"`
	Macro      string   `json:"macro,omitempty"`
}

type AppRule struct {
	AppID      int         `json:"app_id"`
	Name       string      `json:"name"`
	Enabled    bool        `json:"enabled"`
	Match      string      `json:"match"`
	Position   int         `json:"position,omitempty"`
	Conditions []Condition `json:"conditions"`
	Actions    []Action    `json:"actions"`
}

type AppRuleQuery struct {
	Limit            string  `json:"limit,omitempty"`
	Page             string  `json:"page,omitempty"`
	Cursor           string  `json:"cursor,omitempty"`
	Enabled          bool    `json:"enabled,omitempty"`
	HasCondition     *string `json:"has_condition,omitempty"`
	HasConditionType *string `json:"has_condition_type,omitempty"`
	HasAction        *string `json:"has_action,omitempty"`
	HasActionType    *string `json:"has_action_type,omitempty"`
}

func (q *AppRuleQuery) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit":              validateString,
		"page":               validateString,
		"cursor":             validateString,
		"enabled":            validateBool,
		"has_condition":      validateString,
		"has_condition_type": validateString,
		"has_action":         validateString,
		"has_action_type":    validateString,
	}
}
