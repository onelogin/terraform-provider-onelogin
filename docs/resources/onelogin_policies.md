---
layout: "onelogin"
page_title: "OneLogin: onelogin_policies"
sidebar_current: "docs-onelogin-resource-policies"
description: |-
  Create and configure OneLogin security policies.
---

# onelogin_policies

Manages a OneLogin security policy.

A policy is one of two kinds, set by `kind` and fixed for the life of the policy:

* A **user** policy governs how people sign in — passwords, MFA, lockout, session length,
  self-service password reset, the portal, and so on. Assign it to users through the
  `policy_id` on a role or group, or make it the account default in the OneLogin admin UI.
* An **app** policy governs a single application — whether opening it forces
  re-authentication, and how long that lasts.

The two share one API resource, so most arguments belong to only one kind. Each argument
below says which. Setting one on the wrong kind is refused during `terraform plan`, before
anything reaches OneLogin.

Custom security policies are a licensed feature. If your plan does not include
**Custom Security Policies**, creating one returns `406 Not Acceptable`.

## Example Usage

```hcl
resource "onelogin_policies" "engineering" {
  name = "Engineering"
  kind = "user"

  # Passwords
  minimum_password_length          = 12
  password_expiration_days         = 90
  passwords_remembered             = 5
  password_complexity_requirements = 3

  # Lockout
  maximum_invalid_login_attempts = 5
  lock_effective_minutes         = 15

  # MFA
  otp_auth_enabled          = true
  mfa_registration_enabled  = true
  authentication_factor_ids = [12345, 67890]

  # Session
  session_timeout_by_inactivity_value = 30
  session_timeout_by_inactivity_unit  = 0

  terms_and_conditions {
    enabled = true
    content = "Access is monitored and logged."
  }
}

resource "onelogin_policies" "finance_app" {
  name                   = "Finance app step-up"
  kind                   = "app"
  force_authn            = true
  app_force_authn_offset = 60
}
```

## Argument Reference

* `name` - (Required, String) Name of the policy.
* `kind` - (Required, String) Either `user` or `app`. Changing it replaces the policy, because
  the API refuses to move an existing one between kinds.
* `authentication_factor_ids` - (Optional, Set of Number) IDs of the authentication factors this
  policy accepts. Setting it replaces the whole list; setting it to `[]` clears it. Applies to
  both kinds.
* `reset_password_authentication_factor_ids` - (Optional, Set of Number) IDs of the authentication
  factors accepted for password reset. Applies to user policies only.
* `terms_and_conditions` - (Optional, Block, max 1) Terms users must accept before signing in.
  Applies to user policies only.
  * `enabled` - (Optional, Boolean) Whether the terms are shown and must be accepted.
  * `content` - (Optional, String) Text of the terms.

Every remaining argument is optional, and OneLogin supplies a default for most of them. An
argument you do not set keeps whatever value the API reports, so the plan stays empty rather
than showing a diff for a default you never chose. The consequence is that **removing an
argument from your configuration does not reset it** — state keeps the last value OneLogin
returned and nothing is sent. To undo a setting, set it explicitly to the value you want.

### Login flow

* `preferred_auth_state_machine` - (Optional, Number) Login flow the policy prefers. Applies to user policies only.

### Password

* `minimum_password_length` - (Optional, Number) Minimum password length. One of 5, 6, 8, 10, 12 or 16. Applies to user policies only.
* `password_expiration_days` - (Optional, Number) Days before a password expires. 0 never expires. Applies to user policies only.
* `passwords_remembered` - (Optional, Number) How many previous passwords cannot be reused. 0, 3 or 5. Applies to user policies only.
* `password_complexity_requirements` - (Optional, Number) Password complexity: 0 none, 1 letters and digits, 2 mixed case and digits, 3 mixed case, digits and special characters, 4 any three of those four. Applies to user policies only.
* `dynamic_blacklist_attributes` - (Optional, String) User attributes whose values may not appear in a password. Applies to user policies only.
* `enforce_account_password_blacklist` - (Optional, Boolean) Reject passwords on the account's blacklist. Applies to user policies only.
* `enable_password_change` - (Optional, Boolean) Let users change their own password. Applies to user policies only.
* `enable_unlock_via_password_reset` - (Optional, Boolean) Unlock a locked account when the user resets their password. Applies to user policies only.
* `enable_email_password_reset` - (Optional, Boolean) Offer password reset by email. Applies to user policies only.
* `enable_sms_password_reset` - (Optional, Boolean) Offer password reset by SMS. Applies to user policies only.
* `enable_question_password_reset` - (Optional, Boolean) Offer password reset by security question. Applies to user policies only.
* `require_security_questions` - (Optional, Boolean) Require users to set security questions. Applies to user policies only.
* `password_redirect_enabled` - (Optional, Boolean) Send password changes to an external URL instead. Applies to user policies only.
* `password_redirect_url` - (Optional, String) URL to send password changes to. Required when password_redirect_enabled is true. Applies to user policies only.
* `password_redirect_message` - (Optional, String) Message shown alongside the password redirect. Applies to user policies only.
* `enforce_compromised_credentials_check` - (Optional, Boolean) Check credentials against known breaches. Applies to user policies only.

