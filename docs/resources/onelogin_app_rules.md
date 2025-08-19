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

## Example Usage - Strict Ordering

```hcl
resource onelogin_app_rules check{
  app_id = onelogin_saml_apps.my_saml_app.id
  position = 1
  enabled = true
  match = "all"
  name = "first rule"
  conditions = {
    operator = "ri"
    source = "has_role"
    value = "340475"
  }
  actions = {
    action = "set_amazonusername"
    expression = ".*"
    value = ["member_of"]
  }
}
```

## Example Usage - Dependency Based Ordering

```hcl
resource onelogin_app_rules test{
  app_id = onelogin_saml_apps.my_saml_app.id
  enabled = true
  match = "all"
  name = "first rule"
  conditions = {
    operator = "ri"
    source = "has_role"
    value = "340475"
  }
  actions = {
    action = "set_amazonusername"
    expression = ".*"
    value = ["member_of"]
  }
}

resource onelogin_app_rules check{
  app_id = onelogin_saml_apps.my_saml_app.id
  depends_on = [onelogin_app_rules.test]
  enabled = true
  match = "all"
  name = "second rule"
  conditions = {
    operator = "ri"
    source = "has_role"
    value = "340475"
  }
  actions = {
    action = "set_amazonusername"
    expression = ".*"
    value = ["member_of"]
  }
}
```

## Important Note Regarding Position

The position field indicates the order in which rules are applied. They behave like progressive filters and as such, their positioning is strictly enforced. Your options for this field are to either:

* Accept any ordering - Do not fill out any position field and each rule will be inserted in the order received by the API.

* Strict Ordering - Enter a position number for each app rule. You'll need to ensure there are no duplicates and no gaps in numbering.

* Dependency based ordering - Use the `depends_on` field to specify an app rule's predecessor to ensure rules are received by the API in the order in which they should be applied. e.g. `depends_on = [onelogin_app_rules.test]`

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The id of the App resource to which the rule should belong.

* `enabled` - (Required) Indicate if the rule should go into effect.

* `match` - (Required) Indicates how conditions should be matched. Must be one of `all` or `any`.

* `name` - (Required) The Rule's name

* `position` - (Optional) Indicates the ordering of the rule. When not supplied the rule will be put at the end of the list on create and managed by the provider. '0' can be supplied to consistently push this rule to the end of the list on every update.

* `conditions` - (Required) An array of conditions that the user must meet in order for the rule to be applied.
  * `source` - The source field to check. See [List Conditions](https://developers.onelogin.com/api-docs/2/app-rules/list-conditions) for possible values.
  * `operator` - A valid operator for the selected condition source. See [List Condition Operators](https://developers.onelogin.com/api-docs/2/app-rules/list-condition-operators) for possible values.
  * `value` - A plain text string or valid value for the selected condition source. See [List Condition Values](https://developers.onelogin.com/api-docs/2/app-rules/list-condition-values) for possible values.

* `actions` - (Required) An array of actions that will be applied to the users that are matched by the conditions.
  * `action` - The action to apply. See [List Actions](https://developers.onelogin.com/api-docs/2/app-rules/list-conditions) for possible values. *Note*: The action `set_role_from_existing` may also be used, however doing so will always clear the `expression` field as it is not accepted when mapping a rule from existing roles.
  * `value` - An array of strings. Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action. See [List Action Values](https://developers.onelogin.com/api-docs/2/app-rules/list-action-values) for possible values. In most cases only a single item will be accepted in the array.
  * `expression` - A regular expression to extract a value. Applies to provisionable, multi-selects, and string actions.
  * `scriptlet` - A hash containing scriptlet code that returns a value. Scriptlets can not be modified and the same hash should not be applied to other applications.
  * `macro` - A template to construct a value. Applies to default, string, and list actions.

## Attributes Reference

No further attributes are exported.

## Import

An App Rule cannot be imported at this time.
