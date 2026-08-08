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

  # Several values in one parameter -- the shape #239 reported the provider
  # could not represent at all. A single value is still written as a list;
  # the provider sends it to the API as a bare string, which is the shape the
  # API has always stored for it.
  parameters {
    param_key_name = "multivalued"
    label          = "Groups"
    default_values = ["one", "two"]
  }

  configuration = {
    signature_algorithm = "SHA-1"
  }
}
