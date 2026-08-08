resource onelogin_saml_apps saml{
  connector_id = 50534
  name =  "SAML App"
  description = "SAML"

  configuration = {
    signature_algorithm = "SHA-1"
  }
}

# has_role takes a role ID. The example creates the role it checks for rather
# than naming an ID from somebody else's account: the rules endpoint rejects an
# ID it does not recognise with a bare 422 and no indication of which value it
# objected to.
resource onelogin_roles rule_role{
  name = "App Rule Example Role acctest"
}

resource onelogin_app_rules test{
  app_id = onelogin_saml_apps.saml.id
  enabled = true
  match = "all"
  name = "first rule"
  conditions {
    operator = ">"
    source = "last_login"
    value = "90"
  }
  actions {
    action = "set_amazonusername"
    expression = ".*"
    value = ["member_of"]
  }
}
resource onelogin_app_rules check{
  app_id = onelogin_saml_apps.saml.id
  enabled = true
  match = "all"
  name = "second rule"
  conditions {
    operator = "ri"
    source = "has_role"
    value = onelogin_roles.rule_role.id
  }
  actions {
    action = "set_amazonusername"
    expression = ".*"
    value = ["member_of"]
  }
}
