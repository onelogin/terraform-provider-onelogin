// Package policy holds the Terraform schema for onelogin_policies and the
// translation between it and the /api/2/policies representation.
package policy

import (
	"sort"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// scope records which kind of policy a field belongs to.
//
// OneLogin keeps user policies and app policies in one table, told apart only
// by `kind`, and most fields mean something to just one of the two:
// password_expiration_days is meaningless on an app policy, force_authn is
// meaningless on a user policy. The API enforces that -- it rejects a write to
// the wrong kind with a 422 and leaves the field out of responses entirely --
// so the same split has to exist here for a plan to be able to say what is
// wrong before the apply reaches OneLogin.
//
// The source of truth is Policy::APP_ONLY_FIELDS and Policy::SHARED_FIELDS in
// the OneLogin API. As there, user-only is not a list: it is whatever is left
// over, so a field added to the table below without an explicit scope counts as
// user-only. That is the safe direction -- a new field is checked rather than
// quietly unchecked -- and it is why scopeUser is the zero value.
type scope int

const (
	scopeUser   scope = iota // user policies only
	scopeApp                 // app policies only
	scopeShared              // meaningful on both kinds
)

// field is one policy attribute that maps straight onto an API field of the
// same name. The attributes that need more than a rename -- name, kind,
// is_default, the two factor ID lists and terms_and_conditions -- are handled
// separately in Schema, RequestBody and Flatten.
type field struct {
	Name        string
	Type        schema.ValueType
	Scope       scope
	Description string
}

// fields lists every plain policy attribute the API presents, in the order and
// grouping Api::V5::PoliciesPresenter uses, so the two can be read side by side.
var fields = []field{
	// Login flow
	{Name: "preferred_auth_state_machine", Type: schema.TypeInt, Description: "Login flow the policy prefers."},

	// Password
	{Name: "minimum_password_length", Type: schema.TypeInt, Description: "Minimum password length. One of 5, 6, 8, 10, 12 or 16."},
	{Name: "password_expiration_days", Type: schema.TypeInt, Description: "Days before a password expires. 0 never expires."},
	{Name: "passwords_remembered", Type: schema.TypeInt, Description: "How many previous passwords cannot be reused. 0, 3 or 5."},
	{Name: "password_complexity_requirements", Type: schema.TypeInt, Description: "Password complexity: 0 none, 1 letters and digits, 2 mixed case and digits, 3 mixed case, digits and special characters, 4 any three of those four."},
	{Name: "dynamic_blacklist_attributes", Type: schema.TypeString, Description: "User attributes whose values may not appear in a password."},
	{Name: "enforce_account_password_blacklist", Type: schema.TypeBool, Description: "Reject passwords on the account's blacklist."},
	{Name: "enable_password_change", Type: schema.TypeBool, Description: "Let users change their own password."},
	{Name: "enable_unlock_via_password_reset", Type: schema.TypeBool, Description: "Unlock a locked account when the user resets their password."},
	{Name: "enable_email_password_reset", Type: schema.TypeBool, Description: "Offer password reset by email."},
	{Name: "enable_sms_password_reset", Type: schema.TypeBool, Description: "Offer password reset by SMS."},
	{Name: "enable_question_password_reset", Type: schema.TypeBool, Description: "Offer password reset by security question."},
	{Name: "require_security_questions", Type: schema.TypeBool, Description: "Require users to set security questions."},
	{Name: "password_redirect_enabled", Type: schema.TypeBool, Description: "Send password changes to an external URL instead."},
	{Name: "password_redirect_url", Type: schema.TypeString, Description: "URL to send password changes to. Required when password_redirect_enabled is true."},
	{Name: "password_redirect_message", Type: schema.TypeString, Description: "Message shown alongside the password redirect."},
	{Name: "enforce_compromised_credentials_check", Type: schema.TypeBool, Description: "Check credentials against known breaches."},

	// Invite link timeout
	{Name: "invite_expiration_time_value", Type: schema.TypeInt, Description: "How long an invite link stays valid, in invite_expiration_time_unit units. Must be greater than 0."},
	{Name: "invite_expiration_time_unit", Type: schema.TypeInt, Description: "Unit for invite_expiration_time_value: 0 minutes, 1 hours, 2 days."},

	// Lockout
	{Name: "maximum_invalid_login_attempts", Type: schema.TypeInt, Description: "Failed logins before lockout. 3 to 10, or 0 for no limit."},
	{Name: "lock_effective_minutes", Type: schema.TypeInt, Description: "Minutes an account stays locked. 15, 30, 60, or 0 to require an admin."},

	// Session
	{Name: "session_timeout_minutes", Type: schema.TypeInt, Description: "Session length in minutes, in the older single-value format."},
	{Name: "session_timeout_type", Type: schema.TypeInt, Description: "Which session timeout applies: by inactivity or at a fixed time."},
	{Name: "session_timeout_by_inactivity_value", Type: schema.TypeInt, Description: "Inactivity timeout, in session_timeout_by_inactivity_unit units."},
	{Name: "session_timeout_by_inactivity_unit", Type: schema.TypeInt, Description: "Unit for session_timeout_by_inactivity_value."},
	{Name: "session_timeout_by_fixed_time_value", Type: schema.TypeInt, Description: "Fixed session length, in session_timeout_by_fixed_time_unit units."},
	{Name: "session_timeout_by_fixed_time_unit", Type: schema.TypeInt, Description: "Unit for session_timeout_by_fixed_time_value."},
	{Name: "persistent_session_enabled", Type: schema.TypeBool, Description: "Let sessions survive a browser restart."},

	// MFA
	{Name: "otp_auth_enabled", Type: schema.TypeBool, Scope: scopeShared, Description: "Require multi-factor authentication."},
	{Name: "otp_config", Type: schema.TypeInt, Description: "Which factors MFA accepts."},
	{Name: "otp_trigger_condition", Type: schema.TypeInt, Description: "When MFA is triggered."},
	{Name: "otp_security_token_expiration_days", Type: schema.TypeInt, Description: "Days a remembered MFA device stays trusted. 1 to 99999."},
	{Name: "mfa_registration_enabled", Type: schema.TypeBool, Description: "Prompt users to register a factor. Combined with voluntary_mfa_registration_enabled: true/false is required, true/true voluntary, false/false not prompted."},
	{Name: "voluntary_mfa_registration_enabled", Type: schema.TypeBool, Description: "Make factor registration voluntary rather than required. See mfa_registration_enabled."},
	{Name: "user_phone_update_allowed", Type: schema.TypeBool, Description: "Let users change their registered phone number."},
	{Name: "disable_protect_push_notifications", Type: schema.TypeBool, Description: "Turn off OneLogin Protect push notifications."},
	{Name: "disable_protect_push_recovery", Type: schema.TypeBool, Description: "Turn off OneLogin Protect push recovery."},
	{Name: "enable_number_match", Type: schema.TypeBool, Description: "Require number matching on push notifications."},

	// IP restriction
	{Name: "ip_addr_restriction", Type: schema.TypeString, Scope: scopeShared, Description: "Newline-separated list of allowed IP addresses or CIDR ranges."},
	{Name: "ignore_xff", Type: schema.TypeBool, Scope: scopeShared, Description: "Ignore the X-Forwarded-For header when matching ip_addr_restriction."},

	// Device trust
	{Name: "browser_cert_required", Type: schema.TypeBool, Scope: scopeShared, Description: "Require a browser certificate."},
	{Name: "browser_pki_expiration", Type: schema.TypeInt, Description: "Days a browser certificate stays valid."},
	{Name: "self_install_cert", Type: schema.TypeBool, Description: "Let users install their own browser certificate."},
	{Name: "third_party_device_trust", Type: schema.TypeBool, Scope: scopeShared, Description: "Require a third-party device trust check."},
	{Name: "gdt_required", Type: schema.TypeBool, Scope: scopeApp, Description: "Require OneLogin Desktop for this app."},
	{Name: "trusted_device_login_enabled", Type: schema.TypeBool, Description: "Allow login from trusted devices."},
	{Name: "trusted_device_login_mfa_allowed", Type: schema.TypeBool, Description: "Allow MFA on trusted device login."},

	// Portal
	{Name: "new_portal_setting", Type: schema.TypeString, Description: "Access to the new portal: required, allowed or forbidden."},
	{Name: "allow_add_company_app", Type: schema.TypeBool, Description: "Let users add company apps to their portal."},
	{Name: "allow_add_personal_app", Type: schema.TypeBool, Description: "Let users add personal apps to their portal."},
	{Name: "enable_browser_extensions", Type: schema.TypeBool, Description: "Allow the OneLogin browser extension."},
	{Name: "disable_browser_password_manager", Type: schema.TypeBool, Description: "Stop the browser offering to save passwords."},
	{Name: "enable_email_hint", Type: schema.TypeBool, Description: "Prefill the email field on the login page."},

	// Social sign-in
	{Name: "social_sign_in", Type: schema.TypeBool, Description: "Allow social sign-in."},
	{Name: "google", Type: schema.TypeBool, Description: "Allow sign-in with Google."},
	{Name: "facebook", Type: schema.TypeBool, Description: "Allow sign-in with Facebook."},
	{Name: "linkedin", Type: schema.TypeBool, Description: "Allow sign-in with LinkedIn."},
	{Name: "twitter", Type: schema.TypeBool, Description: "Allow sign-in with Twitter."},

	// Advanced
	{Name: "euba_enabled", Type: schema.TypeBool, Description: "Enable end-user behaviour analytics."},
	{Name: "euba_risk_threshold", Type: schema.TypeInt, Description: "Risk score above which EUBA acts."},
	{Name: "enable_smart_access", Type: schema.TypeBool, Scope: scopeShared, Description: "Enable SmartAccess risk scoring."},
	{Name: "smart_access_risk_threshold", Type: schema.TypeInt, Scope: scopeShared, Description: "Risk score above which SmartAccess acts."},
	{Name: "track_inactive_users", Type: schema.TypeBool, Description: "Track users who have not logged in recently."},
	{Name: "enable_system_use_notification", Type: schema.TypeBool, Description: "Show a system use notification before login."},
	{Name: "system_use_notification", Type: schema.TypeString, Description: "Text of the system use notification."},

	// App policy settings
	{Name: "force_authn", Type: schema.TypeBool, Scope: scopeApp, Description: "Force re-authentication when the app is opened."},
	{Name: "app_otp_offset", Type: schema.TypeInt, Scope: scopeApp, Description: "Minutes an app MFA prompt is remembered."},
	{Name: "app_otp_offset_enabled", Type: schema.TypeBool, Scope: scopeApp, Description: "Remember an app MFA prompt for app_otp_offset minutes."},
	{Name: "app_force_authn_offset", Type: schema.TypeInt, Scope: scopeApp, Description: "Minutes before force_authn applies again."},

	// Secure admin and profile areas
	{Name: "secure_admin", Type: schema.TypeBool, Description: "Require step-up authentication to reach the admin area."},
	{Name: "secure_profile", Type: schema.TypeBool, Description: "Require step-up authentication to reach the user profile area."},
	{Name: "admin_policy_id", Type: schema.TypeInt, Description: "App policy that governs step-up authentication for the admin area."},
	{Name: "profile_policy_id", Type: schema.TypeInt, Description: "App policy that governs step-up authentication for the profile area."},
	{Name: "secure_area_otp_timeout_minutes", Type: schema.TypeInt, Description: "Minutes a step-up authentication lasts. One of 5, 10, 15, 20, 30, 45 or 60."},
}

// Schema returns the Terraform schema for a OneLogin security policy.
func Schema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the policy.",
		},
		// ForceNew because the API refuses to change it: the update adapter
		// rejects a body whose kind differs from the stored one, so the only
		// way to move a policy between kinds is to replace it.
		"kind": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"user", "app"}, false),
			Description:  "Whether this is a `user` policy or an `app` policy. Changing it replaces the policy.",
		},
		// Read-only. Which policy is the account default is a property of the
		// account, set through a separate endpoint, and there is no way to
		// stop being the default -- only to make some other policy it.
		"is_default": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether this is the account's default user policy. Set elsewhere; read-only here.",
		},
		"authentication_factor_ids": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "IDs of the authentication factors this policy accepts. Setting it replaces the whole list.",
		},
		"reset_password_authentication_factor_ids": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "IDs of the authentication factors accepted for password reset. Setting it replaces the whole list.",
		},
		"terms_and_conditions": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Whether users must accept the terms before signing in.",
					},
					"content": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "Text of the terms and conditions.",
					},
				},
			},
			Description: "Terms and conditions users must accept.",
		},
	}

	// Every plain field is Optional and Computed. Optional because none of them
	// has to be set, and Computed because the API fills in a default for most
	// of them on create -- without it, a policy the practitioner configured
	// with two fields would show a diff for every other field on the next plan.
	//
	// The cost is that removing an attribute from the configuration no longer
	// resets it: state keeps the last value the API returned and nothing is
	// sent. Changing a value back explicitly is the way to undo one.
	for _, f := range fields {
		s[f.Name] = &schema.Schema{
			Type:        f.Type,
			Optional:    true,
			Computed:    true,
			Description: f.Description,
		}
	}

	return s
}

