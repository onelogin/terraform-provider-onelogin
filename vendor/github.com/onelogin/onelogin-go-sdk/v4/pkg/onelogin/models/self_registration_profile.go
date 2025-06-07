package models

// SelfRegistrationProfile represents a OneLogin Self-Registration Profile
type SelfRegistrationProfile struct {
	ID                   int32                         `json:"id,omitempty"`
	Name                 string                        `json:"name"`
	URL                  string                        `json:"url"`
	Enabled              bool                          `json:"enabled"`
	Moderated            bool                          `json:"moderated"`
	DefaultRoleID        int32                         `json:"default_role_id,omitempty"`
	DefaultGroupID       int32                         `json:"default_group_id,omitempty"`
	Helptext             string                        `json:"helptext,omitempty"`
	ThankyouMessage      string                        `json:"thankyou_message,omitempty"`
	DomainBlacklist      string                        `json:"domain_blacklist,omitempty"`
	DomainWhitelist      string                        `json:"domain_whitelist,omitempty"`
	DomainListStrategy   int32                         `json:"domain_list_strategy,omitempty"`
	EmailVerificationType string                       `json:"email_verification_type,omitempty"`
	Fields               []SelfRegistrationProfileField `json:"fields,omitempty"`
}

// SelfRegistrationProfileField represents a field in a Self-Registration Profile
type SelfRegistrationProfileField struct {
	ID                      int32  `json:"id,omitempty"`
	CustomAttributeID       int32  `json:"custom_attribute_id"`
	Name                    string `json:"name,omitempty"`
	Position                int32  `json:"position,omitempty"`
	SelfRegistrationProfileID int32 `json:"self_registration_profile_id,omitempty"`
}

// SelfRegistrationProfileQuery represents available query parameters for self-registration profiles
type SelfRegistrationProfileQuery struct {
	Limit  string `json:"limit,omitempty"`
	Page   string `json:"page,omitempty"`
}

// Domain list strategy constants
const (
	DomainBlacklistStrategy int32 = 0
	DomainWhitelistStrategy int32 = 1
)

// Email verification type constants
const (
	EmailMagicLink string = "Email MagicLink"
	EmailOTP       string = "Email OTP"
)

// GetKeyValidators returns the validators for the query parameters
func (q *SelfRegistrationProfileQuery) GetKeyValidators() map[string]func(interface{}) bool {
	return map[string]func(interface{}) bool{
		"limit": validateString,
		"page":  validateString,
	}
}
