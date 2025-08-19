---
layout: "onelogin"
page_title: "OneLogin: onelogin_roles"
sidebar_current: "docs-onelogin-resource-roles"
description: |-
  Manage Role resources in OneLogin.
---

# onelogin_roles

Manage Role resources in OneLogin.

This resource allows you to create and configure Roles, including assigning applications, users, and administrators to roles.

## Example Usage

```hcl
resource onelogin_roles executive_admin {
  name = "Executive Admin"
  
  # Optional: assign apps to this role
  apps = [123, 456, 787]
  
  # Optional: assign users to this role
  users = [543, 213, 420]
  
  # Optional: assign administrators to this role
  admins = [777]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.

* `apps` - (Optional) A list of app IDs to associate with this role. Users assigned to this role will have access to these applications. If not specified, no apps will be associated with the role.

* `users` - (Optional) A list of user IDs to assign to this role. These users will gain access to any apps associated with this role. If not specified, no users will be assigned to the role.

* `admins` - (Optional) A list of user IDs who will be administrators for this role. These users can manage the role settings. If not specified, no admins will be assigned to the role.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The ID of the role.

## Import

A role can be imported using the OneLogin Role ID.

```
$ terraform import onelogin_roles.executive_admin <role id>
```

## Notes

When updating a role, you must specify all fields you want to maintain. For example, if you want to add a new user to a role while keeping the existing users, you must include both the existing and new user IDs in the `users` attribute. Otherwise, the existing users will be removed from the role.