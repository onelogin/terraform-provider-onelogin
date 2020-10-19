resource onelogin_smarthooks basic_test {
  type = "pre-authentication"
  packages = {
    mysql = "2.18.1"
  }
  env_vars = {
    "API_KEY"
  }
  retries = 0
  timeout = 2
  disabled = false
  status = "ready"
  risk_enabled = false
  location_enabled = false
  function = <<EOF
function myFunc() {
  console.log("WOO WOO")
}
EOF
}
