resource onelogin_oidc_apps oidc{
  connector_id = 108419
  name =  "Updated OIDC APP"
  description = "OIDC"

  configuration = {
    redirect_uri = "https://localhost:3000/callback"
    refresh_token_expiration_minutes = 1
    login_url = "https://www.updated.com"
    oidc_application_type = 0
    token_endpoint_auth_method = 1
    access_token_expiration_minutes = 1
  }
}
