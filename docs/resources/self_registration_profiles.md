# onelogin_self_registration_profiles Resource

This resource allows you to create and configure Self-Registration Profiles in OneLogin.

## Example Usage

```hcl
resource "onelogin_self_registration_profiles" "example" {
  name                  = "Example Profile"
  url                   = "example-profile"
  enabled               = true
  moderated             = false
  default_role_id       = 123456
  default_group_id      = 789012
  helptext              = "Welcome to our community. Please follow the guidelines."
  thankyou_message      = "Thank you for joining!"
  domain_blacklist      = "spam.com, banned.com"
  domain_whitelist      = "example.com, trusted.com"
  domain_list_strategy  = 1
  email_verification_type = "Email MagicLink"
  
  fields {
    custom_attribute_id = 12345
  }
  
  fields {
    custom_attribute_id = 67890
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the self-registration profile.
* `url` - (Required) The URL path for the self-registration profile.
* `enabled` - (Optional) Whether the self-registration profile is enabled. Defaults to `true`.
* `moderated` - (Optional) Whether the self-registration profile requires moderation. Defaults to `false`.
* `default_role_id` - (Optional) The default role ID to assign to users who register through this profile.
* `default_group_id` - (Optional) The default group ID to assign to users who register through this profile.
* `helptext` - (Optional) Help text displayed on the registration page.
* `thankyou_message` - (Optional) Thank you message displayed after registration.
* `domain_blacklist` - (Optional) Comma-separated list of domains to blacklist.
* `domain_whitelist` - (Optional) Comma-separated list of domains to whitelist.
* `domain_list_strategy` - (Optional) Domain list strategy: `0` for blacklist, `1` for whitelist. Defaults to `0`.
* `email_verification_type` - (Optional) Email verification type: `Email MagicLink` or `Email OTP`. Defaults to `Email MagicLink`.
* `fields` - (Optional) Custom fields for the self-registration profile.
  * `custom_attribute_id` - (Required) ID of the custom attribute.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the self-registration profile.
* `fields` - Custom fields for the self-registration profile.
  * `id` - ID of the field.
  * `name` - Name of the field.

## Import

Self-Registration Profiles can be imported using the ID, e.g.

```
$ terraform import onelogin_self_registration_profiles.example 12345
```
