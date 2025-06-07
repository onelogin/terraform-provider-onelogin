resource onelogin_self_registration_profiles basic_test {
  name = "Updated Test Profile"
  url = "test-profile"
  enabled = false
  moderated = true
  helptext = "Updated welcome message."
  thankyou_message = "Updated thank you message!"
  domain_blacklist = "spam.com, banned.com, malicious.com"
  domain_whitelist = "example.com, trusted.com, safe.com"
  domain_list_strategy = 1
  email_verification_type = "Email OTP"
}
