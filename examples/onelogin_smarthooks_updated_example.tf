resource onelogin_smarthooks basic_test {
  type = "pre-authentication"
  env_vars = []
  packages = {
    mysql = "2.18.1"
  }
  retries = 0
  timeout = 1
  runtime =  "nodejs12.x"
  disabled = false
  options = {
    risk_enabled = false
    location_enabled = false
  }
  function = "ICAgIGV4cG9ydHMuaGFuZGxlciA9IGFzeW5jIGNvbnRleHQgPT4gewogICAgICBjb25zb2xlLmxvZygiUHJlLWF1dGggZXhlY3V0aW5nIGZvciAiICsgY29udGV4dC51c2VyLnVzZXJfaWRlbnRpZmllcik7CiAgICAgIHJldHVybiB7IHVzZXI6IGNvbnRleHQudXNlciB9OwogICAgfTsK"
}
