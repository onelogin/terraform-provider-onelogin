resource onelogin_auth_servers test {
  name = "updated"
  description = "updated test"
  configuration = {
    resource_identifier = "https://example.com/users/contacts"
    audiences = ["https://example.com/contacts", "https://example.com/users/contacts"]
    refresh_token_expiration_minutes = 30
    access_token_expiration_minutes = 10
  }
}
