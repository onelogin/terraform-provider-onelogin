---
layout: "onelogin"
page_title: "OneLogin: onelogin_app_rules"
sidebar_current: "docs-onelogin-resource-app_rules"
description: |-
  Manage App Rule resources.
---

# onelogin_app_rules

Manage App Rule resources.

This resource allows you to create and configure App Rules.

## Example Usage

```hcl
resource onelogin_app_rules check{
  app_id = onelogin_saml_apps.my_saml_app.id
  enabled = true
  match = "all"
  name = "second rule"
  conditions {
    operator = "ri"
    source = "has_role"
    value = "340475"
  }
  actions {
    action = "set_amazonusername"
    expression = ".*"
    value = ["member_of"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The id of the App resource to which the rule should belong.

* `enabled` - (Required) Indicate if the rule should go into effect.

* `match` - (Required) Indicates how conditions should be matched. Must be one of `all` or `any`.

* `name` - (Required) The Rule's name

* `position` - (Required) Indicates the order of the rule. When `null` this will default to last position.

* `conditions` - (Required) An array of conditions that the user must meet in order for the rule to be applied.
  * `source` - The source field to check. See [List Conditions](https://developers.onelogin.com/api-docs/2/app-rules/list-conditions) for possible values.
  * `operator` - A valid operator for the selected condition source. See [List Condition Operators](https://developers.onelogin.com/api-docs/2/app-rules/list-condition-operators) for possible values.
  * `value` - A plain text string or valid value for the selected condition source. See [List Condition Values](https://developers.onelogin.com/api-docs/2/app-rules/list-condition-values) for possible values.

* `actions` - (Required) An array of actions that will be applied to the users that are matched by the conditions.
  * `action` - The action to apply. See [List Actions](https://developers.onelogin.com/api-docs/2/app-rules/list-conditions) for possible values.
  * `value` - An array of strings. Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action. See [List Action Values](https://developers.onelogin.com/api-docs/2/app-rules/list-action-values) for possible values. In most cases only a single item will be accepted in the array.
  * `expression` - A regular expression to extract a value. Applies to provisionable, multi-selects, and string actions.
  * `scriptlet` - A hash containing scriptlet code that returns a value. Scriptlets can not be modified and the same hash should not be applied to other applications.
  * `macro` - A template to construct a value. Applies to default, string, and list actions.

## Attributes Reference

No further attributes are exported.

## Import

An App Rule cannot be imported at this time.
