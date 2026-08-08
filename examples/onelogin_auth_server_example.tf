resource onelogin_auth_servers test {
  name = "test"
  description = "test"
  configuration {
    resource_identifier = "https://example.com/contacts"
    audiences = ["https://example.com/contacts"]
    refresh_token_expiration_minutes = 30
    access_token_expiration_minutes = 10
  }
}
