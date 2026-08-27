---
layout: "onelogin"
page_title: "OneLogin: onelogin_apps"
sidebar_current: "docs-onelogin-resource-apps"
description: |-
  Creates a Basic Application.
---

# onelogin_apps

Creates a Basic Application.

This resource allows you to create and configure a Basic (non-SAML non-OIDC) Application.

## Example Usage

```hcl
resource onelogin_apps my_app {
  connector_id = 12345
  description = "basic app"
  name = "example"
  notes = "basic app"
  visible = true
  allow_assumed_signin = false

  provisioning = {
    enabled = false
  }

  parameters {
    safe_entitlements_enabled = false
    user_attribute_mappings = ""
    provisioned_entitlements = false
    skip_if_blank = false
    user_attribute_macros = ""
    attributes_transformations = ""
    default_values = ["default"]
    include_in_saml_assertion = false
    label = "username"
    param_key_name = "user name"
    values = ["user.email"]
  }
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The app's name.

* `connector_id` - (Required) The ID for the app connector, dictates the type of app (e.g. AWS Multi-Role App).

* `description` - (Optional) App description.

* `notes` - (Optional) Notes about the app.

* `visible` - (Optional) Determine if app should be visible in OneLogin portal. Defaults to `true`.

* `allow_assumed_signin` - (Optional) Enable sign in when user has been assumed by the account owner. Defaults to `false`.

* `policy_id` - (Optional, Number) ID of the app policy enforced when users sign in to
  this app. Create one with [`onelogin_policies`](onelogin_policies.md).

  Only app policies can be assigned. A user policy is refused with
  `The associated Policy with ID <id> could not be found`, the same message the API gives
  for an ID that does not exist.

  The attribute is computed as well as optional, so a policy assigned in the OneLogin
  admin UI is left alone by a configuration that does not mention `policy_id`. The cost
  of that is that *removing* the argument does not unassign the policy — it leaves the
  last value in place. Write `policy_id = 0` to unassign. OneLogin refuses a literal `0`
  and wants `null` for this, so the provider sends the `null` on your behalf.

  ```hcl
  resource "onelogin_policies" "finance_app" {
    name                   = "Finance app step-up"
    kind                   = "app"
    force_authn            = true
    app_force_authn_offset = 60
  }

  resource "onelogin_apps" "finance" {
    name         = "Finance"
    connector_id = 108419
    policy_id    = onelogin_policies.finance_app.id
  }
  ```

* `brand_id` - (Optional, Number) ID of the brand whose login page this app uses. Omit to
  fall back to the account default brand.

  Like `policy_id`, the attribute is computed as well as optional, so a brand assigned in
  the OneLogin admin UI is left alone by a configuration that does not mention
  `brand_id`, and *removing* the argument does not unassign it. Write `brand_id = 0` to
  unassign; OneLogin refuses a literal `0` and wants `null`, which the provider sends on
  your behalf.


* `provisioning` - (Optional) Settings regarding the app's provisioning ability.
  * `enabled` - (Required) Indicates if provisioning is enabled for this app.


* `parameters` - (Optional) a list of custom parameters for this app.
  * `param_key_name` - (Required) Name to represent the parameter in OneLogin.

  * `safe_entitlements_enabled` - (Optional) Indicates that the parameter is used to support creating entitlements using OneLogin Mappings. Defaults to `false`.

  * `user_attribute_mappings` - (Optional) A user attribute to map values from. For custom attributes prefix the name of the attribute with `custom_attribute_`.

  * `provisioned_entitlements` -  (Optional) Provisioned access entitlements for the app. Defaults to `false`.

  * `skip_if_blank` - (Optional)  Flag to let the SCIM provisioner know not include this value if it's blank. Defaults to `false`.

  * `user_attribute_macros` - (Optional) When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.

  * `attributes_transformations` - (Optional) Describes how the app's attributes should be transformed.

  * `default_values` - (Optional) A list of default parameter values. A single value is sent to OneLogin as a plain string and several as an array, matching how the API stores them.

  * `include_in_saml_assertion` - (Optional) When true, this parameter will be included in a SAML assertion payload.

  * `label` - (Optional) The can only be set when creating a new parameter. It can not be updated.

  * `values` - (Optional) A list of parameter values. A single value is sent to OneLogin as a plain string and several as an array, matching how the API stores them.

## Attributes Reference

* `id` - App's unique ID in OneLogin.

* `allow_assumed_signin` - App sign in allowed when user assumed by account administrator.

* `auth_method` - The apps auth method. Refer to the [OneLogin Apps Documentation](https://developers.onelogin.com/api-docs/2/apps/app-resource) for a comprehensive list of available auth methods.

* `connector_id` - ID of the apps underlying connector. Dictates the type of app (e.g. AWS Multi-Role App).

* `description` - App description.

* `icon_url` - The url for the app's icon.

* `name` - The app's name.

* `notes` - Notes about the app.

* `tab_id` - The tab in which to display in OneLogin portal.

* `updated_at` - Timestamp for app's last update.

* `created_at` - Timestamp for app's creation.

* `policy_id` - The app policy assigned to the app. Settable; see the argument above.

* `brand_id` - The brand assigned to the app. Settable; see the argument above.

* `visible` - Indicates if the app is visible in the OneLogin portal.

* `parameters` - The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes.
    * `attributes_transformations` - Describes how the app's attributes should be transformed.

    * `default_values` - The list of default parameter values.

    * `include_in_saml_assertion` - Dictates if the parameter needs to be included in a SAML assertion

    * `label` - The attribute label.

    * `param_id` - The parameter ID.

    * `param_key_name` - The name of the parameter stored in OneLogin.

    * `provisioned_entitlements` - Provisioned access entitlements for the app.

    * `safe_entitlements_enabled` -  Indicates whether entitlements can be created.

    * `skip_if_blank` - Flag to let the SCIM provisioner know not include this value if it's blank.

    * `user_attribute_macros` - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.

    * `user_attribute_mappings` - A user attribute to map values from. For custom attributes the name of the attribute is prefixed with `custom_attribute_`.

    * `values` - The list of parameter values.

* `provisioning` -  Settings regarding the app's provisioning ability.
    * `enabled` - Indicates if provisioning is enabled for this app.

## Import

An App can be imported via the OneLogin App ID.

```
$ terraform import onelogin_apps.my_app <app id>
```
