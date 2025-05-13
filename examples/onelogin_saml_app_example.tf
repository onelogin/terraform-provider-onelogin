resource onelogin_saml_apps saml{
  connector_id = 50534
  name =  "SAML App"
  description = "SAML"

  parameters {
    param_key_name = "email"
    label = "Email"
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

  # Example of using a custom attribute
  parameters {
    param_key_name = "department"
    label = "Department"
    user_attribute_mappings = "custom_attribute_department"
    include_in_saml_assertion = true
  }

  configuration = {
    signature_algorithm = "SHA-1"
  }
}
