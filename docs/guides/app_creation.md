---
layout: "onelogin"
page_title: "OneLogin: App Creation Guide"
sidebar_current: "docs-onelogin-guide-app-creation"
description: |-
  Guide to creating different types of applications in OneLogin.
---

# App Creation Guide

## Choosing the Right Resource Type

The OneLogin Terraform Provider offers several different resources for creating applications:

- `onelogin_apps` - Basic app resource (limited functionality)
- `onelogin_saml_apps` - SAML application with specialized configuration
- `onelogin_oidc_apps` - OIDC/OAuth application with specialized configuration

Always prefer the specialized resources (`onelogin_saml_apps`, `onelogin_oidc_apps`) over the basic `onelogin_apps` resource when possible. The specialized resources are designed to handle the specific requirements of different connector types and provide appropriate validation.

## Common Issues and Solutions

### 1. 422 Unprocessable Entity Errors

When creating apps, you may encounter 422 errors if required fields are missing:

```
Error: request failed with status: 422
```

This typically happens because:
- Required fields for the specific connector type are missing
- The connector ID is invalid or not supported

### 2. Required Fields Based on Connector Type

Different connector types require different configuration fields:

#### SAML Apps
For SAML apps, you typically need:
- `connector_id` - Use an appropriate SAML connector ID (e.g., 110005 for Generic SAML 2.0)
- `configuration` with at least:
  - `signature_algorithm` - (e.g., "SHA-1", "SHA-256")
  - Often requires additional fields like ACS URL, Entity ID, etc.

Example:
```hcl
resource "onelogin_saml_apps" "example" {
  name = "Example SAML App"
  connector_id = 110005
  
  configuration = {
    signature_algorithm = "SHA-1"
  }
  
  # SAML apps typically require specific parameters
  parameters {
    param_key_name = "acs_url"
    label = "ACS URL"
    user_attribute_mappings = "_macro_"
    user_attribute_macros = "https://example.com/acs"
    include_in_saml_assertion = true
  }
}
```

#### OIDC Apps
For OIDC apps, you typically need:
- `connector_id` - Use an appropriate OIDC connector ID
- `configuration` with fields like:
  - `redirect_uri`
  - `token_endpoint_auth_method`
  - `grant_types`

### 3. Computed vs. Required Fields

Some fields in the provider schema are marked as `Computed` (read-only) but may actually be required by the API. This can lead to confusing errors, especially with the basic `onelogin_apps` resource.

The specialized resources (`onelogin_saml_apps`, `onelogin_oidc_apps`) handle this better by setting these fields appropriately behind the scenes.

## Connector IDs

Here are some common connector IDs:

- **110005**: Generic SAML 2.0 (Custom Connector)
- **110003**: OpenID Connect
- **13579**: Generic SCIM Provisioning

You can find more connector IDs by:
1. Creating an app manually in the OneLogin admin portal
2. Looking at the connector ID in the URL or app details

## Testing App Creation

When testing app creation, consider:

1. Start with the minimal required configuration for your connector type
2. Add fields incrementally if you encounter errors
3. Check the OneLogin admin portal for examples of working configurations
4. Use the specialized app resources (`onelogin_saml_apps`, `onelogin_oidc_apps`) instead of the basic `onelogin_apps`