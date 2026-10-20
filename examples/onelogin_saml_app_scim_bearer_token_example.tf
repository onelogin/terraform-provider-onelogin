# A SCIM provisioning app with its bearer token. bearer_token_wo is
# write-only: Terraform sends it and never stores it, so it needs
# Terraform 1.11 or later. Change bearer_token_wo_version to send a new
# token; changing the token alone is not seen.
resource onelogin_saml_apps scim {
  connector_id = 45504
  name         = "SCIM bearer token acctest"

  bearer_token_wo         = "acctest-token-one"
  bearer_token_wo_version = 1

  configuration = {
    signature_algorithm = "SHA-256"
  }
}
