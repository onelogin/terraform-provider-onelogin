# configuration does not need signature_algorithm. Left out, a new app gets
# OneLogin's default and an existing app keeps its current value.
resource onelogin_saml_apps no_signature_algorithm {
  connector_id = 45504
  name         = "SAML configuration acctest"

  configuration = {
    audience = "https://acctest.example.com/metadata"
  }
}