### Invite link timeout

* `invite_expiration_time_value` - (Optional, Number) How long an invite link stays valid, in invite_expiration_time_unit units. Must be greater than 0. Applies to user policies only.
* `invite_expiration_time_unit` - (Optional, Number) Unit for invite_expiration_time_value: 0 minutes, 1 hours, 2 days. Applies to user policies only.

### Lockout

* `maximum_invalid_login_attempts` - (Optional, Number) Failed logins before lockout. 3 to 10, or 0 for no limit. Applies to user policies only.
* `lock_effective_minutes` - (Optional, Number) Minutes an account stays locked. 15, 30, 60, or 0 to require an admin. Applies to user policies only.

### Session

* `session_timeout_minutes` - (Optional, Number) Session length in minutes, in the older single-value format. Applies to user policies only.
* `session_timeout_type` - (Optional, Number) Which session timeout applies: by inactivity or at a fixed time. Applies to user policies only.
* `session_timeout_by_inactivity_value` - (Optional, Number) Inactivity timeout, in session_timeout_by_inactivity_unit units. Applies to user policies only.
* `session_timeout_by_inactivity_unit` - (Optional, Number) Unit for session_timeout_by_inactivity_value. Applies to user policies only.
* `session_timeout_by_fixed_time_value` - (Optional, Number) Fixed session length, in session_timeout_by_fixed_time_unit units. Applies to user policies only.
* `session_timeout_by_fixed_time_unit` - (Optional, Number) Unit for session_timeout_by_fixed_time_value. Applies to user policies only.
* `persistent_session_enabled` - (Optional, Boolean) Let sessions survive a browser restart. Applies to user policies only.

### MFA

* `otp_auth_enabled` - (Optional, Boolean) Require multi-factor authentication. Applies to both kinds.
* `otp_config` - (Optional, Number) Which factors MFA accepts. Applies to user policies only.
* `otp_trigger_condition` - (Optional, Number) When MFA is triggered. Applies to user policies only.
* `otp_security_token_expiration_days` - (Optional, Number) Days a remembered MFA device stays trusted. 1 to 99999. Applies to user policies only.
* `mfa_registration_enabled` - (Optional, Boolean) Prompt users to register a factor. Combined with voluntary_mfa_registration_enabled: true/false is required, true/true voluntary, false/false not prompted. Applies to user policies only.
* `voluntary_mfa_registration_enabled` - (Optional, Boolean) Make factor registration voluntary rather than required. See mfa_registration_enabled. Applies to user policies only.
* `user_phone_update_allowed` - (Optional, Boolean) Let users change their registered phone number. Applies to user policies only.
* `disable_protect_push_notifications` - (Optional, Boolean) Turn off OneLogin Protect push notifications. Applies to user policies only.
* `disable_protect_push_recovery` - (Optional, Boolean) Turn off OneLogin Protect push recovery. Applies to user policies only.
* `enable_number_match` - (Optional, Boolean) Require number matching on push notifications. Applies to user policies only.

### IP restriction

* `ip_addr_restriction` - (Optional, String) Newline-separated list of allowed IP addresses or CIDR ranges. Applies to both kinds.
* `ignore_xff` - (Optional, Boolean) Ignore the X-Forwarded-For header when matching ip_addr_restriction. Applies to both kinds.

### Device trust

