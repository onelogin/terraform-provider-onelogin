resource onelogin_saml_apps saml{
  connector_id = 50534
  name =  "Updated SAML App"
  description = "Updated SAML"

  configuration {
    signature_algorithm = "SHA-256"
  }
}
