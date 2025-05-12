resource onelogin_saml_apps saml{
  connector_id = 50534
  name =  "Updated SAML App"
  description = "Updated SAML"

  parameters {
    param_key_name = "email"
    label = "Email Address"  # Updated label
    user_attribute_mappings = "email"
    include_in_saml_assertion = true
  }

  parameters {
    param_key_name = "firstname"
    label = "First Name"
    user_attribute_mappings = "firstname"
    include_in_saml_assertion = true
  }

  parameters {
    param_key_name = "lastname"
    label = "Last Name"
    user_attribute_mappings = "lastname"
    include_in_saml_assertion = true
  }

  # Updated custom attribute
  parameters {
    param_key_name = "department"
    label = "Department Name"  # Updated label
    user_attribute_mappings = "custom_attribute_department"
    include_in_saml_assertion = true
  }

  # Added new parameter
  parameters {
    param_key_name = "title"
    label = "Job Title"
    user_attribute_mappings = "title"
    include_in_saml_assertion = true
  }

  configuration = {
    signature_algorithm = "SHA-256"
  }
}
