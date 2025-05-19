---
layout: "onelogin"
page_title: "OneLogin: Resource Naming Conventions"
sidebar_current: "docs-onelogin-guide-resource-naming-conventions"
description: |-
  Guide to understanding resource naming conventions in the OneLogin Terraform Provider.
---

# Resource Naming Conventions

## Plural Resource Names

The OneLogin Terraform Provider uses plural resource names for all resources, even when they represent a single entity. This is a convention used consistently throughout the provider.

For example:
- Use `onelogin_roles` (not `onelogin_role`)
- Use `onelogin_apps` (not `onelogin_app`)
- Use `onelogin_users` (not `onelogin_user`)

```hcl
# Correct usage
resource "onelogin_roles" "executive_admin" {
  name = "Executive Admin"
}

# Incorrect usage
resource "onelogin_role" "executive_admin" {
  name = "Executive Admin"
}
```

## App-Role Attachments

The `onelogin_app_role_attachments` resource creates a one-to-one relationship between an app and a role, despite the plural name. To attach multiple roles to an app, you need to create multiple `onelogin_app_role_attachments` resources.

```hcl
# To attach multiple roles to an app, create multiple attachment resources
resource "onelogin_apps" "my_app" {
  name = "My Application"
  # other app properties...
}

resource "onelogin_roles" "admin_role" {
  name = "Admin Role"
  # other role properties...
}

resource "onelogin_roles" "user_role" {
  name = "User Role"
  # other role properties...
}

# Create one attachment resource for each role
resource "onelogin_app_role_attachments" "admin_attachment" {
  app_id = onelogin_apps.my_app.id
  role_id = onelogin_roles.admin_role.id
}

resource "onelogin_app_role_attachments" "user_attachment" {
  app_id = onelogin_apps.my_app.id
  role_id = onelogin_roles.user_role.id
}
```

## Specialized App Resources

For creating applications, the provider offers specialized resources that handle specific app types:

- `onelogin_apps` - Basic app resource (limited functionality)
- `onelogin_saml_apps` - SAML application with specialized configuration
- `onelogin_oidc_apps` - OIDC/OAuth application with specialized configuration

Always prefer the specialized resources when possible, as they handle the specific requirements of each connector type better than the basic resource.

For example:

```hcl
# Preferred: Use specialized SAML app resource
resource "onelogin_saml_apps" "example" {
  name = "Example SAML App"
  connector_id = 110005
  # SAML-specific configuration...
}

# Not recommended for SAML apps
resource "onelogin_apps" "example" {
  name = "Example App"
  connector_id = 110005
  # May not work for all connector types
}
```

## Computed vs. Optional/Required Fields

The provider schema defines certain fields as "Computed" (meaning they're read-only), but the API may actually require values for these fields. The specialized resources handle this better by setting these fields appropriately behind the scenes.

For example, with the basic `onelogin_apps` resource, you can't set `auth_method` or `policy_id` directly because they're marked as "Computed", but the API may require them for certain connector types. This is one reason to prefer the specialized resources for app creation.

## Understanding the Design Choice

This naming convention follows the principle that resource types should describe collections of entities, while resource instances identify specific entities. While this might be unintuitive at first, it's applied consistently throughout the provider.