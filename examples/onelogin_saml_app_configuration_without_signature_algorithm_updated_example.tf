# The audience changed, and signature_algorithm is still not set.
resource onelogin_saml_apps no_signature_algorithm {
  connector_id = 45504
  name         = "SAML configuration acctest"

  configuration = {
    audience = "https://acctest.example.com/metadata/v2"
  }
}
