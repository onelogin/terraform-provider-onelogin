resource onelogin_saml_apps saml{
  connector_id = 50534
  name =  "SAML App"
  description = "SAML"

  configuration {
    provider_arn = null
    signature_algorithm = "SHA-1"
  }
}