// ConfiguredKeys returns the names of the top-level attributes the practitioner
// actually wrote, so that an attribute left out can be told from one written as
// false, 0 or "". d.Get cannot make that distinction -- it returns the zero
// value for both -- and the difference decides whether the attribute is sent.
//
// raw comes from d.GetRawConfig(). It is null outside a plan or an apply, and
// nil is returned in that case to mean "no configuration information", which is
// not the same answer as "nothing was configured"; see requestValue.
func ConfiguredKeys(raw cty.Value) map[string]bool {
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsObjectType() {
		return nil
	}

	keys := map[string]bool{}
	for name, value := range raw.AsValueMap() {
		if value.IsNull() {
			continue
		}
		keys[name] = true
	}
	return keys
}

// FieldsNotApplicableTo returns the configured attributes that do not belong on
// a policy of the given kind, sorted, mirroring Policy.fields_not_applicable_to
// in the OneLogin API. An unrecognised kind has nothing inapplicable, matching
// the API and leaving the complaint to the kind validator.
func FieldsNotApplicableTo(kind string, configured map[string]bool) []string {
	var wrong []string

	for _, f := range fields {
		if !configured[f.Name] || f.Scope == scopeShared {
			continue
		}
		switch {
		case kind == "app" && f.Scope == scopeUser:
			wrong = append(wrong, f.Name)
		case kind == "user" && f.Scope == scopeApp:
			wrong = append(wrong, f.Name)
		}
	}

	// The two composite attributes are not in the table. Both belong to user
	// policies only: reset password factors are part of account recovery, and
	// so are the terms a user accepts when signing in.
	if kind == "app" {
		for _, name := range []string{"reset_password_authentication_factor_ids", "terms_and_conditions"} {
			if configured[name] {
				wrong = append(wrong, name)
			}
		}
	}

	sort.Strings(wrong)
	return wrong
}

