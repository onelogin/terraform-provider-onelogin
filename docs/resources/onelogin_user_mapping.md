---
layout: "onelogin"
page_title: "OneLogin: onelogin_user_mapping"
sidebar_current: "docs-onelogin-resource-user_mapping"
description: |-
  Manage User Mappings resources.
---

# onelogin_user_mappings

Manage User Mappings resources.

This resource allows you to create and configure User Mappings.

## Example Usage

```hcl
resource onelogin_user_mappings example {
  name = "Select Login"
  enabled = true
  match = "all"
  position = 1

  actions = {
    value = ["1"]
    action = "set_status"
  }

  conditions = {
    operator = ">"
    source = "last_login"
    value = "90"
  }
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The resource's name.

* `enabled` - (Required) Indicates if a mapping is enabled.

* `match` - (Required) Indicates how conditions should be matched. Must be one of `all` or `any`.

* `position` - (Required) Indicates the ordering of the mapping. When `null` this will be placed last.

* `conditions` - (Required) An array of conditions that the user must meet in order for the mapping to be applied.
  * `source` - (Required) The source field to check. See [List Conditions](https://developers.onelogin.com/api-docs/2/user-mappings/list-conditions) for possible values.

  * `operator` - (Required) A valid operator for the selected condition source. See [List Condition Operators](https://developers.onelogin.com/api-docs/2/user-mappings/list-condition-operators) for possible values.

  * `value` - (Required) A plain text string or valid value for the selected condition source. See [List Condition Values](https://developers.onelogin.com/api-docs/2/user-mappings/list-condition-values) for possible values.

* `actions` - (Required) The number of minutes until the token expires
  * `action` - (Required) The action to apply. See [List Actions](https://developers.onelogin.com/api-docs/2/user-mappings/list-conditions) for possible values.

  * `value` - (Required) An array of strings. Items in the array will be a plain text string or valid value for the selected action. See [List Action Values](https://developers.onelogin.com/api-docs/2/user-mappings/list-action-values) for possible values. In most cases only a single item will be accepted in the array.



## Attributes Reference

No further attributes are exported

## Import

A User Mapping can be imported via the OneLogin User Mapping.

```
$ terraform import onelogin_user_mappings.example <user_mapping_id>
```
