---
layout: "onelogin"
page_title: "OneLogin: onelogin_groups"
sidebar_current: "docs-onelogin-resource-groups"
description: |-
  Manages a OneLogin group.
---

# onelogin_groups

Manages a OneLogin group. Supports full CRUD operations via the V2 API.

## Example Usage

```hcl
resource "onelogin_groups" "engineering" {
  name      = "Engineering"
  reference = "eng-group"
}
```

Applying a security policy to the group's members:

```hcl
resource "onelogin_policies" "engineering" {
  name                     = "Engineering"
  kind                     = "user"
  password_expiration_days = 90
}

resource "onelogin_groups" "engineering" {
  name      = "Engineering"
  policy_id = onelogin_policies.engineering.id
}
```

## Argument Reference

* `name` - (Required) The name of the group.
* `reference` - (Optional) A reference identifier for the group.
* `policy_id` - (Optional, Number) ID of the user policy applied to this group's members.

  Only user policies can be assigned. An app policy is refused with
  `Policy must reference a user policy`.

  Removing the argument clears the assignment, and the group's members fall back to the
  account default policy. Leaving it out of a configuration that never set it does nothing:
  an unassigned group and an absent argument are both zero, so a plan stays empty.

## Attribute Reference

* `id` - The ID of the group.

## Import

Groups can be imported using the group ID:

```
$ terraform import onelogin_groups.engineering 123456
```