// Getter is the part of *schema.ResourceData that RequestBody reads, so that
// the body can be built from anything holding the same values.
type Getter interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
}

// RequestBody returns the JSON body for a create or an update, carrying only the
// attributes the configuration sets.
//
// Sending only configured attributes matters for two reasons. Every attribute
// here is Optional and Computed, so one the practitioner left out still holds
// whatever the API last returned; sending that back would have the provider
// assert a value nobody asked for. And the API rejects any field belonging to
// the other kind of policy, so a body built from everything in state would be a
// 422 on every app policy rather than a wrong value on one field.
//
// configured comes from ConfiguredKeys and may be nil; see requestValue.
func RequestBody(d Getter, configured map[string]bool) map[string]interface{} {
	// name is Required, so it is always configured and always belongs in the
	// body. kind is Required too but is not sent here: the API refuses to
	// change it, and the resource adds it to the create body alone.
	body := map[string]interface{}{"name": d.Get("name")}

	for _, f := range fields {
		if value, ok := requestValue(d, configured, f.Name); ok {
			body[f.Name] = value
		}
	}

	for _, name := range []string{"authentication_factor_ids", "reset_password_authentication_factor_ids"} {
		if _, ok := requestValue(d, configured, name); !ok {
			continue
		}
		ids := []int{}
		if set, ok := d.Get(name).(*schema.Set); ok {
			for _, id := range set.List() {
				ids = append(ids, id.(int))
			}
			sort.Ints(ids)
		}
		body[name] = ids
	}

	if _, ok := requestValue(d, configured, "terms_and_conditions"); ok {
		if terms, ok := d.Get("terms_and_conditions").([]interface{}); ok && len(terms) > 0 {
			if block, ok := terms[0].(map[string]interface{}); ok {
				body["terms_and_conditions"] = map[string]interface{}{
					"enabled": block["enabled"],
					"content": block["content"],
				}
			}
		}
	}

	return body
}