* `browser_cert_required` - (Optional, Boolean) Require a browser certificate. Applies to both kinds.
* `browser_pki_expiration` - (Optional, Number) Days a browser certificate stays valid. Applies to user policies only.
* `self_install_cert` - (Optional, Boolean) Let users install their own browser certificate. Applies to user policies only.
* `third_party_device_trust` - (Optional, Boolean) Require a third-party device trust check. Applies to both kinds.
* `gdt_required` - (Optional, Boolean) Require OneLogin Desktop for this app. Applies to app policies only.
* `trusted_device_login_enabled` - (Optional, Boolean) Allow login from trusted devices. Applies to user policies only.
* `trusted_device_login_mfa_allowed` - (Optional, Boolean) Allow MFA on trusted device login. Applies to user policies only.

### Portal

* `new_portal_setting` - (Optional, String) Access to the new portal: required, allowed or forbidden. Applies to user policies only.
* `allow_add_company_app` - (Optional, Boolean) Let users add company apps to their portal. Applies to user policies only.
* `allow_add_personal_app` - (Optional, Boolean) Let users add personal apps to their portal. Applies to user policies only.
* `enable_browser_extensions` - (Optional, Boolean) Allow the OneLogin browser extension. Applies to user policies only.
* `disable_browser_password_manager` - (Optional, Boolean) Stop the browser offering to save passwords. Applies to user policies only.
* `enable_email_hint` - (Optional, Boolean) Prefill the email field on the login page. Applies to user policies only.

### Social sign-in

* `social_sign_in` - (Optional, Boolean) Allow social sign-in. Applies to user policies only.
* `google` - (Optional, Boolean) Allow sign-in with Google. Applies to user policies only.
* `facebook` - (Optional, Boolean) Allow sign-in with Facebook. Applies to user policies only.
* `linkedin` - (Optional, Boolean) Allow sign-in with LinkedIn. Applies to user policies only.
* `twitter` - (Optional, Boolean) Allow sign-in with Twitter. Applies to user policies only.

### Advanced

* `euba_enabled` - (Optional, Boolean) Enable end-user behaviour analytics. Applies to user policies only.
* `euba_risk_threshold` - (Optional, Number) Risk score above which EUBA acts. Applies to user policies only.
* `enable_smart_access` - (Optional, Boolean) Enable SmartAccess risk scoring. Applies to both kinds.
* `smart_access_risk_threshold` - (Optional, Number) Risk score above which SmartAccess acts. Applies to both kinds.
* `track_inactive_users` - (Optional, Boolean) Track users who have not logged in recently. Applies to user policies only.
* `enable_system_use_notification` - (Optional, Boolean) Show a system use notification before login. Applies to user policies only.
* `system_use_notification` - (Optional, String) Text of the system use notification. Applies to user policies only.

### App policy settings

* `force_authn` - (Optional, Boolean) Force re-authentication when the app is opened. Applies to app policies only.
* `app_otp_offset` - (Optional, Number) Minutes an app MFA prompt is remembered. Applies to app policies only.
* `app_otp_offset_enabled` - (Optional, Boolean) Remember an app MFA prompt for app_otp_offset minutes. Applies to app policies only.
* `app_force_authn_offset` - (Optional, Number) Minutes before force_authn applies again. Applies to app policies only.

### Secure admin and profile areas

* `secure_admin` - (Optional, Boolean) Require step-up authentication to reach the admin area. Applies to user policies only.
* `secure_profile` - (Optional, Boolean) Require step-up authentication to reach the user profile area. Applies to user policies only.
* `admin_policy_id` - (Optional, Number) App policy that governs step-up authentication for the admin area. Applies to user policies only.
* `profile_policy_id` - (Optional, Number) App policy that governs step-up authentication for the profile area. Applies to user policies only.
* `secure_area_otp_timeout_minutes` - (Optional, Number) Minutes a step-up authentication lasts. One of 5, 10, 15, 20, 30, 45 or 60. Applies to user policies only.
## Attribute Reference

* `id` - The policy ID.
* `is_default` - Whether this is the account's default user policy. Read-only: which policy is
  the default belongs to the account and is set elsewhere, and a policy cannot stop being the
  default except by another policy becoming it.

## Import

Policies are imported by ID:

```
terraform import onelogin_policies.engineering 123456
```
