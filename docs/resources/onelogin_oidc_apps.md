---
layout: "onelogin"
page_title: "OneLogin: onelogin_oidc_apps"
sidebar_current: "docs-onelogin-resource-oidc-apps"
description: |-
  Creates a OIDC Application.
---

# onelogin_oidc_apps

Creates an OIDC Application.

This resource allows you to create and configure an OIDC Application.

## Example Usage

```hcl
resource onelogin_oidc_apps my_oidc_app {
  name = "my OIDC APP"
  notes = "example"
  visible = true
  allow_assumed_signin = false
  connector_id = 123456
  description = "example OIDC app"

  configuration = {
    post_logout_redirect_uri = "https://www.example.com/logout"
    login_url = "https://www.example.com"
    oidc_application_type = 0
    redirect_uri = "https://example.com/example"
    refresh_token_expiration_minutes = 1
    token_endpoint_auth_method = 1
    access_token_expiration_minutes = 1
  }

  provisioning = {
    enabled = false
  }

  parameters {
    provisioned_entitlements = false
    user_attribute_macros = ""
    user_attribute_mappings = ""
    values = ["user.email"]
    attributes_transformations = ""
    default_values = ["default"]
    include_in_saml_assertion = false
    label = "example app parameter "
    param_key_name = "example"
    safe_entitlements_enabled = false
    skip_if_blank = false
  }
}
```

## Argument Reference

The following arguments are supported:
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

  resource "onelogin_oidc_apps" "finance" {
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


* `configuration` - OIDC settings that control the authentication flow e.g. redirect urls and token settings.
  * `post_logout_redirect_uri` - (Optional) The redirect_uri for the app to send the user to after logout. To allow more than one, give a single string with the URIs separated by newlines or commas. Set it to `""` to remove every URI from the app. Omitting the attribute entirely leaves whatever the app already has in place, so that removing it from your configuration is not the same as clearing it.

  * `redirect_uri` - (Optional) The redirect_uri for the OIDC flow. Will be computed by API if not given.

  * `refresh_token_expiration_minutes` - (Optional) Number of minutes for the refresh token to be valid. Defaults to 1 minute.

  * `login_url` - (Optional) The login_url for the OIDC flow. Will be computed by API if not given.

  * `oidc_application_type` - (Optional) Must be one of one of `0` (Web) or `1` (Native/Mobile). Defaults to 0.

  * `token_endpoint_auth_method` - (Optional) Must be one of one of `0` (Basic) `1` (POST) `2` (Nonce/PKCE).

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

    * `label` - The attribute label

    * `param_id` - The parameter ID.

    * `param_key_name` - The name of the parameter stored in OneLogin.

    * `provisioned_entitlements` -  Provisioned access entitlements for the app.

    * `safe_entitlements_enabled` -  Indicates whether entitlements can be created.

    * `skip_if_blank` - Flag to let the SCIM provisioner know not include this value if it's blank.

    * `user_attribute_macros` - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.

    * `user_attribute_mappings` - A user attribute to map values from. For custom attributes the name of the attribute is prefixed with `custom_attribute_`.

    * `values` - The list of parameter values.

* `provisioning` -  Settings regarding the app's provisioning ability.
    * `enabled` - Indicates if provisioning is enabled for this app.


* `configuration`
  * `redirect_uri` - The redirect_uri for the OIDC flow.

  * `refresh_token_expiration_minutes` - Number of minutes for the refresh token to be valid.

  * `login_url` - The login_url for the OIDC flow.

  * `oidc_application_type` - Indicates OIDC app type.

  * `token_endpoint_auth_method` - Indicates the token endpoint authentication method.

* `sso` - The OIDC client credentials OneLogin generated for the app. `sso` is read only and cannot be set in configuration; OneLogin supplies every value.

  * `client_id` - The OIDC client ID.

  * `client_secret` - The OIDC client secret. Per the [OneLogin Apps Documentation](https://developers.onelogin.com/api-docs/2/apps/app-resource), `client_secret` is only ever returned by OneLogin when an app is created, so the provider captures it at create time and retains it in state thereafter. It is therefore available for apps **created by Terraform**, and **not available for imported apps** -- an import has no create response, and no later read can recover the secret.

  If OneLogin re-issues the app's credentials, the following read returns a different `client_id`. The provider treats the retained secret as stale and drops it, rather than pairing it with credentials it does not belong to, and emits a warning. A rotation that leaves `client_id` unchanged cannot be detected, because the provider has no way to read the current secret; in that case the value in state stays stale silently. Either way the app must be recreated for Terraform to hold a valid secret again.

  `client_secret` is omitted from the map -- rather than present with an empty value -- whenever it has not been captured, so indexing it directly (`sso.client_secret`) fails with an "Invalid index" error instead of quietly returning an empty string. For a resource already in state that failure comes at plan time; for one not yet created the map is unknown until apply, so the error lands there instead.

  **If you feed the secret into another system, index it directly.** A default-returning form such as `lookup(..., "client_secret", "")` turns every one of the cases above into an empty string, and an empty string will be written to whatever consumes it -- silently replacing a working credential with a blank one. Failing at plan time is the safer outcome. Use `lookup()` or `try()` with a default only where an absent secret is genuinely acceptable, such as a module that must also accept imported apps and public clients.

  The whole `sso` map is marked sensitive because it carries the client secret. Terraform cannot mark individual map keys, so `client_id` is redacted too and needs `nonsensitive()` to be used in a non-sensitive output:

  ```hcl
  output "oidc_client_id" {
    value = nonsensitive(onelogin_oidc_apps.my_oidc_app.sso.client_id)
  }

  # Direct index: fails loudly at plan time if the secret was never captured,
  # rather than handing an empty string to whatever consumes this output.
  output "oidc_client_secret" {
    value     = onelogin_oidc_apps.my_oidc_app.sso.client_secret
    sensitive = true
  }
  ```

  Note that sensitivity controls display only. The client secret is written in plaintext to Terraform state, to `terraform.tfstate.backup`, to remote state, and to any saved plan file -- `terraform plan -out=...` embeds prior state, so plan artifacts kept by CI contain the secret too. `terraform show -json` exposes it as well. Treat all of them as secrets.

## Import

A OIDC App can be imported via the OneLogin App ID.

```
$ terraform import onelogin_oidc_apps.my_oidc_app <app id>
```
