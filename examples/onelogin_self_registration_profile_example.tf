resource onelogin_self_registration_profiles basic_test {
  name = "Test Profile"
  url = "test-profile"
  enabled = true
  moderated = false
  helptext = "Welcome to our community. Please follow the guidelines."
  thankyou_message = "Thank you for joining!"
  domain_blacklist = "spam.com, banned.com"
  domain_whitelist = "example.com, trusted.com"
  domain_list_strategy = 1
  email_verification_type = "Email MagicLink"
}
