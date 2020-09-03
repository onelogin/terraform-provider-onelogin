resource onelogin_rules test{
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
resource onelogin_rules check{
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
