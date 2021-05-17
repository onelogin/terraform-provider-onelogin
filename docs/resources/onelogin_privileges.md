---
layout: "onelogin"
page_title: "OneLogin: onelogin_privileges"
sidebar_current: "docs-onelogin-resource-privileges"
description: |-
  Manage App Rule resources.
---

# onelogin_privileges

Manage Privilege resources.

This resource allows you to create and configure Privilege.

## Example Usage - Strict Ordering

```hcl
resource onelogin_privileges super_admin {
    name = "super duper admin"
    description = "description"
    user_ids = [123, 345]
    role_ids = [987, 654]
    privilege {
        statement {
            effect = "Allow"
            action = ["apps:List"]
            scope = ["*"]
        }
        statement {
            effect = "Allow"
            action = ["users:List", "users:Update"]
            scope = ["users/123", "users/345]
        }
    }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the privilege.

* `description` - (Optional) Description for the Privilege.

* `user_ids` - (Optional) A list of user IDs for whom the privilege applies.

* `role_ids` - (Optional) A list of role IDs for whom the role applies.

* `privilege` - (Required) A list of statements that describe what the privilege grants access to.
  
  * `statement` - (Required) At least one `statement` is required. Statements describe the effect granted to a resource type. In this case it allow's the privilege holder to lisst apps and users.
  
    *  `effect` - (Required) The effect the privilege grants for the resource. Must be "Allow".
    
    *  `action` - (Required) List of actions the privilege holder can do. Must be one of those [listed in the docs](https://developers.onelogin.com/api-docs/1/privileges/create-privilege)

    * `scope` - (Required) Target the privileged action against specific resources with the scope. In this case, the privilege only grants update access to users 123 and 345.

## Attributes Reference

No further attributes are exported.

## Import

A privilege can be imported using the OneLogin Privilege ID.

```
$ terraform import onelogin_privilegess.super_admin <privilege id>
```
