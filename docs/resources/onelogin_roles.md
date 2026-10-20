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

* `apps` - (Optional) A list of app IDs to associate with this role. Users assigned to this role will have access to these applications. Leave it out to leave the role's apps unmanaged, for example when they are attached with `onelogin_app_role_attachments` or an app's `role_ids`. Set `apps = []` to remove every app.

* `users` - (Optional) A list of user IDs to assign to this role. These users will gain access to any apps associated with this role. Leave it out to leave the role's users unmanaged, for example when a mapping assigns them. Set `users = []` to remove every user.

* `admins` - (Optional) A list of user IDs who will be administrators for this role. These users can manage the role settings. Leave it out to leave the role's admins unmanaged. Set `admins = []` to remove every admin.

* `skip_membership_refresh` - (Optional) A list of membership attributes to leave
  un-refreshed: any of `apps`, `users`, `admins`. Their sub-endpoints are not
  queried during read, and whatever state already holds is kept.

  Membership lives on three separate paginated endpoints, so a role with
  thousands of users costs a page walk on every `plan` or `refresh` — whether or
  not your configuration manages that membership. `lifecycle { ignore_changes }`
  cannot prevent it, because it is applied to the diff, after the read has
  already made the calls.

  Use this where membership is managed outside Terraform, for example by
  mappings, and leave the same attribute out of your configuration:

  ```hcl
  resource onelogin_roles engineering {
    name = "Engineering"

    # Membership comes from a mapping, so do not read it. Leaving users out
    # means Terraform does not manage it either.
    skip_membership_refresh = ["users"]
  }
  ```

  An attribute left out of the configuration is never diffed, so
  `ignore_changes` for it is no longer needed, though an existing one does no
  harm. Listing an attribute here while also setting it will leave whatever
  state already held, which the plan will then compare against your
  configuration.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The ID of the role.

## Import

A role can be imported using the OneLogin Role ID.

```
$ terraform import onelogin_roles.executive_admin <role id>
```

## Notes

A membership attribute you set is authoritative. For example, if you want to add a new user to a role while keeping the existing users, you must include both the existing and new user IDs in the `users` attribute. Otherwise, the existing users will be removed from the role.

Removing `apps`, `users` or `admins` from your configuration stops managing it and leaves the current members in place. Set it to `[]` to remove them all.