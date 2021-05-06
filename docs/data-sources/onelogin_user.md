---
layout: "onelogin"
page_title: "OneLogin: onelogin_user"
sidebar_current: "docs-onelogin-resource-user"
description: |-
  Returns User resource.
---

# Data source: onelogin_user

Returns User resource.

## Example Usage

```hcl
data onelogin_user example {
  username = "timmy.tester"
}
```

## Argument Reference

The following arguments are supported:

* `username` - The user's username.

* `user_id` - The user's ID.

## Attributes Reference

* `id` - The user's id

* `email` - The user's email.

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
