resource onelogin_apps basic_test {
  connector_id = 20938
  name = "Form-Based Fitbit App"
  description = "Basic Form-Based Application"
  
  # Explicitly set these parameters
  visible = true
  allow_assumed_signin = false
  
  # Add provisioning settings
  provisioning = {
    enabled = false
  }
  
  # Parameters might be required for some connectors
  # Uncomment and modify if needed for your specific connector
  # parameters {
  #   param_key_name = "username"
  #   label = "Username"
  #   user_attribute_mappings = "username"
  # }
  
  # parameters {
  #   param_key_name = "password"
  #   label = "Password" 
  # }
}
