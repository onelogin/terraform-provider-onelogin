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

  provisioning {
    enabled = false
  }

	parameters {
		safe_entitlements_enabled = false
		user_attribute_mappings = ""
		provisioned_entitlements = false
		skip_if_blank = false
		user_attribute_macros = ""
		attributes_transformations = ""
		default_values = ""
		include_in_saml_assertion = false
		label = "username"
		param_key_name = "user name"
		values = ""
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

  * `default_values` - (Optional) Default Parameter values.

  * `include_in_saml_assertion` - (Optional) When true, this parameter will be included in a SAML assertion payload.

  * `label` - (Optional) The can only be set when creating a new parameter. It can not be updated.

  * `values` - (Optional) Parameter values.

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

* `policy_id` - The security policy assigned to the app.

* `visible` - Indicates if the app is visible in the OneLogin portal.

* `parameters` - The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes.
    * `attributes_transformations` - Describes how the app's attributes should be transformed.

    * `default_values` -  Default Parameter values.

    * `include_in_saml_assertion` - Dictates if the parameter needs to be included in a SAML assertion

    * `label` - The attribute label.

    * `param_id` - The parameter ID.

    * `param_key_name` - The name of the parameter stored in OneLogin.

    * `provisioned_entitlements` - Provisioned access entitlements for the app.

    * `safe_entitlements_enabled` -  Indicates whether entitlements can be created.

    * `skip_if_blank` - Flag to let the SCIM provisioner know not include this value if it's blank.

    * `user_attribute_macros` - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.

    * `user_attribute_mappings` - A user attribute to map values from. For custom attributes the name of the attribute is prefixed with `custom_attribute_`.

    * `values` - Parameter values.

* `provisioning` -  Settings regarding the app's provisioning ability.
    * `enabled` - Indicates if provisioning is enabled for this app.

## Import

An App can be imported via the OneLogin App ID.

```
$ terraform import onelogin_apps.my_app <app id>
```
