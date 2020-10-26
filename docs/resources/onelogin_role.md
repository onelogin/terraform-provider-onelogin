---
layout: "onelogin"
page_title: "OneLogin: onelogin_roles"
sidebar_current: "docs-onelogin-resource-roles"
description: |-
  Manage App Rule resources.
---

# onelogin_roles

Manage App Rule resources.

This resource allows you to create and configure App Rules.

## Example Usage - Strict Ordering

```hcl
resource onelogin_roles executive_admin {
  name = "executive admin"
  apps = [123, 456, 787]
  users = [543, 213, 420]
  admins= [777]
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.

* `apps` - (Required) A list of app IDs for which the role applies.

* `users` - (Required) A list of user IDs for whom the role applies.

* `admins` - (Required) A list of IDs of users who administer the role.

## Attributes Reference

No further attributes are exported.

## Import

A role can be imported using the OneLogin Role ID.

```
$ terraform import onelogin_roles.executive_admin <role id>
```
