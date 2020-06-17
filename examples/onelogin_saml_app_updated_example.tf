resource onelogin_saml_apps saml{
  connector_id = 50534
  name =  "Updated SAML App"
  description = "Updated SAML"

  configuration = {
    signature_algorithm = "SHA-256"
  }
  rules {
    enabled = true
    match = "all"
    name = "second rule"
    conditions {
      operator = "ri"
      source = "has_role"
      value = "340475"
    }
    actions {
      action = "set_amazonusername"
      expression = ".*"
      value = ["member_of"]
    }
  }
  rules {
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
}
