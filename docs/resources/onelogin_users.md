---
layout: "onelogin"
page_title: "OneLogin: onelogin_users"
sidebar_current: "docs-onelogin-resource-user"
description: |-
  Manage User resources.
---

# onelogin_users

Manage User resources.

This resource allows you to create and configure Users.

## Example Usage

```hcl
resource onelogin_users example {
  username = "timmy.tester"
  email = "timmy.tester@test.com"
}
```

## Argument Reference

The following arguments are supported:
* `username` - (Required) The user's username.

* `email` - (Required) The user's email.

* `firstname` - The user's first name

* `lastname` - The user's last name

* `distinguished_name` - The user's distinguished name

* `samaccountname` - The user's samaccount name

* `userprincipalname` - The user's user principal name

* `member_of` - The user's member_of

* `phone` - The user's phone number

* `title` - The user's title

* `company` - The user's company

* `department` - The user's department

* `comment` - A comment about the user

* `state` - The user's state. Must be one of `0: Unapproved` `1: Approved` `2: Rejected` `3: Unlicensed`

* `status` - The user's status. Must be one of `0: Unactivated` `1: Active` `2: Suspended` `3: Locked` `4: Password expired` `5: Awaiting password reset` `7: Password Pending` `8: Security questions required`

* `group_id` - The user's group_id

* `directory_id` - The user's directory_id

* `trusted_idp_id` - The user's trusted_idp_id

* `manager_ad_id` - The user's manager_ad_id

* `manager_user_id` - The user's manager_user_id

* `external_id` - The user's external_id

* `password` - (Optional, Sensitive) The user's password. This field is sensitive and will not be displayed in logs or output.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The user's id

## Import

A User can be imported via the OneLogin User ID.

```
$ terraform import onelogin_users.example 12345678
```
