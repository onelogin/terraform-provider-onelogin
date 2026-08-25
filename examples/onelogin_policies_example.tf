# A user policy: how people in this tenant sign in.
resource "onelogin_policies" "engineering" {
  name = "Engineering"
  kind = "user"

  # Passwords
  minimum_password_length          = 12
  password_expiration_days         = 90
  passwords_remembered             = 5
  password_complexity_requirements = 3
  enable_password_change           = true
  enable_email_password_reset      = true

  # Lockout
  maximum_invalid_login_attempts = 5
  lock_effective_minutes         = 15

  # MFA
  otp_auth_enabled         = true
  mfa_registration_enabled = true

  # Session: expire after 30 minutes of inactivity.
  session_timeout_by_inactivity_value = 30
  session_timeout_by_inactivity_unit  = 0

  # Restrict sign-in to the corporate network.
  ip_addr_restriction = "10.0.0.0/8"

  terms_and_conditions {
    enabled = true
    content = "Access to this system is monitored and logged."
  }
}

# An app policy: extra assurance in front of one application.
resource "onelogin_policies" "finance_app" {
  name                   = "Finance app step-up"
  kind                   = "app"
  force_authn            = true
  app_force_authn_offset = 60
}
