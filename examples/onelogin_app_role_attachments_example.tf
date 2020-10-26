resource onelogin_saml_apps saml {
  connector_id = 50534
  name =  "SAML App"
  description = "SAML"

  configuration = {
    signature_algorithm = "SHA-1"
  }
}

resource onelogin_app_role_attachments test {
	app_id = onelogin_saml_apps.saml.id
	role_id = 376940
}

resource onelogin_app_role_attachments check {
	app_id = onelogin_saml_apps.saml.id
	role_id = 376941
}