// requestValue returns an attribute's value and whether it should be sent.
//
// When configured is nil the raw configuration was not available, and the
// fallback is d.GetOk, which treats a zero value as unset. That can drop a
// deliberate `false` or `0`, but the alternative is sending values the
// practitioner never wrote, and on an app policy those are a 422 rather than a
// wrong value. Terraform populates the raw configuration for create and update,
// so this only matters somewhere that does not, such as a test harness.
func requestValue(d Getter, configured map[string]bool, name string) (interface{}, bool) {
	if configured == nil {
		return d.GetOk(name)
	}
	if !configured[name] {
		return nil, false
	}
	return d.Get(name), true
}

// Flatten writes a policy from the API into ResourceData.
//
// Only the keys the response carries are written. The API leaves out every
// field belonging to the other kind of policy, and writing a zero value in
// their place would put a value in state for something the policy does not
// have.
func Flatten(d *schema.ResourceData, policy map[string]interface{}) error {
	for _, name := range []string{"name", "kind", "is_default"} {
		value, ok := policy[name]
		if !ok {
			continue
		}
		if err := d.Set(name, value); err != nil {
			return err
		}
	}

	for _, f := range fields {
		value, ok := policy[f.Name]
		if !ok || value == nil {
			continue
		}
		if f.Type == schema.TypeInt {
			number, ok := toInt(value)
			if !ok {
				continue
			}
			value = number
		}
		if err := d.Set(f.Name, value); err != nil {
			return err
		}
	}

	// Both lists are written whenever the response mentions them, empty
	// included. An Optional+Computed collection left unwritten has no value in
	// state at all, which Terraform plans as "known after apply" on every run
	// -- a plan that is never empty, for a policy that has simply not got any
	// factors.
	for _, name := range []string{"authentication_factor_ids", "reset_password_authentication_factor_ids"} {
		raw, present := policy[name]
		if !present {
			continue
		}
		ids := []interface{}{}
		if list, ok := raw.([]interface{}); ok {
			for _, id := range list {
				if number, ok := toInt(id); ok {
					ids = append(ids, number)
				}
			}
		}
		if err := d.Set(name, ids); err != nil {
			return err
		}
	}

	// A policy with no terms contract presents terms_and_conditions as null,
	// which becomes an empty block list for the same reason as the factor
	// lists above.
	if raw, present := policy["terms_and_conditions"]; present {
		terms := []interface{}{}
		if block, ok := raw.(map[string]interface{}); ok {
			terms = append(terms, map[string]interface{}{
				"enabled": block["enabled"],
				"content": block["content"],
			})
		}
		if err := d.Set("terms_and_conditions", terms); err != nil {
			return err
		}
	}

	return nil
}

// toInt converts a number from a decoded JSON response. Numbers arrive as
// float64 from encoding/json, but a caller that built the map itself may have
// used int, so both are accepted.
func toInt(value interface{}) (int, bool) {
	switch number := value.(type) {
	case float64:
		return int(number), true
	case int:
		return number, true
	case int32:
		return int(number), true
	case int64:
		return int(number), true
	default:
		return 0, false
	}
}
