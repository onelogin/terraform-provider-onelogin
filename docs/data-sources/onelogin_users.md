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

### Looking users up by email

Resources that take user IDs — `onelogin_roles.users` among them — can be given
the results of a lookup rather than IDs written out by hand. `emails` takes a
list, so a whole membership can come from one data source:

```hcl
data onelogin_users engineering {
  emails = [
    "alice@example.com",
    "bob@example.com",
  ]
}

resource onelogin_roles engineering {
  name  = "Engineering"
  users = data.onelogin_users.engineering.users[*].id
}
```

`users[*].id` is already a list of numbers, so it can be passed straight through;
`ids` holds the same values as strings.

Note that an email matching no user contributes nothing rather than failing. If a
typo would be better caught than ignored, compare the counts:

```hcl
locals {
  wanted = ["alice@example.com", "bob@example.com"]
}

data onelogin_users engineering {
  emails = local.wanted
}

check "every_member_resolved" {
  assert {
    condition     = length(data.onelogin_users.engineering.users) == length(local.wanted)
    error_message = "One or more emails in local.wanted matched no OneLogin user."
  }
}
```

## Argument Reference

The following arguments are supported:

* `username` - The user's username.

* `firstname` - The user's first name

* `lastname` - The user's last name

* `email` - The user's email.

* `emails` - A list of emails to look up at once. The API matches one email per
  request, so each is queried in turn and the results combined, in the order
  given, with duplicates removed. Any other argument set here applies to every
  one of those queries. Conflicts with `email`.

* `samaccountname` - The user's samaccount name

* `external_id` - The user's external_id

* `directory_id` - The user's directory_id

## Attributes Reference

* `ids` - List of user's id, as strings

* `users` - List of the matching users, each with `id` (number), `username`,
  `email`, `firstname`, `lastname`, `samaccountname`, `external_id`,
  `directory_id` and `last_login`
