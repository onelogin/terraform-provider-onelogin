# The token rotated: a new value and a new version, so the update sends it.
resource onelogin_saml_apps scim {
  connector_id = 45504
  name         = "SCIM bearer token acctest"

  bearer_token_wo         = "acctest-token-two"
  bearer_token_wo_version = 2

  configuration = {
    signature_algorithm = "SHA-256"
  }
}
