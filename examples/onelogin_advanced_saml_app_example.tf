resource onelogin_saml_apps saml_advanced {
    connector_id = 110016 # SAML 2.0 Advanced connector ID
    name =  local.hostname
    description = "Advanced SAML app with encrypted asserts"

    configuration = {
        signature_algorithm = "SHA-256"
        logout_url = "https://${local.hostname}/auth/?sls"
        saml_notonorafter = "3"
        audience = "https://${local.hostname}/auth/metadata/"
        generate_attribute_value_tags = "0"
        saml_initiater_id = "0"
        saml_notbefore = "3"
        saml_issuer_type = "0"
        saml_sign_element = "0"
        encrypt_assertion = "1"
        consumer_url = "https://${local.hostname}/auth/?acs"
        login = "https://${local.hostname}/"
        saml_sessionnotonorafter = "1440"
        saml_encryption_method_id = "0"
        recipient = "https://${local.hostname}/auth/?acs"
        validator = ".*"
        relaystate = "https://${local.hostname}/"
        saml_nameid_format_id = "0"
    }
}