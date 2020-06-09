resource onelogin_user_mappings basic_test {
  name = "Select Login"
  enabled = true
  match = "all"
  position = 1

  actions {
    value = ["1"]
    action = "set_status"
  }

  conditions {
    operator = ">"
    source = "last_login"
    value = "90"
  }
}
