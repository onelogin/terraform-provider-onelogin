---
layout: "onelogin"
page_title: "OneLogin: onelogin_users"
sidebar_current: "docs-onelogin-resource-user"
description: |-
  Returns User IDs matching the given attributes.
---

# Data source: onelogin_users

Returns User IDs matching the given attributes.

## Example Usage

```hcl
data onelogin_users example {
  firstname = "tom"
}
```

## Argument Reference

The following arguments are supported:

* `username` - The user's username.

* `firstname` - The user's first name

* `lastname` - The user's last name

* `email` - The user's email.

* `samaccountname` - The user's samaccount name

* `external_id` - The user's external_id

* `directory_id` - The user's directory_id

## Attributes Reference

* `ids` - List of user's id
